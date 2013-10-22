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

Prime the queue for this demonstration using publish.go.

	Examples:

		# Prime a queue with messages:
		STOMP_PORT=61613 STOMP_NMSGS=10 go run publish.go

		# Review ActiveMQ balancing characteristics:
		STOMP_PORT=61613 go run recv_mds.go

		# Prime a queue with messages again:
		STOMP_PORT=62613 STOMP_NMSGS=10 go run publish.go

		# Review Apollo balancing characteristics:
		STOMP_PORT=62613 go run recv_mds.go

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
	id := stompngo.Uuid() // Use package convenience function for unique ID
	d := sngecomm.Dest()
	var r <-chan stompngo.MessageData
	if ack {
		r = sngecomm.Subscribe(conn, d, id, "client-individual")
	} else {
		r = sngecomm.Subscribe(conn, d, id, "auto")
	}
	// Receive loop
	for {
		d := <-r // Read a messagedata struct
		if d.Error != nil {
			panic(d.Error)
		}
		m := d.Message.BodyString()
		fmt.Println(exampid, "subnumber", s, m, id)
		runtime.Gosched()
		time.Sleep(1 * time.Second)
		if ack {
			sngecomm.Ack(conn, d.Message.Headers, id)
			fmt.Println(exampid + "ACK complete ...")
		}
		runtime.Gosched()
	}
	wgrecv.Done() // Never get here, cancel via ^C
}

// Connect to a STOMP broker, receive and ack some messages.
// Disconnect never occurs, kill via ^C.
func main() {
	fmt.Println(exampid, "starts ...")

	// Set up the connection.
	h, port := sngecomm.HostAndPort() //
	n, e := net.Dial("tcp", net.JoinHostPort(h, port))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "dial complete ...")
	ch := sngecomm.ConnectHeaders()
	conn, e = stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid, "stomp connect complete ...", conn.Protocol())

	wgrecv.Add(ns) // Number of subscriptions, hard coded in this demonstartion
	for i := 1; i <= ns; i++ {
		go recv(i)
	}
	fmt.Println(exampid, "receivers started ...")

	wgrecv.Wait() // This will never complete, use ^C to cancel

}
