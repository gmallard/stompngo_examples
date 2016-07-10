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
	"github.com/davecheney/profile"
	//
	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "srmgor_1smrconn:"

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
)

/*
openSconn opens a stompngo Connection.
*/
func openSconn() (net.Conn, *stompngo.Connection) {
	// Open the net and stomp connection
	h, p := senv.HostAndPort() // network connection host and port
	var e error
	// Network open
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalln(exampid, "dial error", e) // Handle this ......
	}
	// Stomp connect, 1.1(+)
	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalln(exampid, "connect error", e) // Handle this ......
	}
	ll.Printf("%s connsess:%s vhost:%s protocol:%s\n", exampid, conn.Session(),
		senv.Vhost(), senv.Protocol())
	return n, conn
}

/*
closeSconn closes a stompngo Connection.
*/
func closeSconn(n net.Conn, conn *stompngo.Connection) {
	// Disconnect from Stomp server
	e := conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalln(exampid, "startSender disconnect error", e) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		ll.Fatalln(exampid, "startSender netclose error", e) // Handle this ......
	}
	return
}

/*
runReceive receives all messages from a specified queue.
*/
func runReceive(conn *stompngo.Connection, q int, w *sync.WaitGroup) {
	qns := fmt.Sprintf("%d", q) // queue number
	id := stompngo.Uuid()       // A unique subscription ID
	ll.Printf("%s connsess:%s runRecieve_starts id:%s qns:%s\n",
		exampid, conn.Session(), id, qns)
	//
	d := senv.Dest() + "." + qns
	ll.Printf("%s connsess:%s runRecieve_starts_dest id:%s d:%s qns:%s\n",
		exampid, conn.Session(), id, d, qns)
	// Subscribe (use common helper)
	sc := sngecomm.HandleSubscribe(conn, d, id, sngecomm.AckMode())
	//
	tmr := time.NewTimer(100 * time.Hour)

	pbc := sngecomm.Pbc() // Print byte count

	nmsgs := senv.Nmsgs()

	// Receive loop
	var md stompngo.MessageData
	for mc := 1; mc <= nmsgs; mc++ {
		ll.Printf("%s connsess:%s runReceive_chanchek id:%s qns:%s lensc:%d capsc:%d\n",
			exampid, conn.Session(), id, qns, len(sc), cap(sc))

		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// Frames RECEIPT or ERROR not expected here
			ll.Fatalln(exampid, md) // Handle this
		}

		if md.Error != nil {
			ll.Fatalln(exampid, id, "runReceive error", md.Error, qns)
		}

		// Process the inbound message .................
		ll.Printf("%s connsess:%s runReceive_inbound id:%s qns:%s mc:%d\n",
			exampid, conn.Session(),
			id, qns, mc)

		// Sanity check the message Command, and the queue and message numbers
		mns := fmt.Sprintf("%d", mc) // string message number
		if md.Message.Command != stompngo.MESSAGE {
			ll.Fatalln("Bad Frame", md, qns, mns)
		}
		if !md.Message.Headers.ContainsKV("qnum", qns) || !md.Message.Headers.ContainsKV("msgnum", mns) {
			ll.Fatalln("Bad Headers", md.Message.Headers, qns, mns)
		}

		sl := len(md.Message.Body)
		if pbc > 0 {
			sl = pbc
			if len(md.Message.Body) < sl {
				sl = len(md.Message.Body)
			}
		}

		ll.Printf("%s connsess:%s runReceive_recv_message id:%s body:%s qns:%s msgnum:%s\n",
			exampid, conn.Session(), id, string(md.Message.Body[0:sl]), qns,
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
			ll.Printf("%s connsess:%s runReceive_stagger dt:%v qns:%s mc:%d\n",
				exampid, conn.Session(),
				dt, qns, mc)
			tmr.Reset(dt)
			_ = <-tmr.C
			runtime.Gosched()
		}
	}
	// Unsubscribe
	sngecomm.HandleUnsubscribe(conn, d, id)

	ll.Printf("%s connsess:%s runRecieve_ends id:%s qns:%s\n",
		exampid, conn.Session(), id, qns)
	w.Done()
}

/*
receiverConnection starts individual receivers for this connection.
*/
func receiverConnection(conn *stompngo.Connection, cn, qpc int) {
	ll.Printf("%s connsess:%s receiverConnection_starts cn:%d qpc:%d\n",
		exampid, conn.Session(), cn, qpc)

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
		ll.Printf("%s connsess:%s receiverConnection_startq cn:%d q1:%d ql: %d\n",
			exampid, conn.Session(), cn, q1, ql)
		skipped = false
	} else {
		// Skips are possible, at least with the current calling code, see above
		ll.Printf("%s connsess:%s receiverConnection_startskip cn:%d q1:%d ql: %d\n",
			exampid, conn.Session(), cn, q1, ql)
		skipped = true
	}

	for q := q1; q <= ql; q++ {
		wgrconn.Add(1)
		go runReceive(conn, q, &wgrconn)
	}
	wgrconn.Wait()
	//
	ll.Printf("%s connsess:%s receiverConnection_ends cn:%d qpc:%d skipped:%t\n",
		exampid, conn.Session(), cn, qpc, skipped)
	wgr.Done()
}

/*
startReceivers creates connections per environment variables, and starts each
connection.
*/
func startReceivers() {

	// This was a performance experiment.  With number of connections.
	// My recollection is that it did not work out.
	// However ..... I will leave this code in place for now.

	// Figure out number of receiver connections wanted
	nrc := sngecomm.Nqs() // 1 receiver per each destination
	nqs := nrc            // Number of queues (destinations) starts the same

	// ll.Println("SRDBG00", nrc, nqs)

	if s := os.Getenv("STOMP_RECVCONNS"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			ll.Println("RECVCONNS conversion error", e)
		} else {
			nrc = int(i)
		}
	}

	// ll.Println("SRDBG01", nrc, nqs)

	// Limit max receiver connection count to number of destinations
	if nrc > nqs {
		nrc = nqs
	}

	//ll.Println("SRDBG02", nrc, nqs)

	// Next calc. destinations per receiver
	dpr := nqs / nrc // Calculation first guess.
	//ll.Println("SRDBG03", nrc, nqs, dpr)
	if nqs%nrc != 0 {
		dpr += 1 // Bump destinations per receiver by 1.
		// ll.Println("SRDBG04", nrc, nqs, dpr)
	}
	// Destinations per receiver must be at least 1
	if dpr == 0 {
		dpr = 1
		// ll.Println("SRDBG05", nrc, nqs, dpr)
	}

	// ll.Println("SRDBG06", nrc, nqs, dpr)

	ll.Printf("%s startReceivers_start nrc:%d dpr:%d\n",
		exampid, nrc, dpr)

	// So the idea seems to be allow more than one destination per receiver
	ncm := make([]net.Conn, 0)
	csm := make([]*stompngo.Connection, 0)
	for c := 1; c <= nrc; c++ { // :-)
		n, conn := openSconn()
		ncm = append(ncm, n)
		csm = append(csm, conn)
		wgr.Add(1)
		ll.Printf("%s connsess:%s startReceivers_connnew conn_number:%d nrc:%d dpr:%d\n",
			exampid, conn.Session(), c, nrc, dpr)
		go receiverConnection(conn, c, dpr)
	}
	wgr.Wait()
	ll.Printf("%s startReceivers_wait_done nrc:%d dpr:%d\n",
		exampid, nrc, dpr)
	//
	for c := 1; c <= nrc; c++ {
		ll.Printf("%s connsess:%s startReceivers_connend conn_number:%d nrc:%d dpr:%d\n",
			exampid, csm[c-1].Session(), c, nrc, dpr)
		closeSconn(ncm[c-1], csm[c-1])
	}
	//
	wga.Done()
}

/*
runSender sends all messages to a specified queue.
*/
func runSender(conn *stompngo.Connection, qns string) {
	d := senv.Dest() + "." + qns
	id := stompngo.Uuid() // A unique sender id
	ll.Printf("%s runSender_start connsess:%s id:%s dest:%s\n",
		exampid, conn.Session(), id, d)
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
		ll.Printf("%s runSender_send connsess:%s id:%s qns:%s mc:%d\n",
			exampid, conn.Session(),
			id, qns, mc)
		e := conn.Send(sh, string(sngecomm.Partial()))
		if e != nil {
			ll.Fatalln(exampid, id, "send error", e, qns, mc)
		}
		if mc == nmsgs {
			break
		}
		if sw {
			dt := time.Duration(sngecomm.ValueBetween(min, max, sf))
			ll.Printf("%s runSender_stagger connsess:%s id:%s send_stagger:%v qns:%s mc:%d\n",
				exampid, conn.Session(), id,
				dt, qns, mc)
			tmr.Reset(dt)
			_ = <-tmr.C
			runtime.Gosched()
		}
	}
	ll.Printf("%s runSender_end connsess:%s id:%s dest:%s\n",
		exampid, conn.Session(), id, d)
	//
	wgs.Done()
}

/*
startSender initializes the single send connection, and starts one sender go
for each destination.
*/
func startSender() {
	n, conn := openSconn()
	ll.Printf("%s startSender_starts 1connsess:%s\n", exampid, conn.Session())
	for i := 1; i <= sngecomm.Nqs(); i++ {
		wgs.Add(1)
		go runSender(conn, fmt.Sprintf("%d", i))
	}
	wgs.Wait()
	ll.Printf("%s startSender_ends 1connsess:%s\n", exampid, conn.Session())
	closeSconn(n, conn)
	//
	wga.Done()
}

/*
main is the driver for all logic.
*/
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
	ll.Println(exampid, "main starts")
	ll.Println(exampid, "main profiling", sngecomm.Pprof())
	ll.Println(exampid, "main current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	if sngecomm.SetMAXPROCS() {
		nc := runtime.NumCPU()
		ll.Println(exampid, "main number of CPUs is:", nc)
		c := runtime.GOMAXPROCS(nc)
		ll.Println(exampid, "main previous number of GOMAXPROCS is:", c)
		ll.Println(exampid, "main current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	}
	// Wait flags
	sw = sngecomm.SendWait()
	rw = sngecomm.RecvWait()
	sf = sngecomm.SendFactor()
	rf = sngecomm.RecvFactor()
	ll.Println(exampid, "main Sleep Factors", "send", sf, "recv", rf)

	wga.Add(1)
	go startReceivers()
	wga.Add(1)
	go startSender()
	wga.Wait()

	// The end
	dt := time.Since(start)
	ll.Println(exampid, "main ends", dt)
}
