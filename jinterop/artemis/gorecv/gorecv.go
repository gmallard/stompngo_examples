//
// Copyright © 2011-2016 Guy M. Allard
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
A message receiver, to demonstrate JMS interoperability.
*/
package main

import (
	"log"
	"net"
	"os"
	//
	"github.com/gmallard/stompngo"
)

var (
	exampid = "gorecv: "
	ll      = log.New(os.Stdout, "GOJRCVRT ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	nmsgs   = 1
)

// Connect to a STOMP 1.2 broker, receive some messages and disconnect.
func main() {
	ll.Printf("%s v1:%v\n", exampid, "starts_...")

	// Set up the connection.
	n, e := net.Dial("tcp", "localhost:31613")
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v\n", exampid, "dial_complete_...")
	ch := stompngo.Headers{"login", "userr", "passcode", "passw0rd",
		"host", "localhost", "accept-version", "1.2"}
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v\n", exampid, "stomp_connect_complete_...")

	// Setup Headers ...
	id := stompngo.Uuid() // Use package convenience function for unique ID
	sbh := stompngo.Headers{"destination", "jms.queue.exampleQueue",
		"id", id} // subscribe/unsubscribe headers

	// Subscribe
	sc, e := conn.Subscribe(sbh)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ...
	}
	ll.Printf("%s v1:%v\n", exampid, "stomp_subscribe_complete_...")

	var md stompngo.MessageData
	// Read data from the returned channel
	for i := 1; i <= nmsgs; i++ {
		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// A RECEIPT or ERROR frame is unexpected here
			ll.Fatalf("%s v1:%v\n", exampid, md) // Handle this
		}
		ll.Printf("%s v1:%v\n", exampid, "channel_read_complete_...")
		// MessageData has two components:
		// a) a Message struct
		// b) an Error value.  Check the error value as usual
		if md.Error != nil {
			ll.Fatalf("%s  f4v:%v\n", exampid, md.Error) // Handle this
		}
		//
		ll.Printf("Frame Type: %s\n", md.Message.Command) // Should be MESSAGE
		wh := md.Message.Headers
		for j := 0; j < len(wh)-1; j += 2 {
			ll.Printf("Header: %s:%s\n", wh[j], wh[j+1])
		}
		ll.Printf("Payload: %s\n", string(md.Message.Body)) // Data payload
	}
	// It is polite to unsubscribe, although unnecessary if a disconnect follows.
	// With Stomp 1.1+, the same unique ID is required on UNSUBSCRIBE.  Failure
	// to provide it will result in an error return.
	e = conn.Unsubscribe(sbh) // Same headers as Subscribe
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ...
	}
	ll.Printf("%s v1:%v\n", exampid, "stomp_unsubscribe_complete_...")

	// Disconnect from the Stomp server
	dh := stompngo.Headers{}
	e = conn.Disconnect(dh)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v\n", exampid, "stomp_disconnect_complete_...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v\n", exampid, "network_close_complete_...")

	ll.Printf("%s v1:%v\n", exampid, "ends_...")
}
