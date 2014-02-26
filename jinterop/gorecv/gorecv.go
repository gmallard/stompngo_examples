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
A message receiver, to demonstrate JMS interoperability.
*/
package main

import (
	"fmt"
	"log"
	"net"
	//
	"github.com/gmallard/stompngo")

var exampid = "gorecv: "

var nmsgs = 5

// Connect to a STOMP 1.1 broker, receive some messages and disconnect.
func main() {
	fmt.Println(exampid + "starts ...")

	// Set up the connection.
	n, e := net.Dial("tcp", "localhost:61613")
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "dial complete ...")
	eh := stompngo.Headers{"login", "userr", "passcode", "passw0rd"}
	conn, e := stompngo.Connect(n, eh)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "stomp connect complete ...")

	// Setup Headers ...
	u := stompngo.Uuid() // Use package convenience function for unique ID
	s := stompngo.Headers{"destination", "/queue/allards.queue",
		"id", u} // subscribe/unsubscribe headers

	// Subscribe
	r, e := conn.Subscribe(s)
	if e != nil {
		log.Fatalln(e) // Handle this ...
	}
	fmt.Println(exampid + "stomp subscribe complete ...")
	// Read data from the returned channel
	for i := 1; i <= nmsgs; i++ {
		m := <-r
		fmt.Println(exampid + "channel read complete ...")
		// MessageData has two components:
		// a) a Message struct
		// b) an Error value.  Check the error value as usual
		if m.Error != nil {
			log.Fatalln(m.Error) // Handle this
		}
		//
		fmt.Printf("Frame Type: %s\n", m.Message.Command) // Will be MESSAGE or ERROR!
		h := m.Message.Headers
		for j := 0; j < len(h)-1; j += 2 {
			fmt.Printf("Header: %s:%s\n", h[j], h[j+1])
		}
		fmt.Printf("Payload: %s\n", string(m.Message.Body)) // Data payload
	}
	// It is polite to unsubscribe, although unnecessary if a disconnect follows.
	// With Stomp 1.1, the same unique ID is required on UNSUBSCRIBE.  Failure
	// to provide it will result in an error return.
	e = conn.Unsubscribe(s)
	if e != nil {
		log.Fatalln(e) // Handle this ...
	}
	fmt.Println(exampid + "stomp unsubscribe complete ...")

	// Disconnect from the Stomp server
	eh = stompngo.Headers{}
	e = conn.Disconnect(eh)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "stomp disconnect complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "network close complete ...")

	fmt.Println(exampid + "ends ...")
}
