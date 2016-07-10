//
// Copyright © 2013-2016 Guy M. Allard
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

/*
Connect and Disconnect from a STOMP broker with a TLS connection, use case 3.

	TLS Use Case 3 - broker *does* authenticate client, client does *not* authenticate broker

	Subcase 3.A - Message broker configuration does *not* require client authentication

	- Expect connection success

	Subcase 3.B - Message broker configuration *does* require client authentication

	- Expect connection success if the broker can authenticate the client certificate

	Example use might be:

		go build
		./tlsuc3 -cliCertFile=/ad3/gma/sslwork/2013/client.crt -cliKeyFile=/ad3/gma/sslwork/2013/client.key

*/
package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	//
	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid     = "tlsuc3:"
	tc          *tls.Config
	cliCertFile string
	cliKeyFile  string
	ll          = log.New(os.Stdout, "TLSU3 ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

func init() {
	flag.StringVar(&cliCertFile, "cliCertFile", "DUMMY_CERT", "Name of client cert file")
	flag.StringVar(&cliKeyFile, "cliKeyFile", "DUMMY_KEY", "Name of client key file")
}

// Connect to a STOMP broker using TLS and disconnect.
func main() {
	ll.Printf("%s starts\n", exampid)

	flag.Parse() // Parse flags
	ll.Printf("%s using cliCertFile cliCertFile:%s\n",
		exampid, cliCertFile)
	ll.Printf("%s using cliKeyFile cliKeyFile:%s\n",
		exampid, cliKeyFile)

	// TLS Configuration.
	tc = new(tls.Config)
	tc.InsecureSkipVerify = true // Do *not* check the broker's certificate

	// Get host and port
	h, p := senv.HostAndPort()
	ll.Printf("%s host_and_port host:%s port:%s\n", exampid, h, p)

	// Be polite, allow SNI (Server Virtual Hosting)
	tc.ServerName = h

	// Finish TLS Config initialization, so broker can authenticate client.
	// cc -> tls.Certificate
	cc, e := tls.LoadX509KeyPair(cliCertFile, cliKeyFile)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	tc.Certificates = append(tc.Certificates, cc) // Add cert to config

	// This is OK, but does not seem to be required
	tc.BuildNameToCertificate() // Build names map

	// Connect logic: use net.Dial and tls.Client
	t, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s dial_complete\n", exampid)

	n := tls.Client(t, tc)
	e = n.Handshake()
	if e != nil {
		if e.Error() == "EOF" {
			ll.Printf("%s handshake EOF, Is the broker port TLS enabled? port:%s\n",
				exampid, p)
		}
		ll.Fatalln(exampid, "netHandshake", e) // Handle this .....
	}

	ll.Printf("%s handshake_complete\n", exampid)
	sngecomm.DumpTLSConfig(exampid, tc, n)

	// Connect Headers
	ch := sngecomm.ConnectHeaders()

	// Get a stomp connection.  Parameters are:
	// a) the opened net connection
	// b) the connect Headers
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s connsess:%s stomp_connect_complete\n",
		exampid, conn.Session())

	// *NOTE* your application functionaltiy goes here!

	// Polite Stomp disconnects are not required, but highly recommended.
	// Empty headers here.
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s connsess:%s stomp_disconnect_complete\n",
		exampid, conn.Session())

	// Close the net connection.
	e = n.Close()
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s connsess:%s net_close_complete\n",
		exampid, conn.Session())

	ll.Printf("%s ends\n", exampid)
}
