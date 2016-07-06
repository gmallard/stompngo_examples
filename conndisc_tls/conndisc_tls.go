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
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid    = "conndisc_tls: "
	testConfig *tls.Config
)

// Connect to a STOMP 1.0 broker using TLS and disconnect.
func main() {
	log.Println(exampid + "starts ...")

	// TLS Configuration.  This configuration assumes that:
	// a) The server used does *not* require client certificates
	// b) This client has no need to authenticate the server
	// Note that the tls.Config structure can be modified to support any
	// authentication scheme, including two-way/mutual authentication. Examples
	// are provided elsewhere in this project.
	testConfig = new(tls.Config)
	testConfig.InsecureSkipVerify = true // Do *not* check the server's certificate

	// Open a net connection
	h, p := sngecomm.HostAndPort()
	// Use tls.Dial, not net.Dial
	n, e := tls.Dial("tcp", net.JoinHostPort(h, p), testConfig)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "dial complete ...")

	// All stomp API methods require 'Headers'.  Stomp headers are key/value
	// pairs.  The stompngo package implements them using a string slice.
	ch := sngecomm.ConnectHeaders()

	// Get a stomp connection.  Parameters are:
	// a) the opened net connection
	// b) the Headers
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "stomp connect complete ...")

	// *NOTE* your application functionaltiy goes here!

	// Polite Stomp disconnects are not required, but highly recommended.
	// Empty headers for this DISCONNECT.
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "stomp disconnect complete ...")

	// Close the net connection.
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "network close complete ...")

	log.Println(exampid + "ends ...")
}
