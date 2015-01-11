//
// Copyright © 2011-2015 Guy M. Allard
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
	"log"
	"net"
	"runtime"
	"sync"
	"time"
	//
	"github.com/davecheney/profile"
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var exampid = "srmgor_2conn:"

var wgsend sync.WaitGroup
var wgrecv sync.WaitGroup
var wgall sync.WaitGroup

// We 'stagger' between each message send and message receive for a random
// amount of time.
// Vary these for experimental purposes.  YMMV.
var max int64 = 1e9      // Max stagger time (nanoseconds)
var min int64 = max / 10 // Min stagger time (nanoseconds)

// Wait flags
var send_wait = true
var recv_wait = true

// Sleep multipliers
var sendFact float64 = 1.0
var recvFact float64 = 1.0

// Possible profile file
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

// Send messages to a particular queue
func sender(conn *stompngo.Connection, qn, c int) {
	qns := fmt.Sprintf("%d", qn) // queue number
	fmt.Println(sngecomm.ExampIdNow(exampid), "send starts", qn)
	//
	qp := sngecomm.Dest() // queue name prefix
	q := qp + "." + qns
	fmt.Println(sngecomm.ExampIdNow(exampid), "send queue name", q)
	h := stompngo.Headers{"destination", q,
		"qnum", qns} // send Headers
	if sngecomm.Persistent() {
		h = h.Add("persistent", "true")
	}
	//
	tmr := time.NewTimer(100 * time.Hour)
	// Send loop
	for i := 1; i <= c; i++ {
		si := fmt.Sprintf("%d", i)
		sh := append(h, "msgnum", si)
		// Generate a message to send ...............
		fmt.Println(sngecomm.ExampIdNow(exampid), "send message", qns, si)
		e := conn.Send(sh, string(sngecomm.Partial()))

		if e != nil {
			log.Fatalln(sngecomm.ExampIdNow(exampid), "send error", e, qn)
			break
		}
		if send_wait {
			runtime.Gosched() // yield for this example
			d := time.Duration(sngecomm.ValueBetween(min, max, sendFact))
			fmt.Println(sngecomm.ExampIdNow(exampid), "send", "stagger", int64(d)/1000000, "ms")
			tmr.Reset(d)
			_ = <-tmr.C
		}
	}
	// Sending is done
	fmt.Println(sngecomm.ExampIdNow(exampid), "send ends", qn)
	wgsend.Done()
}

// Asynchronously process all messages for a given subscription.
func receiveWorker(mc <-chan stompngo.MessageData, qns string, c int,
	d chan<- bool, conn *stompngo.Connection, id string) {
	//
	tmr := time.NewTimer(100 * time.Hour)
	// Receive loop
	for i := 1; i <= c; i++ {
		d := <-mc
		if d.Error != nil {
			log.Fatalln(sngecomm.ExampIdNow(exampid), "recv read error", d.Error, qns)
		}

		// Sanity check the queue and message numbers
		mns := fmt.Sprintf("%d", i) // message number
		if !d.Message.Headers.ContainsKV("qnum", qns) || !d.Message.Headers.ContainsKV("msgnum", mns) {
			log.Fatalln("Bad Headers", d.Message.Headers, qns, mns)
		}

		// Process the inbound message .................
		sl := 16
		if len(d.Message.Body) < sl {
			sl = len(d.Message.Body)
		}
		// Handle ACKs if needed
		if sngecomm.AckMode() != "auto" {
			ah := []string{}
			switch conn.Protocol() {
			case stompngo.SPL_11:
				ah = append(ah, "subscription", id, "message-id", d.Message.Headers.Value("message-id"))
			default: // 1.2 (NB: 1.0 not supported here)
				ah = append(ah, "id", d.Message.Headers.Value("ack"))
			}
			ah = append(ah, "qnum", qns, "msgnum", mns) // For tracking
			e := conn.Ack(ah)
			if e != nil {
				log.Fatalln("ACK Error", e)
			}
		}
		fmt.Println(sngecomm.ExampIdNow(exampid), "recv message", string(d.Message.Body[0:sl]), qns, d.Message.Headers.Value("msgnum"))
		if i == c {
			break
		}
		if recv_wait {
			runtime.Gosched() // yield for this example
			d := time.Duration(sngecomm.ValueBetween(min, max, recvFact))
			fmt.Println(sngecomm.ExampIdNow(exampid), "recv", "stagger", int64(d)/1000000, "ms")
			tmr.Reset(d)
			_ = <-tmr.C
		}
	}
	//
	d <- true
}

// Receive messages from a particular queue
func receiver(conn *stompngo.Connection, qn, c int) {
	qns := fmt.Sprintf("%d", qn) // queue number
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv starts", qn)
	//
	qp := sngecomm.Dest() // queue name prefix
	q := qp + "." + qns
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv queue name", q, qn)
	id := stompngo.Uuid() // A unique subscription ID
	r := sngecomm.Subscribe(conn, q, id, sngecomm.AckMode())
	// Many receivers running under the same connection can cause
	// (wire read) performance issues.  This is *very* dependent on the broker
	// being used, specifically the broker's algorithm for putting messages on
	// the wire.
	// To alleviate those issues, this strategy insures that messages are
	// received from the wire as soon as possible.  Those messages are then
	// buffered internally for (possibly later) application processing.

	// Process all inputs async .......
	var mc chan stompngo.MessageData
	nb := c // Buffer size
	if sngecomm.Conn2Buffer() > 0 {
		nb = sngecomm.Conn2Buffer() // User spec'd bufsize
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv", "mdbuffersize", nb, qns)
	mc = make(chan stompngo.MessageData, nb)   // MessageData Buffer size
	dc := make(chan bool)                      // Receive processing done channel
	go receiveWorker(mc, qns, c, dc, conn, id) // Start async processor
	for i := 1; i <= c; i++ {
		mc <- <-r // Receive message data as soon as possible, and internally queue it
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv", "waitforWorkersBegin", qns)
	<-dc // Wait until receive processing is done for this queue
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv", "waitforWorkersEnd", qns)

	// Unsubscribe
	sngecomm.Unsubscribe(conn, q, id)

	// Receiving is done
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv ends", qns)
	wgrecv.Done()
}

func startSenders(qn int) {
	fmt.Println(sngecomm.ExampIdNow(exampid), "startSenders starts", qn)

	// Open
	h, p := sngecomm.HostAndPort() // host and port
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "startSenders netconnect error", e, qn) // Handle this ......
	}

	// Stomp connect
	ch := sngecomm.ConnectHeaders()
	log.Println(sngecomm.ExampIdNow(exampid), "startSenders", "vhost:", sngecomm.Vhost(), "protocol:", sngecomm.Protocol())
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "startSenders stompconnect error", e, qn) // Handle this ......
	}
	log.Println(sngecomm.ExampIdNow(exampid), "startSenders connection", conn, qn)
	c := sngecomm.Nmsgs() // message count
	fmt.Println(sngecomm.ExampIdNow(exampid), "startSenders message count", c, qn)
	for i := 1; i <= qn; i++ { // all queues
		wgsend.Add(1)
		go sender(conn, i, c)
	}
	wgsend.Wait()

	// Disconnect from Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Println(sngecomm.ExampIdNow(exampid), "startSenders disconnect error", e, qn) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "startSenders netclose error", e, qn) // Handle this ......
	}

	fmt.Println(sngecomm.ExampIdNow(exampid), "startSenders ends", qn)
	sngecomm.ShowStats(exampid, "startSenders", conn)
	wgall.Done()
}

func startReceivers(qn int) {
	fmt.Println(sngecomm.ExampIdNow(exampid), "startReceivers starts", qn)

	// Open
	h, p := sngecomm.HostAndPort() // host and port
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "startReceivers nectonnr:", e, qn) // Handle this ......
	}
	ch := sngecomm.ConnectHeaders()
	log.Println(sngecomm.ExampIdNow(exampid), "startReceivers", "vhost:", sngecomm.Vhost(), "protocol:", sngecomm.Protocol())
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "startReceivers stompconnectr:", e, qn) // Handle this ......
	}
	log.Println("startReceivers Receive connection is:", conn, qn)
	c := sngecomm.Nmsgs() // get message count
	fmt.Println(sngecomm.ExampIdNow(exampid), "startReceivers message count", c, qn)
	for i := 1; i <= qn; i++ { // all queues
		wgrecv.Add(1)
		go receiver(conn, i, c)
	}
	wgrecv.Wait()

	// Disconnect from Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Println(sngecomm.ExampIdNow(exampid), "startReceivers disconnect error", e, qn) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		log.Println(sngecomm.ExampIdNow(exampid), "startReceivers netclose error", e, qn) // Handle this ......
	}

	fmt.Println(sngecomm.ExampIdNow(exampid), "startReceivers ends", qn)
	sngecomm.ShowStats(exampid, "startReceivers", conn)
	wgall.Done()
}

// Show a number of writers and readers operating concurrently from unique
// destinations.
func main() {
	sngecomm.ShowRunParms(exampid)

	if sngecomm.Pprof() {
		cfg := profile.Config{
			MemProfile:     true,
			CPUProfile:     true,
			BlockProfile:   true,
			NoShutdownHook: false, // Hook SIGINT
		}
		defer profile.Start(&cfg).Stop()
	}

	tn := time.Now()
	fmt.Println(sngecomm.ExampIdNow(exampid), "main starts")

	if sngecomm.SetMAXPROCS() {
		nc := runtime.NumCPU()
		fmt.Println(sngecomm.ExampIdNow(exampid), "main number of CPUs is:", nc)
		c := runtime.GOMAXPROCS(nc)
		fmt.Println(sngecomm.ExampIdNow(exampid), "main previous number of GOMAXPROCS is:", c)
		fmt.Println(sngecomm.ExampIdNow(exampid), "main current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	}
	//
	send_wait = sngecomm.SendWait()
	recv_wait = sngecomm.RecvWait()
	sendFact = sngecomm.SendFactor()
	recvFact = sngecomm.RecvFactor()
	fmt.Println(sngecomm.ExampIdNow(exampid), "main Sleep Factors", "send", sendFact, "recv", recvFact)
	//
	q := sngecomm.Nqs()
	//
	wgall.Add(2)
	go startReceivers(q)
	go startSenders(q)
	wgall.Wait()

	fmt.Println(sngecomm.ExampIdNow(exampid), "main ends", time.Since(tn))
}
