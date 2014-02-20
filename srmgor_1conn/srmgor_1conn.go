//
// Copyright © 2011-2014 Guy M. Allard
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
Send and receive many STOMP messages using multiple queues and goroutines
to service each send or receive instance. All senders and receivers share the
same STOMP connection.
*/
package main

import (
	"fmt"
	"github.com/davecheney/profile"
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
	"log"
	//	"os"
	"net"
	"runtime"
	"sync"
	"time"
)

var exampid = "srmgor_1conn:"

var wgsend sync.WaitGroup
var wgrecv sync.WaitGroup
var wgall sync.WaitGroup

// We 'stagger' between each message send and message receive for a random
// amount of time.
// Vary these for experimental purposes.  YMMV.
var max int64 = 1e9      // Max stagger time (nanoseconds)
var min int64 = max / 10 // Min stagger time (nanoseconds)

// Wait flags
var sendWait = true
var recvWait = true

// Sleep multipliers
var sendFact float64 = 1.0
var recvFact float64 = 1.0

//
var n net.Conn                // Network Connection
var conn *stompngo.Connection // Stomp Connection

var lhl = 44

// Send messages to a particular queue
func sender(qn, c int) {
	qns := fmt.Sprintf("%d", qn) // queue number
	id := stompngo.Uuid()        // A unique sender id
	fmt.Println(exampid, time.Now(), id, "send start", qn)
	//
	qp := sngecomm.Dest() // queue name prefix
	q := qp + "." + qns
	fmt.Println(exampid, time.Now(), id, "send queue name:", q, qn)
	h := stompngo.Headers{"destination", q, "senderId", id,
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
			log.Fatalln(exampid, time.Now(), id, "send error", e, qn)
		}
		if sendWait {
			runtime.Gosched()
			d := time.Duration(sngecomm.ValueBetween(min, max, sendFact))
			fmt.Println(sngecomm.ExampIdNow(exampid), "send", "stagger", int64(d)/1000000, "ms")
			tmr.Reset(d)
			_ = <-tmr.C
		}
	}
	// Sending is done
	fmt.Println(exampid, time.Now(), id, "send ends", qn)
	wgsend.Done()
}

// Receive messages from a particular queue
func receiver(qn, c int) {
	qns := fmt.Sprintf("%d", qn) // queue number
	id := stompngo.Uuid()        // A unique subscription ID
	fmt.Println(exampid, time.Now(), id, "recv starts", qn)
	//
	qp := sngecomm.Dest() // queue name prefix
	q := qp + "." + qns
	fmt.Println(exampid, time.Now(), id, "recv queue name:", q, qn)
	// Subscribe
	r := sngecomm.Subscribe(conn, q, id, "auto")
	//
	tmr := time.NewTimer(100 * time.Hour)
	// Receive loop
	for i := 1; i <= c; i++ {
		fmt.Println(exampid, time.Now(), id, "recv chanchek", "q", qns, "len", len(r), "cap", cap(r))
		d := <-r
		if d.Error != nil {
			log.Fatalln(exampid, time.Now(), id, "recv error", d.Error, qn)
		}

		// Process the inbound message .................
		osl := lhl
		if len(d.Message.Body) < osl {
			osl = len(d.Message.Body)
		}
		os := string(d.Message.Body[0:osl])
		fmt.Println(exampid, time.Now(), id, "recv message", os, qn)

		// Sanity check the queue and message numbers
		mns := fmt.Sprintf("%d", i) // message number
		if !d.Message.Headers.ContainsKV("qnum", qns) || !d.Message.Headers.ContainsKV("msgnum", mns) {
			log.Fatalln("Bad Headers", d.Message.Headers, qns, mns)
		}
		if recvWait {
			runtime.Gosched()
			d := time.Duration(sngecomm.ValueBetween(min, max, recvFact))
			fmt.Println(sngecomm.ExampIdNow(exampid), "recv", "stagger", int64(d)/1000000, "ms")
			tmr.Reset(d)
			_ = <-tmr.C
		}
	}
	// Unsubscribe
	sngecomm.Unsubscribe(conn, q, id)

	// Receiving is done
	fmt.Println(exampid, time.Now(), id, "recv ends", qn)
	wgrecv.Done()
}

func startSenders(qn int) {
	fmt.Println(exampid, time.Now(), "startSenders starts", qn)

	c := sngecomm.Nmsgs() // message count
	fmt.Println(exampid, time.Now(), "startSenders message count", c, qn)
	for i := 1; i <= qn; i++ { // all queues
		wgsend.Add(1)
		go sender(i, c)
	}
	wgsend.Wait()

	fmt.Println(exampid, time.Now(), "startSenders ends", qn)
	wgall.Done()
}

func startReceivers(qn int) {
	fmt.Println(exampid, time.Now(), "startReceivers starts", qn)

	c := sngecomm.Nmsgs() // get message count
	fmt.Println(exampid, time.Now(), "startReceivers message count:", c, qn)
	for i := 1; i <= qn; i++ { // all queues
		wgrecv.Add(1)
		go receiver(i, c)
	}
	wgrecv.Wait()

	fmt.Println(exampid, time.Now(), "startReceivers ends", qn)
	wgall.Done()
}

// Show a number of writers and readers operating concurrently from unique
// destinations.
func main() {
	sngecomm.ShowRunParms()
	cfg := profile.Config{
		MemProfile:     true,
		CPUProfile:     true,
		BlockProfile:   true,
		NoShutdownHook: false, // Hook SIGINT
	}
	pp := profile.Start(&cfg)
	//
	start := time.Now()
	fmt.Println(exampid, time.Now(), "main starts")
	if sngecomm.SetMAXPROCS() {
		nc := runtime.NumCPU()
		fmt.Println(exampid, time.Now(), "main number of CPUs is:", nc)
		c := runtime.GOMAXPROCS(nc)
		fmt.Println(exampid, time.Now(), "main previous number of GOMAXPROCS is:", c)
		fmt.Println(exampid, time.Now(), "main current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	}
	// Wait flags
	sendWait = sngecomm.SendWait()
	recvWait = sngecomm.RecvWait()
	sendFact = sngecomm.SendFactor()
	recvFact = sngecomm.RecvFactor()
	fmt.Println(sngecomm.ExampIdNow(exampid), "Sleep Factors", "send", sendFact, "recv", recvFact)
	// Number of queues
	q := sngecomm.Nqs()
	fmt.Println(exampid, time.Now(), "main Nqs:", q)
	// Open net and stomp connections
	h, p := sngecomm.HostAndPort() // network connection host and port
	var e error
	// Network open
	n, e = net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(exampid, time.Now(), "main dial error", e) // Handle this ......
	}
	// Stomp connect, 1.1(+)
	ch := sngecomm.ConnectHeaders()
	log.Println(exampid, time.Now(), "vhost:", sngecomm.Vhost(), "protocol:", sngecomm.Protocol())
	conn, e = stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(exampid, time.Now(), "main connect error", e) // Handle this ......
	}

	// Many receivers running under the same connection can cause
	// (wire read) performance issues.  This is *very* dependent on the broker
	// being used, specifically the broker's algorithm for putting messages on
	// the wire.
	// To alleviate those issues, this strategy insures that messages are
	// received from the wire as soon as possible.  Those messages are then
	// buffered internally for (possibly later) application processing. In
	// this example, buffering occurs in the stompngo package.
	conn.SetSubChanCap(sngecomm.SubChanCap()) // Experiment with this value, YMMV

	// Set up a logger
	// l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	// conn.SetLogger(l)

	// Run everything

	wgall.Add(2)
	go startReceivers(q)
	go startSenders(q)
	wgall.Wait()
	pp.Stop()

	// Disconnect from Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(exampid, time.Now(), "main disconnect error", e) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		log.Fatalln(exampid, time.Now(), "main netclose error", e) // Handle this ......
	}
	tn := time.Now().String()
	sngecomm.ShowStats(exampid+" "+tn, "done", conn)
	dur := time.Since(start)
	fmt.Println(exampid, time.Now(), "main ends", dur)
}
