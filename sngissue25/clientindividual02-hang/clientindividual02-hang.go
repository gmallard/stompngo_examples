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

See run.sh in this directory for usage example.

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

var exampid = "subscribe: "

// Connect to a STOMP broker, subscribe and receive some messages and disconnect.
func main() {
	fmt.Println(exampid + "starts ...")

	// Set up the connection.
	h, p := sngecomm.HostAndPort()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "dial complete ...")
	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid+"stomp connect complete ...", conn.Protocol())

	fmt.Println(exampid+"connected headers", conn.ConnectResponse.Headers)
	// *NOTE* your application functionaltiy goes here!
	// With Stomp, you must SUBSCRIBE to a destination in order to receive.
	// Subscribe returns a channel of MessageData struct.
	// Here we use a common utility routine to handle the differing subscribe
	// requirements of each protocol level.
	d := sngecomm.Dest()
	id := stompngo.Uuid()
	r := sngecomm.Subscribe(conn, d, id, "client-individual")
	fmt.Println(exampid + "stomp subscribe complete ...")
	// Read data from the returned channel
	for i := 1; i <= sngecomm.Nmsgs(); i++ {
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
		if m.Message.Command != stompngo.MESSAGE {
			log.Fatalln(m) // Handle this ...
		}
		h := m.Message.Headers
		for j := 0; j < len(h)-1; j += 2 {
			fmt.Printf("Header: %s:%s\n", h[j], h[j+1])
		}
		fmt.Printf("Payload: %s\n", string(m.Message.Body)) // Data payload
		// Issue the final ACK
		sngecomm.Ack(conn, m.Message.Headers, id)
		fmt.Println("Ending ACK", i)
	}

	// Do *NOT* UNSUBSCRIBE, proceed directly to DISCONNECT

	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
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
