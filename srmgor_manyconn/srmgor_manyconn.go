//
// Copyright © 2012-2016 Guy M. Allard
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
to service each send or receive instance.  Each sender and receiver
operates under a unique network connection.

	Examples:

		# A few queues and a few messages:
		STOMP_NQS=5 STOMP_NMSGS=10 go run srmgor_manyconn.go

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
	exampid = "srmgor_manyconn: "

	wgs sync.WaitGroup
	wgr sync.WaitGroup

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

	// Number of messages
	nmsgs = senv.Nmsgs()

	ll = log.New(os.Stdout, "EMSMR ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	tag = "manyconn"
)

func sendMessages(conn *stompngo.Connection, qnum int, nc net.Conn) {
	ltag := tag + "-sendmessages"

	qns := fmt.Sprintf("%d", qnum) // queue number
	d := sngecomm.Dest() + "." + qns
	ll.Printf("%stag:%s connsess:%s start d:%s qnum:%d\n",
		exampid, ltag, conn.Session(),
		d, qnum)
	wh := stompngo.Headers{"destination", d,
		"qnum", qns} // send Headers
	if senv.Persistent() {
		wh = wh.Add("persistent", "true")
	}
	//
	tmr := time.NewTimer(100 * time.Hour)
	// Send messages
	for mc := 1; mc <= nmsgs; mc++ {
		mcs := fmt.Sprintf("%d", mc)
		sh := append(wh, "msgnum", mcs)
		// Generate a message to send ...............

		ll.Printf("%stag:%s connsess:%s message mc:%d qnum:%d\n",
			exampid, ltag, conn.Session(),
			mc, qnum)
		e := conn.Send(sh, string(sngecomm.Partial()))
		if e != nil {
			ll.Fatalf("%stag:%s connsess:%s send_error qnum:%v error:%v",
				exampid, ltag, conn.Session(),
				qnum, e.Error()) // Handle this ......
		}
		if mc == nmsgs {
			break
		}
		if sw {
			runtime.Gosched() // yield for this example
			dt := time.Duration(sngecomm.ValueBetween(min, max, sf))
			ll.Printf("%stag:%s connsess:%s send_stagger dt:%v qnum:%s mc:%d\n",
				exampid, ltag, conn.Session(),
				dt, qnum, mc)
			tmr.Reset(dt)
			_ = <-tmr.C
		}
	}
}

func receiveMessages(conn *stompngo.Connection, qnum int, nc net.Conn) {
	ltag := tag + "-receivemessages"

	qns := fmt.Sprintf("%d", qnum) // queue number
	d := sngecomm.Dest() + "." + qns
	id := stompngo.Uuid() // A unique subscription ID

	ll.Printf("%stag:%s connsess:%s receiveMessages_start id:%s d:%s qnum:%d nmsgs:%d\n",
		exampid, ltag, conn.Session(),
		id, d, qnum, nmsgs)
	// Subscribe
	sc := sngecomm.HandleSubscribe(conn, d, id, sngecomm.AckMode())

	pbc := sngecomm.Pbc() // Print byte count

	//
	tmr := time.NewTimer(100 * time.Hour)
	var md stompngo.MessageData
	for mc := 1; mc <= nmsgs; mc++ {

		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// Frames RECEIPT or ERROR not expected here
			ll.Fatalf("%stag:%s connsess:%s send_error qns:%v md:%v",
				exampid, ltag, conn.Session(),
				qns, md) // Handle this ......
		}
		if md.Error != nil {
			ll.Fatalf("%stag:%s connsess:%s receive_error qns:%v error:%v\n",
				exampid, ltag, conn.Session(),
				qns, md.Error)
		}

		if md.Message.Command != stompngo.MESSAGE {
			ll.Fatalf("%stag:%s connsess:%s bad_frame qns:%s mc:%d md:%v\n",
				exampid, ltag, conn.Session(),
				qns, mc, md)
		}

		mcs := fmt.Sprintf("%d", mc) // message number
		if !md.Message.Headers.ContainsKV("qnum", qns) || !md.Message.Headers.ContainsKV("msgnum", mcs) {
			ll.Fatalf("%stag:%s connsess:%s dirty_message qns:%v msgnum:%v md:%v",
				exampid, tag, conn.Session(),
				qns, mcs, md) // Handle this ......
		}

		// Process the inbound message .................
		sl := len(md.Message.Body)
		if pbc > 0 {
			sl = pbc
			if len(md.Message.Body) < sl {
				sl = len(md.Message.Body)
			}
		}
		ll.Printf("%stag:%s connsess:%s receiveMessages_msg d:%s body:%s qnum:%d msgnum:%s\n",
			exampid, ltag, conn.Session(),
			d, string(md.Message.Body[0:sl]), qnum,
			md.Message.Headers.Value("msgnum"))
		if mc == nmsgs {
			break
		}
		// Handle ACKs if needed
		if sngecomm.AckMode() != "auto" {
			ah := []string{}
			sngecomm.HandleAck(conn, ah, id)
		}
		if mc == nmsgs {
			break
		}
		//
		if rw {
			runtime.Gosched() // yield for this example
			dt := time.Duration(sngecomm.ValueBetween(min, max, rf))
			ll.Printf("%stag:%s connsess:%s recv_stagger dt:%v qns:%s mc:%d\n",
				exampid, ltag, conn.Session(),
				dt, qns, mc)
			tmr.Reset(dt)
			_ = <-tmr.C
		}
	}
	ll.Printf("%stag:%s connsess:%s end d:%s qnum:%d nmsgs:%d\n",
		exampid, ltag, conn.Session(),
		d, qnum, nmsgs)

	// Unsubscribe
	sngecomm.HandleUnsubscribe(conn, d, id)
	//
}

func runReceiver(qnum int) {
	ltag := tag + "-runreceiver"

	ll.Printf("%stag:%s connsess:%s start qnum:%d\n",
		exampid, ltag, sngecomm.Lcs,
		qnum)

	// Standard example connect sequence
	n, conn, e := sngecomm.CommonConnect(exampid, ltag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s error:%s\n",
			exampid, ltag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	//
	conn.SetSubChanCap(senv.SubChanCap()) // Experiment with this value, YMMV
	// Receives
	receiveMessages(conn, qnum, n)

	ll.Printf("%stag:%s connsess:%s receives_complete qnum:%d\n",
		exampid, ltag, conn.Session(),
		qnum)

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, ltag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s error:%s\n",
			exampid, ltag, conn.Session(),
			e.Error()) // Handle this ......
	}

	sngecomm.ShowStats(exampid, "recv_"+fmt.Sprintf("%d", qnum), conn)
	wgr.Done()
}

func runSender(qnum int) {

	ltag := tag + "-runsender"

	ll.Printf("%stag:%s connsess:%s start qnum:%d\n",
		exampid, ltag, sngecomm.Lcs,
		qnum)
	// Standard example connect sequence
	n, conn, e := sngecomm.CommonConnect(exampid, ltag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s error:%s\n",
			exampid, ltag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	//
	sendMessages(conn, qnum, n)

	ll.Printf("%stag:%s connsess:%s sends_complete qnum:%d\n",
		exampid, ltag, conn.Session(),
		qnum)

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, ltag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s error:%s\n",
			exampid, ltag, conn.Session(),
			e.Error()) // Handle this ......
	}
	sngecomm.ShowStats(exampid, "send_"+fmt.Sprintf("%d", qnum), conn)
	wgs.Done()
}

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
	//
	sw = sngecomm.SendWait()
	rw = sngecomm.RecvWait()
	sf = sngecomm.SendFactor()
	rf = sngecomm.RecvFactor()
	ll.Printf("%stag:%s connsess:%s main_wait_sleep_factors sw:%v rw:%v sf:%v rf:%v\n",
		exampid, tag, sngecomm.Lcs,
		sw, rw, sf, rf)
	//
	numq := sngecomm.Nqs()
	nmsgs = senv.Nmsgs() // message count
	//
	ll.Printf("%stag:%s connsess:%s main_starting_receivers\n",
		exampid, tag, sngecomm.Lcs)
	for q := 1; q <= numq; q++ {
		wgr.Add(1)
		go runReceiver(q)
	}
	ll.Printf("%stag:%s connsess:%s main_started_receivers\n",
		exampid, tag, sngecomm.Lcs)
	//
	ll.Printf("%stag:%s connsess:%s main_starting_senders\n",
		exampid, tag, sngecomm.Lcs)
	for q := 1; q <= numq; q++ {
		wgs.Add(1)
		go runSender(q)
	}
	ll.Printf("%stag:%s connsess:%s main_started_senders\n",
		exampid, tag, sngecomm.Lcs)
	//
	wgs.Wait()
	ll.Printf("%stag:%s connsess:%s main_senders_complete\n",
		exampid, tag, sngecomm.Lcs)
	wgr.Wait()
	ll.Printf("%stag:%s connsess:%s main_receivers_complete\n",
		exampid, tag, sngecomm.Lcs)
	//

	// The end
	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, sngecomm.Lcs,
		time.Now().Sub(st))

}
