//
// Copyright © 2011-2013 Guy M. Allard
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
Connect and Disconnect from a STOMP 1.1 broker.
*/
package main

import (
	"fmt"
	"github.com/gmallard/stompngo"
	. "github.com/gmallard/stompngo_examples/sngecomm"
	"log"
	"net"
)

var exampid = "connect_11: "

// Connect to a STOMP 1.1 broker and disconnect.
func main() {
	fmt.Println(exampid + "starts ...")

	// Open a net connection
	h, p := HostAndPort11()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "dial complete ...")

	// All stomp API methods require 'Headers'.  Stomp headers are key/value
	// pairs.  The stompngo package implements them using a string slice.
	//
	// To connect to a Stomp 1.1 broker, you must:
	// a) Use the correct host and port of course
	// b) Pass the 'accept-version' header per specification requirements
	// c) Pass the 'host' header per specification requirements
	//
	// We demand a 1.1 connection here.
	// Note that the 1.1 vhost _could_ be different than the host name used
	// for the connection, but in this example is the same.
	//
	ch := stompngo.Headers{"accept-version", "1.1",
		"host", Vhost()}

	// Get a stomp connection.  Parameters are:
	// a) the opened net connection
	// b) the (empty) Headers
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "stomp connect complete ...")

	// *NOTE* your application functionaltiy goes here!

	// Polite Stomp disconnects are not required, but highly recommended.
	// Empty headers again in this example.
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
