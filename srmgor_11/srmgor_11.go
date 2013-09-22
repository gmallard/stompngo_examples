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

/*
Send and receive many STOMP messages using multiple queues and goroutines
to service each send or receive instance.  All senders share a single
STOMP connection, as do all receivers.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
	"log"
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"time"
)

var exampid = "srmgor_11:"

var wgsend sync.WaitGroup
var wgrecv sync.WaitGroup
var wgall sync.WaitGroup

// We 'stagger' between each message send and message receive for a random
// amount of time.
// Vary these for experimental purposes.  YMMV.
var max int64 = 1e9      // Max stagger time (nanoseconds)
var min int64 = max / 10 // Min stagger time (nanoseconds)
// Vary these for experimental purposes.  YMMV.
var send_factor int64 = 1 // Send factor time
var recv_factor int64 = 1 // Receive factor time

// Wait flags
var send_wait = true
var recv_wait = true

// Possible profile file
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// Send messages to a particular queue
func sender(conn *stompngo.Connection, qn, c int) {
	qns := fmt.Sprintf("%d", qn) // queue number
	fmt.Println(exampid, "send starts", qn)
	//
	qp := sngecomm.Dest() // queue name prefix
	q := qp + "." + qns
	fmt.Println(exampid, "send queue name", q, qn)
	h := stompngo.Headers{"destination", q} // send Headers
	if sngecomm.Persistent() {
		h = h.Add("persistent", "true")
	}
	// Send loop
	for i := 1; i <= c; i++ {
		si := fmt.Sprintf("%d", i)
		// Generate a message to send ...............
		mp := exampid + "|" + "payload" + "|qnum:" + qns + "|msgnum:" + si + " :"
		m := mp + sngecomm.Partial() // Variable length
		fmt.Println(exampid, "send message", mp, qn)
		e := conn.Send(h, m)
		if e != nil {
			log.Fatalln(exampid, "send error", e, qn)
			break
		}
		if send_wait {
			runtime.Gosched()                                                              // yield for this example
			time.Sleep(time.Duration(send_factor * (sngecomm.ValueBetween(min, max) / 2))) // Time to build next message
		}
	}
	// Sending is done
	fmt.Println(exampid, "send ends", qn)
	wgsend.Done()
}

// Asynchronously process all messages for a given subscription.
func receiveWorker(mc <-chan stompngo.MessageData, qns string, c int, d chan<- bool) {
	// Receive loop
	for i := 1; i <= c; i++ {
		d := <-mc
		if d.Error != nil {
			log.Fatalln(exampid, "recv read error", d.Error, qns)
		}

		// Process the inbound message .................
		m := d.Message.BodyString()
		li := strings.LastIndex(m, ":")
		fmt.Println(exampid, "recv message", string(m[0:li]), qns)

		// Sanity check the queue and message numbers
		mns := fmt.Sprintf("%d", i) // message number
		t := "|qnum:" + qns + "|msgnum:" + mns
		if !strings.Contains(m, t) {
			log.Fatalln(exampid, "recv bad message", m, t, qns)
		}
		if recv_wait {
			runtime.Gosched()                                                              // yield for this example
			time.Sleep(time.Duration(send_factor * (sngecomm.ValueBetween(min, max) / 2))) // Time to build next message
		}
	}
	//
	d <- true
}

// Receive messages from a particular queue
func receiver(conn *stompngo.Connection, qn, c int) {
	qns := fmt.Sprintf("%d", qn) // queue number
	fmt.Println(exampid, "recv starts", qn)
	//
	qp := sngecomm.Dest() // queue name prefix
	q := qp + "." + qns
	fmt.Println(exampid, "recv queue name", q, qn)
	u := stompngo.Uuid() // A unique subscription ID
	h := stompngo.Headers{"destination", q, "id", u}
	// Subscribe
	r, e := conn.Subscribe(h)
	if e != nil {
		log.Fatalln(exampid, "recv subscribe error", e, qn)
	}

	// Many receivers running under the same connection can cause
	// (wire read) performance issues.  This is *very* dependent on the broker
	// being used, specifically the broker's algorithm for putting messages on
	// the wire.
	// To alleviate those issues, this strategy insures that messages are
	// received from the wire as soon as possible.  Those messages are then
	// buffered internally for (possibly later) application processing.

	// Process all inputs async .......
	mc := make(chan stompngo.MessageData, c) // All messages we could possibly get, YMMV
	dc := make(chan bool)                    // Receive processing done channel
	go receiveWorker(mc, qns, c, dc)         // Start async processor
	for i := 1; i <= c; i++ {
		mc <- <-r // Receive message data as soon as possible, and internally queue it
	}
	<-dc // Wait until receive processing is done

	// Unsubscribe
	e = conn.Unsubscribe(h)
	if e != nil {
		log.Fatalln(exampid, "recv unsubscribe error", e, qn)
	}
	// Receiving is done
	fmt.Println(exampid, "recv ends", qn)
	wgrecv.Done()
}

func startSenders(qn int) {
	fmt.Println(exampid, "startSenders starts", qn)

	// Open
	h, p := sngecomm.HostAndPort11() // a 1.1 connect
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(exampid, "startSenders netconnect error", e, qn) // Handle this ......
	}

	// Stomp connect, 1.1
	ch := stompngo.Headers{"host", sngecomm.Vhost(), "accept-version", "1.1"}
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(exampid, "startSenders stompconnect error", e, qn) // Handle this ......
	}
	log.Println(exampid, "startSenders connection", conn, qn)
	c := sngecomm.Nmsgs() // message count
	fmt.Println(exampid, "startSenders message count", c, qn)
	for i := 1; i <= qn; i++ { // all queues
		wgsend.Add(1)
		go sender(conn, i, c)
	}
	wgsend.Wait()

	// Disconnect from Stomp server
	eh := stompngo.Headers{}
	e = conn.Disconnect(eh)
	if e != nil {
		log.Println(exampid, "startSenders disconnect error", e, qn) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		log.Fatalln(exampid, "startSenders netclose error", e, qn) // Handle this ......
	}

	fmt.Println(exampid, "startSenders ends", qn)
	sngecomm.ShowStats(exampid, "startSenders", conn)
	wgall.Done()
}

func startReceivers(qn int) {
	fmt.Println(exampid, "startReceivers starts", qn)

	// Open
	h, p := sngecomm.HostAndPort11() // a 1.1 connect
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(exampid, "startReceivers nectonnr:", e, qn) // Handle this ......
	}
	ch := stompngo.Headers{"host", sngecomm.Vhost(), "accept-version", "1.1"}
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln("startReceivers stompconnectr:", e, qn) // Handle this ......
	}
	log.Println("startReceivers Receive connection is:", conn, qn)
	c := sngecomm.Nmsgs() // get message count
	fmt.Println(exampid, "startReceivers message count", c, qn)
	for i := 1; i <= qn; i++ { // all queues
		wgrecv.Add(1)
		go receiver(conn, i, c)
	}
	wgrecv.Wait()

	// Disconnect from Stomp server
	eh := stompngo.Headers{}
	e = conn.Disconnect(eh)
	if e != nil {
		log.Println(exampid, "startReceivers disconnect error", e, qn) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		log.Println(exampid, "startReceivers netclose error", e, qn) // Handle this ......
	}

	fmt.Println(exampid, "startReceivers ends", qn)
	sngecomm.ShowStats(exampid, "startReceivers", conn)
	wgall.Done()
}

// Show a number of writers and readers operating concurrently from unique
// destinations.
func main() {
	fmt.Println(exampid, "main starts")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if sngecomm.SetMAXPROCS() {
		nc := runtime.NumCPU()
		fmt.Println(exampid, "main number of CPUs is:", nc)
		c := runtime.GOMAXPROCS(nc)
		fmt.Println(exampid, "main previous number of GOMAXPROCS is:", c)
		fmt.Println(exampid, "main current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	}
	//
	send_wait = sngecomm.SendWait()
	recv_wait = sngecomm.RecvWait()
	//
	q := sngecomm.Nqs()
	fmt.Println(exampid, "main Nqs:", q)
	//
	wgall.Add(2)
	go startReceivers(q)
	go startSenders(q)
	wgall.Wait()

	fmt.Println(exampid, "main ends")
}
