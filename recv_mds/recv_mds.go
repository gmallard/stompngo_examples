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
		# Review ActiveMQ balancing characteristics.  Note:
		# this will eventually block, and the program will have to be
		# forcibly stopped.
		STOMP_PORT=61613 STOMP_ACKMODE="client-individual" go run recv_mds.go

		# Prime a queue with messages again:
		STOMP_PORT=62613 STOMP_NMSGS=10 go run publish.go
		# Review Apollo balancing characteristics.  Note:
		# this will eventually block, and the program will have to be
		# forcibly stopped.
		STOMP_PORT=62613 STOMP_ACKMODE="client-individual" go run recv_mds.go

*/
package main

import (
	"log"
	"net"
	"os"
	"runtime"
	"time"
	//
	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "recv_mds: "
	ns      = 4                           // Number of subscriptions
	n       net.Conn                      // Network Connection
	conn    *stompngo.Connection          // Stomp Connection
	ackMode string               = "auto" // ackMode control
	port    string
	ll      = log.New(os.Stdout, "EMDS ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

func recv(conn *stompngo.Connection, s int) {
	ll.Println(exampid, "receiver", s, "starts")
	// Setup Headers ...
	id := stompngo.Uuid() // Use package convenience function for unique ID
	d := senv.Dest()
	ackMode = sngecomm.AckMode() // get ack mode

	pbc := sngecomm.Pbc() // Print byte count

	sc := sngecomm.HandleSubscribe(conn, d, id, ackMode)
	// Receive loop.
	mc := 0
	var md stompngo.MessageData
	for {
		select {
		case md = <-sc: // Read a messagedata struct, with a MESSAGE frame
		case md = <-conn.MessageData: // Read a messagedata struct, with a ERROR/RECEIPT frame
			// Unexpected here in this example.
			ll.Fatalln(exampid, md) // Handle this
		}
		//
		mc++
		if md.Error != nil {
			panic(md.Error)
		}
		ll.Println(exampid, "subnumber", s, id, mc)
		if pbc > 0 {
			maxlen := pbc
			if len(md.Message.Body) < maxlen {
				maxlen = len(md.Message.Body)
			}
			ss := string(md.Message.Body[0:maxlen])
			ll.Printf("Payload: %s\n", ss) // Data payload
		}

		// time.Sleep(3 * time.Second) // A very arbitrary number
		// time.Sleep(500 * time.Millisecond) // A very arbitrary number
		runtime.Gosched()
		time.Sleep(1500 * time.Millisecond) // A very arbitrary number
		runtime.Gosched()
		if ackMode != "auto" {
			sngecomm.HandleAck(conn, md.Message.Headers, id)
			ll.Println(exampid + "ACK complete ...")
		}
		runtime.Gosched()
	}
}

// Connect to a STOMP broker, receive and ackMode some messages.
// Disconnect never occurs, kill via ^C.
func main() {
	ll.Println(exampid, "starts ...")

	// Set up the connection.
	h, p := senv.HostAndPort() //
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Println(exampid, "dial complete ...", net.JoinHostPort(h, p))
	ch := sngecomm.ConnectHeaders()
	conn, e = stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Println(exampid, "stomp connect complete ...", conn.Protocol())

	for i := 1; i <= ns; i++ {
		go recv(conn, i)
	}
	ll.Println(exampid, ns, "receivers started ...")

	select {} // This will never complete, use ^C to cancel

}
