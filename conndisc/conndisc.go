//
// Copyright © 2013-2015 Guy M. Allard
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
Connect and Disconnect from a STOMP broker using TCP.

	Examples:

		# Connect to a broker with all defaults:
		# Host is "localhost"
		# Port is 61613
		# Login is "guest"
		# Passcode is "guest
		# Virtual Host is "localhost"
		# Protocol is 1.1
		go run conndisc.go

		# Connect to a broker using STOMP protocol level 1.0:
		STOMP_PROTOCOL=1.0 go run conndisc.go

		# Connect to a broker using a custom host and port:
		STOMP_HOST=tjjackson STOMP_PORT=62613 go run conndisc.go

		# Connect to a broker using a custom port and virtual host:
		STOMP_PORT=41613 STOMP_VHOST="/" go run conndisc.go

		# Connect to a broker using a custom login and passcode:
		STOMP_LOGIN="userid" STOMP_PASSCODE="t0ps3cr3t" go run conndisc.go

*/
package main

import (
	"fmt"
	"log"
	"net"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var exampid = "conndisc: "

// Connect to a STOMP broker and disconnect.
func main() {
	fmt.Println(exampid + "starts ...")

	// Open a net connection
	h, p := sngecomm.HostAndPort()
  jhp := net.JoinHostPort(h, p)
	fmt.Println(exampid + "will dial ...", jhp)
	n, e := net.Dial("tcp", jhp)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "dial complete ...")

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
	fmt.Println(exampid+"stomp connect complete, protocol level is:", conn.Protocol())

	// *NOTE* your application functionaltiy goes here!

	// Polite Stomp disconnects are not required, but highly recommended.
	// Empty headers here for the DISCONNECT
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "stomp disconnect complete ...")

	// Close the net connection.
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "network close complete ...")

	fmt.Println(exampid + "ends ...")
}
