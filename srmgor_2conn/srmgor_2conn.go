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
	"os"
	"runtime"
	"strconv"
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
	exampid = "srmgor_2conn: "

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

	// Possible profile file
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

	ll = log.New(os.Stdout, "E1S1R ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	tag = "2conn"
)

// Send messages to a particular queue
func sender(conn *stompngo.Connection, qn, nmsgs int) {
	ltag := tag + "-sender"

	qns := fmt.Sprintf("%d", qn) // queue number
	d := sngecomm.Dest() + "."  + string(exampid[:len(exampid)-2]) + "." + qns
	ll.Printf("%stag:%s connsess:%s starts qn:%d nmsgs:%d d:%s\n",
		exampid, ltag, conn.Session(),
		qn, nmsgs, d)
	//
	wh := stompngo.Headers{"destination", d,
		"qnum", qns} // send Headers
	if senv.Persistent() {
		wh = wh.Add("persistent", "true")
	}
	//
	tmr := time.NewTimer(100 * time.Hour)
	// Send loop
	for i := 1; i <= nmsgs; i++ {
		si := fmt.Sprintf("%d", i)
		sh := append(wh, "msgnum", si)
		// Generate a message to send ...............
		ll.Printf("%stag:%s connsess:%s message qns:%s si:%s\n",
			exampid, ltag, conn.Session(),
			qns, si)
		e := conn.Send(sh, string(sngecomm.Partial()))
		if e != nil {
			ll.Fatalf("%stag:%s connsess:%s send_error qnum:%v error:%v",
				exampid, ltag, conn.Session(),
				qn, e.Error()) // Handle this ......
		}
		if i == nmsgs {
			break
		}
		if sw {
			runtime.Gosched() // yield for this example
			dt := time.Duration(sngecomm.ValueBetween(min, max, sf))
			ll.Printf("%stag:%s connsess:%s send_stagger dt:%v qns:%s\n",
				exampid, ltag, conn.Session(),
				dt, qns)
			tmr.Reset(dt)
			_ = <-tmr.C
		}
	}
	// Sending is done
	ll.Printf("%stag:%s connsess:%s sender_ends qn:%d nmsgs:%d\n",
		exampid, ltag, conn.Session(),
		qn, nmsgs)
	wgs.Done()
}

// Asynchronously process all messages for a given subscription.
func receiveWorker(sc <-chan stompngo.MessageData, qns string, nmsgs int,
	qc chan<- bool, conn *stompngo.Connection, id string) {
	//
	ltag := tag + "-receiveWorker"

	tmr := time.NewTimer(100 * time.Hour)

	pbc := sngecomm.Pbc() // Print byte count

	// Receive loop
	var md stompngo.MessageData
	for i := 1; i <= nmsgs; i++ {

		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// Frames RECEIPT or ERROR not expected here
			ll.Fatalf("%stag:%s connsess:%s bad_frame qns:%v md:%v",
				exampid, ltag, conn.Session(),
				qns, md) // Handle this ......
		}
		if md.Error != nil {
			ll.Fatalf("%stag:%s connsess:%s recv_error qns:%v error:%v",
				exampid, ltag, conn.Session(),
				qns, md.Error) // Handle this ......
		}

		// Sanity check the queue and message numbers
		mns := fmt.Sprintf("%d", i) // message number
		if !md.Message.Headers.ContainsKV("qnum", qns) || !md.Message.Headers.ContainsKV("msgnum", mns) {
			ll.Fatalf("%stag:%s connsess:%s dirty_message qnum:%v msgnum:%v md:%v",
				exampid, ltag, conn.Session(),
				qns, mns, md) // Handle this ......
		}

		// Process the inbound message .................
		sl := len(md.Message.Body)
		if pbc > 0 {
			sl = pbc
			if len(md.Message.Body) < sl {
				sl = len(md.Message.Body)
			}
		}

		// Handle ACKs if needed
		if sngecomm.AckMode() != "auto" {
			ah := []string{}
			sngecomm.HandleAck(conn, ah, id)
		}
		ll.Printf("%stag:%s connsess:%s recv_message body:%s qns:%s msgnum:%s i:%v\n",
			exampid, ltag, conn.Session(),
			string(md.Message.Body[0:sl]),
			qns,
			md.Message.Headers.Value("msgnum"), i)
		if i == nmsgs {
			break
		}
		if rw {
			runtime.Gosched() // yield for this example
			dt := time.Duration(sngecomm.ValueBetween(min, max, rf))
			ll.Printf("%stag:%s connsess:%s recv_stagger dt:%v qns:%s\n",
				exampid, ltag, conn.Session(),
				dt, qns)
			tmr.Reset(dt)
			_ = <-tmr.C
		}
	}
	//
	qc <- true
}

// Receive messages from a particular queue
func receiver(conn *stompngo.Connection, qn, nmsgs int) {
	ltag := tag + "-receiver"

	qns := fmt.Sprintf("%d", qn) // queue number
	ll.Printf("%stag:%s connsess:%s starts qns:%d nmsgs:%d\n",
		exampid, ltag, conn.Session(),
		qn, nmsgs)
	//
	qp := sngecomm.Dest() // queue name prefix
	q := qp + "."   + string(exampid[:len(exampid)-2]) + "." + qns
	ll.Printf("%stag:%s connsess:%s queue_info q:%s qn:%d nmsgs:%d\n",
		exampid, ltag, conn.Session(),
		q, qn, nmsgs)
	id := stompngo.Uuid() // A unique subscription ID
	sc := sngecomm.HandleSubscribe(conn, q, id, sngecomm.AckMode())
	ll.Printf("%stag:%s connsess:%s subscribe_complete\n",
		exampid, ltag, conn.Session())
	// Many receivers running under the same connection can cause
	// (wire read) performance issues.  This is *very* dependent on the broker
	// being used, specifically the broker's algorithm for putting messages on
	// the wire.
	// To alleviate those issues, this strategy insures that messages are
	// received from the wire as soon as possible.  Those messages are then
	// buffered internally for (possibly later) application processing.

	bs := -1 //
	if s := os.Getenv("STOMP_CONN2BUFFER"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			ll.Fatalf("%stag:%s connsess:%s CONN2BUFFER_conversion_error error:%v",
				exampid, ltag, conn.Session(),
				e.Error()) // Handle this ......

		} else {
			bs = int(i)
		}
	}
	if bs < 1 {
		bs = nmsgs
	}
	ll.Printf("%stag:%s connsess:%s mdbuffersize_qnum bs:%d qn:%d\n",
		exampid, ltag, conn.Session(),
		bs, qn)

	// Process all inputs async .......
	// var mc chan stompngo.MessageData
	mdc := make(chan stompngo.MessageData, bs)      // MessageData Buffer size
	dc := make(chan bool)                           // Receive processing done channel
	go receiveWorker(mdc, qns, nmsgs, dc, conn, id) // Start async processor
	for i := 1; i <= nmsgs; i++ {
		mdc <- <-sc // Receive message data as soon as possible, and internally queue it
	}
	ll.Printf("%stag:%s connsess:%s waitforWorkersBegin qns:%s\n",
		exampid, ltag, conn.Session(),
		qns)
	<-dc // Wait until receive processing is done for this queue
	ll.Printf("%stag:%s connsess:%s waitforWorkersEnd qns:%s\n",
		exampid, ltag, conn.Session(),
		qns)

	// Unsubscribe
	sngecomm.HandleUnsubscribe(conn, q, id)
	ll.Printf("%stag:%s connsess:%s unsubscribe_complete\n",
		exampid, ltag, conn.Session())

	// Receiving is done
	ll.Printf("%stag:%s connsess:%s ends qns:%s\n",
		exampid, ltag, conn.Session(),
		qns)
	wgr.Done()
}

func startSenders(qn int) {
	ltag := tag + "-startsenders"

	ll.Printf("%stag:%s connsess:%s queue qn:%v\n",
		exampid, ltag, sngecomm.Lcs,
		qn)

	// Standard example connect sequence
	n, conn, e := sngecomm.CommonConnect(exampid, ltag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s on_connect error:%v",
			exampid, ltag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	nmsgs := senv.Nmsgs() // message count
	ll.Printf("%stag:%s connsess:%s message_count nmsgs:%d qn:%d\n",
		exampid, ltag, conn.Session(),
		nmsgs, qn)
	for i := 1; i <= qn; i++ { // all queues
		wgs.Add(1)
		go sender(conn, i, nmsgs)
	}
	ll.Printf("%stag:%s connsess:%s starts_done\n",
		exampid, ltag, conn.Session())
	wgs.Wait()

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, ltag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s on_disconnect error:%v",
			exampid, ltag, conn.Session(),
			e.Error()) // Handle this ......
	}

	sngecomm.ShowStats(exampid, ltag, conn)
	wga.Done()
}

func startReceivers(qn int) {
	ltag := tag + "-startreceivers"

	ll.Printf("%stag:%s connsess:%s starts qn:%d\n",
		exampid, ltag, sngecomm.Lcs,
		qn)

	// Standard example connect sequence
	n, conn, e := sngecomm.CommonConnect(exampid, ltag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s on_connect error:%v",
			exampid, ltag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	nmsgs := senv.Nmsgs() // get message count
	ll.Printf("%stag:%s connsess:%s message_count nmsgs:%d qn:%d\n",
		exampid, ltag, conn.Session(),
		nmsgs, qn)
	for i := 1; i <= qn; i++ { // all queues
		wgr.Add(1)
		go receiver(conn, i, nmsgs)
	}
	ll.Printf("%stag:%s connsess:%s starts_done\n",
		exampid, ltag, conn.Session())
	wgr.Wait()

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, ltag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s on_disconnect error:%v",
			exampid, ltag, conn.Session(),
			e.Error()) // Handle this ......
	}

	sngecomm.ShowStats(exampid, ltag, conn)
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
	//
	sw = sngecomm.SendWait()
	rw = sngecomm.RecvWait()
	sf = sngecomm.SendFactor()
	rf = sngecomm.RecvFactor()
	ll.Printf("%stag:%s connsess:%s main_wait_sleep_factors sw:%v rw:%v sf:%v rf:%v\n",
		exampid, tag, sngecomm.Lcs,
		sw, rw, sf, rf)
	//
	q := sngecomm.Nqs()
	//
	wga.Add(2)
	go startReceivers(q)
	go startSenders(q)
	wga.Wait()

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, sngecomm.Lcs,
		time.Now().Sub(st))

}
