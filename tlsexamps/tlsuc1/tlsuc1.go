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
Connect and Disconnect from a STOMP broker with a TLS connection, use case 1.

	TLS Use Case 1 - client does *not* authenticate broker.

	Subcase 1.A - Message broker configuration does *not* require client authentication

	- Expect connection success

	Subcase 1.B - Message broker configuration *does* require client authentication

	- Expect connection failure (broker must be sent a valid client certificate)

	Example use might be:

		go build
		./tlsuc1

*/
package main

import (
	"crypto/tls"
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
	exampid = "tlsuc1:"
	tc      *tls.Config

	ll = log.New(os.Stdout, "TLSU1 ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

// Connect to a STOMP broker using TLS and disconnect.
func main() {
	ll.Printf("%s starts\n", exampid)

	// TLS Configuration.
	tc = new(tls.Config)
	tc.InsecureSkipVerify = true // Do *not* check the server's certificate

	// Get host and port
	h, p := senv.HostAndPort()
	ll.Printf("%s host_and_port host:%s port:%s\n", exampid, h, p)

	// Be polite, allow SNI (Server Virtual Hosting)
	tc.ServerName = h

	// Connect logic: use net.Dial and tls.Client
	hap := net.JoinHostPort(h, p)
	n, e := net.Dial("tcp", hap)
	if e != nil {
		ll.Fatalln(exampid, "nedDial", e) // Handle this ......
	}
	ll.Printf("%s dial_complete\n", exampid)

	nc := tls.Client(n, tc)
	e = nc.Handshake()
	if e != nil {
		if e.Error() == "EOF" {
			ll.Printf("%s handshake EOF, Is the broker port TLS enabled? port:%s\n",
				exampid, p)
		}
		ll.Fatalln(exampid, "netHandshake", e) // Handle this ......
	}
	ll.Printf("%s handshake_complete\n", exampid)
	sngecomm.DumpTLSConfig(exampid, tc, nc)

	// Connect Headers
	ch := sngecomm.ConnectHeaders()
	// Get a stomp connection.  Parameters are:
	// a) the opened net connection
	// b) the connect Headers
	conn, e := stompngo.Connect(nc, ch)
	if e != nil {
		ll.Fatalln(exampid, "sngConnect", e) // Handle this ......
	}
	ll.Printf("%s connsess:%s stomp_connect_complete\n",
		exampid, conn.Session())

	// *NOTE* your application functionaltiy goes here!

	// Polite Stomp disconnects are not required, but highly recommended.
	// Empty headers here.
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalln(exampid, "sngDisconnect", e) // Handle this ......
	}
	ll.Printf("%s connsess:%s stomp_disconnect_complete\n",
		exampid, conn.Session())

	// Close the net connection.
	e = nc.Close()
	if e != nil {
		ll.Fatalln(exampid, "netClose", e) // Handle this ......
	}
	ll.Printf("%s connsess:%s net_close_complete\n",
		exampid, conn.Session())

	ll.Printf("%s ends\n", exampid)
}
