//
// Copyright © 2015 Guy M. Allard
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
sngissue25 getters
*/
package main

import (
	"fmt"
	"log"
	"os"
	"net"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var exampid = "getters: "

var (
	qbase = "/queue/sngissue25."
	conn  *stompngo.Connection
)

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
	conn, e = stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid+"stomp connect complete ...", conn.Protocol())
	fmt.Println(exampid+"connected headers", conn.ConnectResponse.Headers)
	// *NOTE* your application functionaltiy goes here!

	runApp(conn)

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

func runApp(c *stompngo.Connection) {

	// Handle each queue in turn ....
	for cq := 1; cq <= sngecomm.Nqs(); cq++ {
		//	for cq := 1; cq <= 1; cq++ { // TEMP
		cqs := fmt.Sprintf("%d", cq)
		nqn := qbase + cqs // Next queue name
		fmt.Println(exampid+"next queue name ...", nqn)
//		id := stompngo.Uuid()
		id := "Q" + cqs
		fmt.Println(exampid+"next queue id ...", id)
		// Subscribe
		sh := stompngo.Headers{"destination", nqn,
			"id", id,
			"ack", "client"}
		if os.Getenv("STOMP_DRAIN") != "" {
			sh = sh.Add("subdrain", "yep")
		}
		mdc, e := c.Subscribe(sh)
		if e != nil {
			log.Fatalln(exampid+"fatal subscribe", e)
		}

		// OK, now get *some* messages from this queue, issue an ACK, and then
		// UNSUBSCRIBE.  From each queue get *only* this number of messages:
		// queue 1, get 1 message
		// queue 2, get 2 messages
		// queue 3, get 3 messages
		// .... and so on.

		rmc := 0
		var lmd stompngo.MessageData

		for mn := 1; ; mn++ { // Only get queue number of MessageData elements
			select {
			case lmd = <-mdc: // Get the next MessageData for this queue
				//
			case lmd = <-c.MessageData: // Possible ERROR frame
				//
			}
			if lmd.Error != nil {
				log.Fatalln(exampid+"receive error", cq, mn)
			}
			fmt.Println(exampid+"received message number", mn, "->")
			fmt.Println(exampid+"received COMMAND", lmd.Message.Command)
			fmt.Println(exampid+"received HEADERS", lmd.Message.Headers)
			fmt.Println(exampid+"received BODY", string(lmd.Message.Body))
			rmc++
			if lmd.Message.Command == "ERROR" {
				log.Fatalln("fatal ERROR", lmd)
			}
			if rmc == cq {
				break
			}
		} // End of message processing for the current queue

		fmt.Println(exampid+"received message count:", rmc)
		// The ACK of the last message
		fmt.Println(exampid+"starting ACK, queue number", cq, "headers", lmd.Message.Headers, "dest", nqn)
		sngecomm.Ack(c, lmd.Message.Headers, id)
		fmt.Println(exampid+"ending ACK, queue number", cq, "headers", lmd.Message.Headers, "dest", nqn)
		// The UNSUBSCRIBE
		fmt.Println(exampid+"start UNSUBSCRIBE, queue:", nqn)
		sngecomm.Unsubscribe(c, nqn, id)
		fmt.Println(exampid+"end UNSUBSCRIBE, queue:", nqn)
	} // End of queue scan
}
