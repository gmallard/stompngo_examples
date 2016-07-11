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
Send and receive many STOMP messages using multiple queues and goroutines
to service each send or receive instance. All senders and receivers share the
same STOMP connection.
*/
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
	//
	"github.com/davecheney/profile"
	//
	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	ll = log.New(os.Stdout, "ECNDS ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	exampid = "srmgor_1conn:"

	wgs sync.WaitGroup
	wgr sync.WaitGroup
	wga sync.WaitGroup

	// We 'stagger' between each message send and message receive for a random
	// amount of time.
	// Vary these for experimental purposes.  YMMV.
	max int64 = 1e9      // Max stagger time (nanoseconds)
	min int64 = max / 10 // Min stagger time (nanoseconds)

	// Wait flags
	sw = true
	rw = true

	// Sleep multipliers
	sf float64 = 1.0
	rf float64 = 1.0

	//
	n    net.Conn             // Network Connection
	conn *stompngo.Connection // Stomp Connection

	lhl = 44
)

// Send messages to a particular queue
func sender(qn, mc int) {
	qns := fmt.Sprintf("%d", qn) // string queue number
	id := stompngo.Uuid()        // A unique sender id
	ll.Printf("%s id:%s send_start_queue_number qn:%d\n", exampid, id, qn)
	//
	d := senv.Dest() + "." + qns
	ll.Printf("%s id:%s send_queue_name:%s qn:%d\n", exampid, id, d, qn)
	wh := stompngo.Headers{"destination", d, "senderId", id,
		"qnum", qns} // send Headers
	if senv.Persistent() {
		wh = wh.Add("persistent", "true")
	}
	//
	tmr := time.NewTimer(100 * time.Hour)
	// Send loop
	for i := 1; i <= mc; i++ {
		si := fmt.Sprintf("%d", i)
		sh := append(wh, "msgnum", si)
		// Generate a message to send ...............
		ll.Printf("%s id:%s send_message qn:%d msgnum:%s\n", exampid, id, qn, si)
		e := conn.Send(sh, string(sngecomm.Partial()))
		if e != nil {
			ll.Fatalf("%s v1:%v v2:%v v3:%v v4:%v\n", exampid, id, "send error", e, qn)
		}
		if sw {
			dt := time.Duration(sngecomm.ValueBetween(min, max, sf))
			ll.Printf("%s send_stagger dt:%v qn:%d id:%s\n",
				exampid, dt,
				qn, id)
			tmr.Reset(dt)
			_ = <-tmr.C
			runtime.Gosched()
		}
	}
	// Sending is done
	ll.Printf("%s id:%s send_end_queue_number qn:%d\n", exampid, id, qn)
	wgs.Done()
}

// Receive messages from a particular queue
func receiver(qn, mc int) {
	qns := fmt.Sprintf("%d", qn) // string queue number
	pbc := sngecomm.Pbc()
	id := stompngo.Uuid() // A unique subscription ID
	//
	d := senv.Dest() + "." + qns
	ll.Printf("%s id:%s recv_queue_name:%s qn:%s\n", exampid, id, d, qns)
	// Subscribe
	sc := sngecomm.HandleSubscribe(conn, d, id, sngecomm.AckMode())
	//
	tmr := time.NewTimer(100 * time.Hour)
	var md stompngo.MessageData
	// Receive loop
	for i := 1; i <= mc; i++ {
		ll.Printf("%s id:%s recv_ranchek qn:%d chlen:%d chcap:%d\n", exampid, id,
			qn, len(sc), cap(sc))

		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// A RECEIPT or ERROR frame is unexpected here
			ll.Fatalf("%s v1:%v\n", exampid, md) // Handle this
		}
		if md.Error != nil {
			ll.Fatalf("%s v1:%v v2:%v v3:%v v4:%v\n", exampid, id, "recv error", md.Error, qn)
		}

		// Process the inbound message .................
		ll.Printf("%s id:%s recv_message qn:%d msgnum:%d\n", exampid, id, qn, i)
		if pbc > 0 {
			maxlen := pbc
			if len(md.Message.Body) < maxlen {
				maxlen = len(md.Message.Body)
			}
			ss := string(md.Message.Body[0:maxlen])
			ll.Printf("%s Payload: %s\n", exampid, ss) // Data payload
		}

		// Sanity check the message Command, and the queue and message numbers
		mns := fmt.Sprintf("%d", i) // message number
		if md.Message.Command != stompngo.MESSAGE {
			ll.Fatalf("%s v1:%v v2:%v v3:%v v4:%v\n", exampid, "Bad Frame", md, qn, mns)
		}
		if !md.Message.Headers.ContainsKV("qnum", qns) || !md.Message.Headers.ContainsKV("msgnum", mns) {
			ll.Fatalf("%s v1:%v v2:%v v3:%v v4:%v\n", exampid, "Bad Headers", md.Message.Headers, qn, mns)
		}

		if rw {
			dt := time.Duration(sngecomm.ValueBetween(min, max, rf))
			ll.Printf("%s recv_stagger dt:%v qn:%d id:%s\n",
				exampid, dt,
				qn, id)
			tmr.Reset(dt)
			_ = <-tmr.C
			runtime.Gosched()
		}

		// Handle ACKs if needed
		if sngecomm.AckMode() != "auto" {
			sngecomm.HandleAck(conn, md.Message.Headers, id)
		}
	}
	// Unsubscribe
	sngecomm.HandleUnsubscribe(conn, d, id)

	// Receiving is done
	ll.Printf("%s id:%s recv_end_queue_number qn:%d\n", exampid, id, qn)
	wgr.Done()
}

/*
	Start all sender go routines.
*/
func startSenders(nqs int) {
	ll.Printf("%s startSenders_starts nqs:%d\n", exampid, nqs)

	mc := senv.Nmsgs() // message count
	//	nqs)
	ll.Printf("%s startSenders_message_count mc:%d nqs:%d\n", exampid, mc, nqs)
	for i := 1; i <= nqs; i++ { // all queues
		wgs.Add(1)
		go sender(i, mc)
	}
	wgs.Wait()

	ll.Printf("%s startSenders_endsexampid nqs:%d\n", exampid, nqs)
	wga.Done()
}

/*
	Start all receiver go routines.
*/
func startReceivers(nqs int) {
	ll.Printf("%s startReceivers_starts nqs:%d\n", exampid, "startReceivers starts", nqs)

	mc := senv.Nmsgs() // get message count
	ll.Printf("%s startReceivers_message_count mc:%d nqs:%d\n", exampid, mc, nqs)
	for i := 1; i <= nqs; i++ { // all queues
		wgr.Add(1)
		go receiver(i, mc)
	}
	wgr.Wait()

	ll.Printf("%s startReceivers_ends nqs:%d\n", exampid, nqs)
	wga.Done()
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

	start := time.Now()
	ll.Printf("%s v1:%v\n", exampid, "main starts")
	ll.Printf("%s v1:%v v2:%v\n", exampid, "main profiling", sngecomm.Pprof())
	ll.Printf("%s v1:%v v2:%v\n", exampid, "main current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	if sngecomm.SetMAXPROCS() {
		nc := runtime.NumCPU()
		ll.Printf("%s v1:%v v2:%v\n", exampid, "main number of CPUs is:", nc)
		gmp := runtime.GOMAXPROCS(nc)
		ll.Printf("%s v1:%v v2:%v\n", exampid, "main previous number of GOMAXPROCS is:", gmp)
		ll.Printf("%s v1:%v v2:%v\n", exampid, "main current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	}
	// Wait flags
	sw = sngecomm.SendWait()
	rw = sngecomm.RecvWait()
	sf = sngecomm.SendFactor()
	rf = sngecomm.RecvFactor()
	ll.Printf("%s v1:%v v2:%v v3:%v v4:%v v5:%v\n", exampid, "main Sleep Factors", "send", sf, "recv", rf)
	// Number of queues
	nqs := sngecomm.Nqs()
	// Open net and stomp connections
	h, p := senv.HostAndPort() // network connection host and port
	hap := net.JoinHostPort(h, p)
	var e error
	// Network open
	n, e = net.Dial("tcp", hap)
	if e != nil {
		ll.Fatalf("%s v1:%v v2:%v\n", exampid, "main dial error", e) // Handle this ......
	}
	// Stomp connect, 1.1(+)
	ch := sngecomm.ConnectHeaders()
	ll.Printf("%s v1:%v v2:%v v3:%v v4:%v\n", exampid, "vhost:", senv.Vhost(), "protocol:", senv.Protocol())
	conn, e = stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalf("%s v1:%v v2:%v\n", exampid, "main connect error", e) // Handle this ......
	}

	// Many receivers running under the same connection can cause
	// (wire read) performance issues.  This is *very* dependent on the broker
	// being used, specifically the broker's algorithm for putting messages on
	// the wire.
	// To alleviate those issues, this strategy insures that messages are
	// received from the wire as soon as possible.  Those messages are then
	// buffered internally for (possibly later) application processing. In
	// this example, buffering occurs in the stompngo package.
	conn.SetSubChanCap(senv.SubChanCap()) // Experiment with this value, YMMV

	// Run everything
	wga.Add(2)
	go startReceivers(nqs)
	go startSenders(nqs)
	wga.Wait()

	// Disconnect from Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalf("%s v1:%v v2:%v\n", exampid, "main disconnect error", e) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		ll.Fatalf("%s v1:%v v2:%v\n", exampid, "main netclose error", e) // Handle this ......
	}
	sngecomm.ShowStats(exampid, "done", conn)
	dur := time.Since(start)
	ll.Printf("%s v1:%v v2:%v\n", exampid, "main ends", dur)
}
