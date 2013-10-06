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

// Show a number of queue writers and readers operating concurrently.
// Try to be realistic about workloads.
// Receiver checks messages for proper queue and message number.
// All senders and receivers use the same Stomp connection.

/*
Output can demonstrate different broker's algorithms for balancing messages across
multiple subscriptions to the same queue.  In this example all subscriptions
share the same connection/session.
*/
package main

import (
	"fmt"
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

var (
	exampid = "recv_mds: "
	ns      = 4 // Number of subscriptions
	wgrecv  sync.WaitGroup
	n       net.Conn                     // Network Connection
	conn    *stompngo.Connection         // Stomp Connection
	ack     bool                 = false // ack mode control
	port    string
)

func recv(s int) {
	// Setup Headers ...
	u := stompngo.Uuid() // Use package convenience function for unique ID
	h := stompngo.Headers{"destination", sngecomm.Dest(),
		"id", u}
	//
	if ack {
		h = h.Add("ack", "client-individual")
		if port == "61613" { // AMQ
			fmt.Println(exampid, "AMQ Port detected")
			h = h.Add("activemq.prefetchSize", "1")
		}
	} else {
		h = h.Add("ack", "auto")
	}
	fmt.Println(exampid, s, h)
	// Subscribe
	r, e := conn.Subscribe(h)
	if e != nil {
		log.Fatalln(exampid, "recv subscribe error:", e, s)
	}
	// Receive loop
	for {
		d := <-r // Read a messagedata struct
		if d.Error != nil {
			panic(d.Error)
		}
		m := d.Message.BodyString()
		fmt.Println(exampid, s, m, u)
		runtime.Gosched()
		time.Sleep(1 * time.Second)
		if ack {
			ah := stompngo.Headers{"message-id", d.Message.Headers.Value("message-id"),
				"subscription", u} // 1.1 ACK headers
			// fmt.Println(exampid, "ACK Headers", ah)
			e := conn.Ack(ah)
			if e != nil {
				log.Fatalln(e) // Handle this
			}
		}
		runtime.Gosched()
	}
	wgrecv.Done() // Never get here, cancel via ^C
}

// Connect to a STOMP 1.1 broker, receive some messages and disconnect.
func main() {
	fmt.Println(exampid, "starts ...")

	// Set up the connection.
	h, port := sngecomm.HostAndPort11() //
	n, e := net.Dial("tcp", net.JoinHostPort(h, port))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "dial complete ...")
	ch := stompngo.Headers{"accept-version", sngecomm.Protocol(),
		"host", sngecomm.Vhost()}
	conn, e = stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "stomp connect complete ...", conn.Protocol())

	wgrecv.Add(ns)
	for i := 1; i <= ns; i++ {
		go recv(i)
	}
	fmt.Println(exampid, "receivers started ...")

	wgrecv.Wait() // This will never complete, use ^C to cancel

	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "stomp disconnect complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "network close complete ...")

	fmt.Println(exampid, "ends ...")
}
