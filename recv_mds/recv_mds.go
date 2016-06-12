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

// Show a number of queue writers and readers operating concurrently.
// Try to be realistic about workloads.
// Receiver checks messages for proper queue and message number.
// All senders and receivers use the same Stomp connection.

/*
Outputs to demonstrate different broker's algorithms for balancing messages across
multiple subscriptions to the same queue.  (One go routine per subscription). In
this example all subscriptions share the same connection/session.  The actual
results can / will likely be slightly surprising.  YMMV.

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
	"log"
	"net"
	"runtime"
	"time"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "recv_mds: "
	ns      = 4                           // Number of subscriptions
	n       net.Conn                      // Network Connection
	conn    *stompngo.Connection          // Stomp Connection
	ackMode string               = "auto" // ackMode control
	port    string
)

func recv(s int) {
	log.Println(exampid, "receiver", s, "starts")
	// Setup Headers ...
	id := stompngo.Uuid() // Use package convenience function for unique ID
	d := sngecomm.Dest()
	ackMode = sngecomm.AckMode() // get ack mode

	pbc := sngecomm.Pbc() // Print byte count

	var r <-chan stompngo.MessageData
	r = sngecomm.Subscribe(conn, d, id, ackMode)
	// Receive loop
	mc := 0
	for {
		d := <-r // Read a messagedata struct
		mc++
		if d.Error != nil {
			panic(d.Error)
		}
		log.Println(exampid, "subnumber", s, id, mc)
		if pbc > 0 {
			maxlen := pbc
			if len(d.Message.Body) < maxlen {
				maxlen = len(d.Message.Body)
			}
			ss := string(d.Message.Body[0:maxlen])
			log.Printf("Payload: %s\n", ss) // Data payload
		}

		time.Sleep(3 * time.Second)
		runtime.Gosched()
		if ackMode != "auto" {
			sngecomm.Ack(conn, d.Message.Headers, id)
			log.Println(exampid + "ACK complete ...")
		}
		runtime.Gosched()
	}
}

// Connect to a STOMP broker, receive and ackMode some messages.
// Disconnect never occurs, kill via ^C.
func main() {
	log.Println(exampid, "starts ...")

	// Set up the connection.
	h, port := sngecomm.HostAndPort() //
	n, e := net.Dial("tcp", net.JoinHostPort(h, port))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid, "dial complete ...")
	ch := sngecomm.ConnectHeaders()
	conn, e = stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid, "stomp connect complete ...", conn.Protocol())

	for i := 1; i <= ns; i++ {
		go recv(i)
	}
	log.Println(exampid, ns, "receivers started ...")

	select {} // This will never complete, use ^C to cancel

}
