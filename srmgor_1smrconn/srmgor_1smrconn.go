//
// Copyright © 2014-2016 Guy M. Allard
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
// Receivers checks messages for proper queue and message number.

/*
Send and receive many STOMP messages using multiple queues and goroutines
to service each send or receive instance. All senders use one STOMP connection.
All receivers are balanced across multiple STOMP connections.  Balancing
configuration is taken from environment variables.
*/
package main

import (
	"fmt"
	"log"
	"net"
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
	exampid = "srmgor_1smrconn: "

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

	lhl = 44

	wgs sync.WaitGroup
	wgr sync.WaitGroup
	wga sync.WaitGroup

	ll = log.New(os.Stdout, "E1SMR ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	tag = "1smrconn"
)

/*
openSconn opens a stompngo Connection.
*/
func openSconn() (net.Conn, *stompngo.Connection) {
	ltag := tag + "-opensconn"

	// Standard example connect sequence
	n, conn, e := sngecomm.CommonConnect(exampid, ltag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s consess:%s connect_error error:%s\n",
			exampid, ltag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}
	return n, conn
}

/*
closeSconn closes a stompngo Connection.
*/
func closeSconn(n net.Conn, conn *stompngo.Connection) {
	ltag := tag + "-closesconn"

	// Standard example disconnect sequence
	e := sngecomm.CommonDisconnect(n, conn, exampid, ltag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s disconnect_error error:%s\n",
			exampid, ltag, conn.Session(),
			e.Error()) // Handle this ......
	}
	return
}

/*
runReceive receives all messages from a specified queue.
*/
func runReceive(conn *stompngo.Connection, q int, w *sync.WaitGroup) {
	ltag := tag + "-runreceive"

	qns := fmt.Sprintf("%d", q) // queue number
	id := stompngo.Uuid()       // A unique subscription ID
	d := sngecomm.Dest() + "." + string(exampid[:len(exampid)-2]) + "." + qns

	ll.Printf("%stag:%s connsess:%s starts id:%s qns:%s d:%s\n",
		exampid, ltag, conn.Session(),
		id, qns, d)

	// Subscribe (use common helper)
	sc := sngecomm.HandleSubscribe(conn, d, id, sngecomm.AckMode())
	ll.Printf("%stag:%s connsess:%s subscribe_done id:%s qns:%s d:%s\n",
		exampid, ltag, conn.Session(),
		id, qns, d)

	//
	tmr := time.NewTimer(100 * time.Hour)

	pbc := sngecomm.Pbc() // Print byte count

	nmsgs := senv.Nmsgs()

	// Receive loop
	var md stompngo.MessageData
	for mc := 1; mc <= nmsgs; mc++ {
		ll.Printf("%stag:%s connsess:%s chanchek id:%s qns:%s lensc:%d capsc:%d\n",
			exampid, ltag, conn.Session(),
			id, qns, len(sc), cap(sc))
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

		// Process the inbound message .................
		ll.Printf("%stag:%s connsess:%s inbound id:%s qns:%s mc:%d\n",
			exampid, ltag, conn.Session(),
			id, qns, mc)
		// Sanity check the message Command, and the queue and message numbers
		mns := fmt.Sprintf("%d", mc) // string message number
		if md.Message.Command != stompngo.MESSAGE {
			ll.Fatalf("%stag:%s connsess:%s bad_frame qns:%s mc:%d md:%v\n",
				exampid, ltag, conn.Session(),
				qns, mc, md)
		}
		if !md.Message.Headers.ContainsKV("qnum", qns) || !md.Message.Headers.ContainsKV("msgnum", mns) {
			ll.Fatalf("%stag:%s connsess:%s dirty_message qns:%v msgnum:%v md:%v",
				exampid, tag, conn.Session(),
				qns, mns, md) // Handle this ......
		}

		sl := len(md.Message.Body)
		if pbc > 0 {
			sl = pbc
			if len(md.Message.Body) < sl {
				sl = len(md.Message.Body)
			}
		}

		ll.Printf("%stag:%s connsess:%s runReceive_recv_message id:%s body:%s qns:%s msgnum:%s\n",
			exampid, ltag, conn.Session(),
			id, string(md.Message.Body[0:sl]), qns,
			md.Message.Headers.Value("msgnum"))

		// Handle ACKs if needed
		if sngecomm.AckMode() != "auto" {
			ah := stompngo.Headers{}
			sngecomm.HandleAck(conn, ah, id)
		}
		if mc == nmsgs {
			break
		}
		if rw {
			dt := time.Duration(sngecomm.ValueBetween(min, max, rf))
			ll.Printf("%stag:%s connsess:%s recv_stagger dt:%v qns:%s mc:%d\n",
				exampid, ltag, conn.Session(),
				dt, qns, mc)
			tmr.Reset(dt)
			_ = <-tmr.C
			runtime.Gosched()
		}
	}
	// Unsubscribe
	sngecomm.HandleUnsubscribe(conn, d, id)

	ll.Printf("%stag:%s connsess:%s runRecieve_ends id:%s qns:%s\n",
		exampid, ltag, conn.Session(),
		id, qns)
	w.Done()
}

/*
receiverConnection starts individual receivers for this connection.
*/
func receiverConnection(conn *stompngo.Connection, cn, qpc int) {
	ltag := tag + "-receiverconnection"

	ll.Printf("%stag:%s connsess:%s starts cn:%d qpc:%d\n",
		exampid, ltag, conn.Session(),
		cn, qpc)

	// cn -> a connection number: 1..n
	// qpc -> destinations per connection
	// Ex:
	// 1, 2
	// 2, 2
	// 3, 2

	// This code runs *once* for each connection

	// These calcs are what causes a skip below.  It is a safety valve to keep
	// from starting one too many connections.
	cb := cn - 1       // this connection number, zero based
	q1 := qpc*cb + 1   // 1st queue number
	ql := q1 + qpc - 1 // last queue number
	if ql > sngecomm.Nqs() {
		ql = sngecomm.Nqs() // truncate last if over max destinations
	}

	var wgrconn sync.WaitGroup

	var skipped bool
	if q1 <= ql {
		ll.Printf("%stag:%s connsess:%s startq cn:%d q1:%d ql: %d\n",
			exampid, ltag, conn.Session(),
			cn, q1, ql)
		skipped = false
	} else {
		// Skips are possible, at least with the current calling code, see above
		ll.Printf("%stag:%s connsess:%s startskip cn:%d q1:%d ql: %d\n",
			exampid, ltag, conn.Session(),
			cn, q1, ql)
		skipped = true
	}

	for q := q1; q <= ql; q++ {
		wgrconn.Add(1)
		go runReceive(conn, q, &wgrconn)
	}
	wgrconn.Wait()
	//
	ll.Printf("%stag:%s connsess:%s ends cn:%d qpc:%d skipped:%t\n",
		exampid, ltag, conn.Session(),
		cn, qpc, skipped)
	wgr.Done()
}

/*
startReceivers creates connections per environment variables, and starts each
connection.
*/
func startReceivers() {

	ltag := tag + "-startreceivers"

	// This was a performance experiment.  With number of connections.
	// My recollection is that it did not work out.
	// However ..... I will leave this code in place for now.

	// Figure out number of receiver connections wanted
	nrc := sngecomm.Nqs() // 1 receiver per each destination
	nqs := nrc            // Number of queues (destinations) starts the same

	if s := os.Getenv("STOMP_RECVCONNS"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			ll.Fatalf("%stag:%s connsess:%s RECVCONNS_conversion_error error:%v\n",
				exampid, ltag, sngecomm.Lcs,
				e.Error())
		} else {
			nrc = int(i)
		}
	}

	// Limit max receiver connection count to number of destinations
	if nrc > nqs {
		nrc = nqs
	}

	// Next calc. destinations per receiver
	dpr := nqs / nrc // Calculation first guess.
	if nqs%nrc != 0 {
		dpr += 1 // Bump destinations per receiver by 1.
	}
	// Destinations per receiver must be at least 1
	if dpr == 0 {
		dpr = 1
	}

	ll.Printf("%stag:%s connsess:%s start nrc:%d dpr:%d\n",
		exampid, ltag, sngecomm.Lcs,
		nrc, dpr)

	// So the idea seems to be allow more than one destination per receiver
	ncm := make([]net.Conn, 0)
	csm := make([]*stompngo.Connection, 0)
	for c := 1; c <= nrc; c++ { // :-)
		n, conn := openSconn()
		ncm = append(ncm, n)
		csm = append(csm, conn)
		wgr.Add(1)
		ll.Printf("%stag:%s connsess:%s connstart conn_number:%d nrc:%d dpr:%d\n",
			exampid, ltag, conn.Session(),
			c, nrc, dpr)
		go receiverConnection(conn, c, dpr)
	}
	wgr.Wait()
	ll.Printf("%stag:%s connsess:%s wait_done nrc:%d dpr:%d\n",
		exampid, ltag, sngecomm.Lcs,
		nrc, dpr)
	//
	for c := 1; c <= nrc; c++ {
		ll.Printf("%stag:%s connsess:%s connend conn_number:%d nrc:%d dpr:%d\n",
			exampid, ltag, csm[c-1].Session(),
			c, nrc, dpr)
		sngecomm.ShowStats(exampid, ltag, csm[c-1])
		closeSconn(ncm[c-1], csm[c-1])
	}
	//
	wga.Done()
}

/*
runSender sends all messages to a specified queue.
*/
func runSender(conn *stompngo.Connection, qns string) {
	ltag := tag + "-runsender"

	d := sngecomm.Dest() + "."  + string(exampid[:len(exampid)-2]) + "." + qns
	id := stompngo.Uuid() // A unique sender id
	ll.Printf("%stag:%s connsess:%s start id:%s dest:%s\n",
		exampid, ltag, conn.Session(),
		id, d)
	wh := stompngo.Headers{"destination", d, "senderId", id,
		"qnum", qns} // basic send Headers
	if senv.Persistent() {
		wh = wh.Add("persistent", "true")
	}
	tmr := time.NewTimer(100 * time.Hour)
	nmsgs := senv.Nmsgs()
	for mc := 1; mc <= nmsgs; mc++ {
		sh := append(wh, "msgnum", fmt.Sprintf("%d", mc))
		// Generate a message to send ...............
		ll.Printf("%stag:%s  connsess:%s send id:%s qns:%s mc:%d\n",
			exampid, ltag, conn.Session(),
			id, qns, mc)
		e := conn.Send(sh, string(sngecomm.Partial()))
		if e != nil {
			ll.Fatalf("%stag:%s connsess:%s send_error qns:%v error:%v",
				exampid, ltag, conn.Session(),
				qns, e.Error()) // Handle this ......
		}
		if mc == nmsgs {
			break
		}
		if sw {
			dt := time.Duration(sngecomm.ValueBetween(min, max, sf))
			ll.Printf("%stag:%s connsess:%s send_stagger dt:%v qns:%s mc:%d\n",
				exampid, ltag, conn.Session(),
				dt, qns, mc)
			tmr.Reset(dt)
			_ = <-tmr.C
			runtime.Gosched()
		}
	}
	ll.Printf("%stag:%s connsess:%s end id:%s dest:%s\n",
		exampid, ltag, conn.Session(),
		id, d)
	//
	wgs.Done()
}

/*
startSender initializes the single send connection, and starts one sender go
for each destination.
*/
func startSender() {
	ltag := tag + "-startsender"

	n, conn := openSconn()
	ll.Printf("%stag:%s connsess:%s start\n",
		exampid, ltag, conn.Session())
	for i := 1; i <= sngecomm.Nqs(); i++ {
		wgs.Add(1)
		go runSender(conn, fmt.Sprintf("%d", i))
	}
	wgs.Wait()
	ll.Printf("%stag:%s connsess:%s end\n",
		exampid, ltag, conn.Session())
	sngecomm.ShowStats(exampid, ltag, conn)
	closeSconn(n, conn)
	//
	wga.Done()
}

/*
main is the driver for all logic.
*/
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

	wga.Add(1)
	go startReceivers()
	wga.Add(1)
	go startSender()
	wga.Wait()

	// The end
	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, sngecomm.Lcs,
		time.Now().Sub(st))
}
