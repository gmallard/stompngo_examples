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
Connect and Disconnect from a STOMP broker using TLS.

All examples in the 'conndisc' directory also apply here.  This example shows
that TLS is requested by using a specific port and tls.Dial.

	Example:

		# Connect to a broker using TLS:
		STOMP_PORT=61612 go run conndisc_tls.go

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
	exampid = "conndisc_tls: "
	tc      *tls.Config
	ll      = log.New(os.Stdout, "TLCD ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

// Connect to a STOMP 1.0 broker using TLS and disconnect.
func main() {
	ll.Printf("%s v1:%v\n", exampid, "starts_...")

	// TLS Configuration.  This configuration assumes that:
	// a) The server used does *not* require client certificates
	// b) This client has no need to authenticate the server
	// Note that the tls.Config structure can be modified to support any
	// authentication scheme, including two-way/mutual authentication. Examples
	// are provided elsewhere in this project.
	tc = new(tls.Config)
	tc.InsecureSkipVerify = true // Do *not* check the server's certificate

	// Open a net connection
	h, p := senv.HostAndPort()
	hap := net.JoinHostPort(h, p)
	// Use tls.Dial, not net.Dial
	n, e := tls.Dial("tcp", hap, tc)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v v2:%v\n", exampid, "dial complete ...", hap)

	// All stomp API methods require 'Headers'.  Stomp headers are key/value
	// pairs.  The stompngo package implements them using a string slice.
	ch := sngecomm.ConnectHeaders()

	// Get a stomp connection.  Parameters are:
	// a) the opened net connection
	// b) the Headers
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v\n", exampid, "stomp_connect_complete_...")

	// *NOTE* your application functionaltiy goes here!

	// Stomp disconnects are not required, but highly recommended.
	// We use empty headers here for the DISCONNECT.  This will trigger the
	// disconnect receipt described below.
	//
	// If you do *not* want a disconnect receipt use:
	// stompngo.Headers{"noreceipt", "true"}
	//
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalf("%s v1:%v v2:%v\n", exampid, "sngDisconnectTLS", e) // Handle this ......
	}
	ll.Printf("%s v1:%v\n", exampid, "stomp_disconnect_complete_...")

	// After a DISCONNECT there is (by default) a 'Disconnect
	// Receipt' available.  The disconnect receipt is an instance of
	// stompngo.MessageData.
	ll.Printf("%s v1:%v v2:%v\n", exampid, "disconnect receipt", conn.DisconnectReceipt)
	ll.Printf("%s v1:%v\n", exampid, "stomp_disconnect_complete_...")

	// Close the net connection.
	e = n.Close()
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v\n", exampid, "network_close_complete_...")

	ll.Printf("%s v1:%v\n", exampid, "ends_...")
}
