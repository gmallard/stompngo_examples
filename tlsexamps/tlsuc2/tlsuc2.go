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
Connect and Disconnect from a STOMP broker with a TLS connection, use case 2.

	TLS Use Case 2 - client *does* authenticate broker.

	Subcase 2.A - Message broker configuration does *not* require client authentication

		- Expect connection success because  the client did authenticate the
  		broker's certificate.

	Subcase 2.B - Message broker configuration *does* require client authentication

	- Expect connection failure (broker must be sent a valid client certificate)

	Example use might be:

		go build
		./tlsuc2 -srvCAFile=/ad3/gma/sslwork/2013/TestCA.crt # PEM format file

*/
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io/ioutil"
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
	exampid   = "tlsuc2:"
	tc        *tls.Config
	srvCAFile string // Name of file with broker's CA certificate, PEM format

	ll = log.New(os.Stdout, "TLSU2 ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

func init() {
	flag.StringVar(&srvCAFile, "srvCAFile", "DUMMY", "Name of file with broker CA certificate")
}

// Connect to a STOMP broker using TLS and disconnect.
func main() {
	ll.Printf("%s starts\n", exampid)

	flag.Parse() // Parse flags
	ll.Printf("%s using srvCAFile srvCAFile:%s\n",
		exampid, srvCAFile)

	// TLS Configuration.
	tc = new(tls.Config)
	tc.InsecureSkipVerify = false // *Do* check the broker's certificate

	// Get host and port
	h, p := senv.HostAndPort()
	ll.Printf("%s host_and_port host:%s port:%s\n", exampid, h, p)

	// Be polite, allow SNI (Server Virtual Hosting)
	tc.ServerName = h

	// Finish TLS Config initialization, so client can authenticate broker.
	b, e := ioutil.ReadFile(srvCAFile) // Read broker's CA cert (PEM)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error())
	}
	k, _ := pem.Decode(b) // Decode PEM format (*pem.Block)
	if k == nil {
		ll.Fatalf("%s %s\n", exampid, e.Error())
	}
	c, e := x509.ParseCertificate(k.Bytes) // Create *x509.Certificate
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error())
	}

	tc.RootCAs = x509.NewCertPool() // Create a cert "pool"
	tc.RootCAs.AddCert(c)           // Add the CA cert to the pool

	// Connect logic: use net.Dial and tls.Client
	hap := net.JoinHostPort(h, p)
	t, e := net.Dial("tcp", hap)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s dial_complete\n", exampid)

	nc := tls.Client(t, tc)
	e = nc.Handshake()
	if e != nil {
		if e.Error() == "EOF" {
			ll.Printf("%s handshake EOF, Is the broker port TLS enabled? port:%s\n",
				exampid, p)
		}
		ll.Fatalf("%s v1:%v v2:%v\n", exampid, "netHandshake", e) // Handle this .....
	}
	ll.Printf("%s handshake_complete\n", exampid)

	sngecomm.DumpTLSConfig(exampid, tc, nc)

	// stompngo Connect Headers
	ch := sngecomm.ConnectHeaders()
	// Get a stomp connection.  Parameters are:
	// a) the opened net connection
	// b) the connect Headers
	conn, e := stompngo.Connect(nc, ch)
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
	e = nc.Close()
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s connsess:%s net_close_complete\n",
		exampid, conn.Session())

	ll.Printf("%s ends\n", exampid)
}
