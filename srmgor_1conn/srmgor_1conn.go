//
// Copyright © 2011-2018 Guy M. Allard
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
	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	ll = log.New(os.Stdout, "ECNDS ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	exampid = "srmgor_1conn: "

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

	tag = "1conn"
)

// Send messages to a particular queue
func sender(qn, mc int) {
	ltag := tag + "-sender"

	qns := fmt.Sprintf("%d", qn) // string queue number
	id := stompngo.Uuid()        // A unique sender id
	d := sngecomm.Dest() + "." + string(exampid[:len(exampid)-2]) + "." + qns

	ll.Printf("%stag:%s connsess:%s queue_info id:%v d:%v qnum:%v mc:%v\n",
		exampid, ltag, conn.Session(),
		id, d, qn, mc)
	//
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
		ll.Printf("%stag:%s connsess:%s send_headers id:%v d:%v qnum:%v headers:%v\n",
			exampid, ltag, conn.Session(),
			id, d, qn, sh)
		e := conn.Send(sh, string(sngecomm.Partial()))
		if e != nil {
			ll.Fatalf("%stag:%s connsess:%s send_error qnum:%v error:%v",
				exampid, tag, conn.Session(),
				qn, e.Error()) // Handle this ......
		}
		if i == mc {
			break
		}
		if sw {
			dt := time.Duration(sngecomm.ValueBetween(min, max, sf))
			ll.Printf("%stag:%s connsess:%s send_stagger id:%v d:%v qnum:%v stagger:%v\n",
				exampid, ltag, conn.Session(),
				id, d, qn, dt)
			tmr.Reset(dt)
			_ = <-tmr.C
			runtime.Gosched()
		}
	}
	// Sending is done
	ll.Printf("%stag:%s connsess:%s finish_info id:%v d:%v qnum:%v mc:%v\n",
		exampid, ltag, conn.Session(),
		id, d, qn, mc)
	wgs.Done()
}

// Receive messages from a particular queue
func receiver(qn, mc int) {
	ltag := tag + "-receiver"

	qns := fmt.Sprintf("%d", qn) // string queue number
	pbc := sngecomm.Pbc()
	id := stompngo.Uuid() // A unique subscription ID
	d := sngecomm.Dest() + "." + string(exampid[:len(exampid)-2]) + "." + qns

	ll.Printf("%stag:%s connsess:%s queue_info id:%v d:%v qnum:%v mc:%v\n",
		exampid, ltag, conn.Session(),
		id, d, qn, mc)
	// Subscribe
	sc := sngecomm.HandleSubscribe(conn, d, id, sngecomm.AckMode())
	ll.Printf("%stag:%s connsess:%s subscribe_complete id:%v d:%v qnum:%v mc:%v\n",
		exampid, ltag, conn.Session(),
		id, d, qn, mc)
	//
	tmr := time.NewTimer(100 * time.Hour)
	var md stompngo.MessageData
	// Receive loop
	for i := 1; i <= mc; i++ {
		ll.Printf("%stag:%s connsess:%s recv_ranchek id:%v d:%v qnum:%v mc:%v chlen:%v chcap:%v\n",
			exampid, ltag, conn.Session(),
			id, d, qn, mc, len(sc), cap(sc))

		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// A RECEIPT or ERROR frame is unexpected here
			ll.Fatalf("%stag:%s connsess:%s bad_frame qnum:%v headers:%v body:%s",
				exampid, tag, conn.Session(),
				qn, md.Message.Headers, md.Message.Body) // Handle this ......
		}
		if md.Error != nil {
			ll.Fatalf("%stag:%s connsess:%s recv_error qnum:%v error:%v",
				exampid, tag, conn.Session(),
				qn, md.Error) // Handle this ......
		}

		// Process the inbound message .................
		ll.Printf("%stag:%s connsess:%s recv_message qnum:%v i:%v\n",
			exampid, tag, conn.Session(),
			qn, i)
		if pbc > 0 {
			maxlen := pbc
			if len(md.Message.Body) < maxlen {
				maxlen = len(md.Message.Body)
			}
			ss := string(md.Message.Body[0:maxlen])
			ll.Printf("%stag:%s connsess:%s payload qnum:%v body:%s\n",
				exampid, tag, conn.Session(),
				qn, ss)

		}

		// Sanity check the message Command, and the queue and message numbers
		mns := fmt.Sprintf("%d", i) // message number
		if md.Message.Command != stompngo.MESSAGE {
			ll.Fatalf("%stag:%s connsess:%s bad_frame qnum:%v command:%v headers:%v body:%v\n",
				exampid, tag, conn.Session(),
				qn, md.Message.Command, md.Message.Headers, string(md.Message.Body)) // Handle this ......

		}
		if !md.Message.Headers.ContainsKV("qnum", qns) || !md.Message.Headers.ContainsKV("msgnum", mns) {
			ll.Fatalf("%stag:%s connsess:%s dirty_message qns:%v msgnum:%v command:%v headers:%v body:%v\n",
				exampid, tag, conn.Session(),
				qns, mns, md.Message.Command, md.Message.Headers, string(md.Message.Body)) // Handle this ......) // Handle this ......
		}

		if i == mc {
			break
		}

		if rw {
			dt := time.Duration(sngecomm.ValueBetween(min, max, rf))
			ll.Printf("%stag:%s connsess:%s recv_stagger id:%v d:%v qnum:%v stagger:%v\n",
				exampid, ltag, conn.Session(),
				id, d, qn, dt)
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
	ll.Printf("%stag:%s connsess:%s unsubscribe_complete id:%v d:%v qnum:%v mc:%v\n",
		exampid, ltag, conn.Session(),
		id, d, qn, mc)

	// Receiving is done
	ll.Printf("%stag:%s connsess:%s recv_end id:%v d:%v qnum:%v mc:%v\n",
		exampid, ltag, conn.Session(),
		id, d, qn, mc)

	wgr.Done()
}

/*
	Start all sender go routines.
*/
func startSenders(nqs int) {
	ltag := tag + "-startsenders"

	ll.Printf("%stag:%s connsess:%s queue_count nqs:%v\n",
		exampid, ltag, conn.Session(),
		nqs)

	mc := senv.Nmsgs() // message count
	ll.Printf("%stag:%s connsess:%s message_count mc:%v\n",
		exampid, ltag, conn.Session(),
		mc)
	for i := 1; i <= nqs; i++ { // all queues
		wgs.Add(1)
		go sender(i, mc)
	}
	wgs.Wait()

	ll.Printf("%stag:%s connsess:%s ends nqs:%v mc:%v\n",
		exampid, ltag, conn.Session(),
		nqs, mc)
	wga.Done()
}

/*
	Start all receiver go routines.
*/
func startReceivers(nqs int) {
	ltag := tag + "-startreceivers"

	ll.Printf("%stag:%s connsess:%s queue_count nqs:%v\n",
		exampid, ltag, conn.Session(),
		nqs)

	mc := senv.Nmsgs() // get message count
	ll.Printf("%stag:%s connsess:%s message_count mc:%v\n",
		exampid, ltag, conn.Session(),
		mc)

	for i := 1; i <= nqs; i++ { // all queues
		wgr.Add(1)
		go receiver(i, mc)
	}
	wgr.Wait()

	ll.Printf("%stag:%s connsess:%s ends nqs:%v mc:%v\n",
		exampid, ltag, conn.Session(),
		nqs, mc)
	wga.Done()
}

// Show a number of writers and readers operating concurrently from unique
// destinations.
func main() {

	st := time.Now()

	sngecomm.ShowRunParms(exampid)

	ll.Printf("%stag:%s connsess:%s main_starts\n",
		exampid, tag, sngecomm.Lcs)

	ll.Printf("%stag:%s connsess:%s main_profiling pprof:%v\n",
		exampid, tag, sngecomm.Lcs,
		sngecomm.Pprof())

	ll.Printf("%stag:%s connsess:%s main_current_GOMAXPROCS gmp:%v\n",
		exampid, tag, sngecomm.Lcs,
		runtime.GOMAXPROCS(-1))

	if sngecomm.SetMAXPROCS() {
		nc := runtime.NumCPU()
		ll.Printf("%stag:%s connsess:%s main_current_num_cpus cncpu:%v\n",
			exampid, tag, sngecomm.Lcs,
			nc)
		gmp := runtime.GOMAXPROCS(nc)
		ll.Printf("%stag:%s connsess:%s main_previous_num_cpus pncpu:%v\n",
			exampid, tag, sngecomm.Lcs,
			gmp)
		ll.Printf("%stag:%s connsess:%s main_current_GOMAXPROCS gmp:%v\n",
			exampid, tag, sngecomm.Lcs,
			runtime.GOMAXPROCS(-1))
	}
	// Wait flags
	sw = sngecomm.SendWait()
	rw = sngecomm.RecvWait()
	sf = sngecomm.SendFactor()
	rf = sngecomm.RecvFactor()
	ll.Printf("%stag:%s connsess:%s main_wait_sleep_factors sw:%v rw:%v sf:%v rf:%v\n",
		exampid, tag, sngecomm.Lcs,
		sw, rw, sf, rf)
	// Number of queues
	nqs := sngecomm.Nqs()

	// Standard example connect sequence
	var e error
	n, conn, e = sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		if conn != nil {
			ll.Printf("%stag:%s  connsess:%s Connect Response headers:%v body%s\n",
				exampid, tag, conn.Session(), conn.ConnectResponse.Headers,
				string(conn.ConnectResponse.Body))
		}
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
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

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_disconnect error:%v",
			exampid, tag, conn.Session(),
			e.Error()) // Handle this ......
	}

	sngecomm.ShowStats(exampid, tag, conn)

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, conn.Session(),
		time.Now().Sub(st))

	time.Sleep(250 * time.Millisecond)
}
