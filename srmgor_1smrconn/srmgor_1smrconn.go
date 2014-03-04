//
// Copyright © 2014 Guy M. Allard
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
	"runtime"
	"sync"
	"time"
	//
	"github.com/davecheney/profile"
	"github.com/gmallard/stompngo"
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
	sendWait = true
	recvWait = true

	// Sleep multipliers
	sendFact float64 = 1.0
	recvFact float64 = 1.0

	lhl = 44

	wgsend sync.WaitGroup
	wgrecv sync.WaitGroup
	wgall  sync.WaitGroup
)

/*
openSconn opens a stompngo Connection.
*/
func openSconn() (net.Conn, *stompngo.Connection) {
	// Open the net and stomp connection
	h, p := sngecomm.HostAndPort() // network connection host and port
	var e error
	// Network open
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "dial error", e) // Handle this ......
	}
	// Stomp connect, 1.1(+)
	ch := sngecomm.ConnectHeaders()
	log.Println(sngecomm.ExampIdNow(exampid), "vhost:", sngecomm.Vhost(), "protocol:", sngecomm.Protocol())
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "connect error", e) // Handle this ......
	}
	return n, conn
}

/*
closeSconn closes a stompngo Connection.
*/
func closeSconn(n net.Conn, conn *stompngo.Connection) {
	// Disconnect from Stomp server
	e := conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "startSender disconnect error", e) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "startSender netclose error", e) // Handle this ......
	}
	return
}

/*
runReceive receives all messages from a specified queue.
*/
func runReceive(conn *stompngo.Connection, q int, w *sync.WaitGroup) {
	qns := fmt.Sprintf("%d", q) // queue number
	id := stompngo.Uuid()       // A unique subscription ID
	fmt.Println(sngecomm.ExampIdNow(exampid), id, "runReceive starts", qns)
	//
	qp := sngecomm.Dest() // queue name prefix
	qname := qp + "." + qns
	fmt.Println(sngecomm.ExampIdNow(exampid), id, "runReceive queue name:", qname, qns)
	// Subscribe (use common helper)
	r := sngecomm.Subscribe(conn, qname, id, sngecomm.AckMode())
	//
	tmr := time.NewTimer(100 * time.Hour)

	// Receive loop
	for m := 1; m <= sngecomm.Nmsgs(); m++ {
		fmt.Println(sngecomm.ExampIdNow(exampid), id, "runReceive chanchek", "q", qns, "len", len(r), "cap", cap(r))
		d := <-r
		if d.Error != nil {
			log.Fatalln(sngecomm.ExampIdNow(exampid), id, "runReceive error", d.Error, qns)
		}

		// Process the inbound message .................
		osl := lhl
		if len(d.Message.Body) < osl {
			osl = len(d.Message.Body)
		}
		os := string(d.Message.Body[0:osl])
		fmt.Println(sngecomm.ExampIdNow(exampid), id, "runReceive message", os, qns, m)

		// Sanity check the message Command, and the queue and message numbers
		mns := fmt.Sprintf("%d", m) // message number
		if d.Message.Command != stompngo.MESSAGE {
			log.Fatalln("Bad Frame", d, qns, mns)
		}
		if !d.Message.Headers.ContainsKV("qnum", qns) || !d.Message.Headers.ContainsKV("msgnum", mns) {
			log.Fatalln("Bad Headers", d.Message.Headers, qns, mns)
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
		if m == sngecomm.Nmsgs() {
			break
		}
		if recvWait {
			d := time.Duration(sngecomm.ValueBetween(min, max, recvFact))
			fmt.Println(sngecomm.ExampIdNow(exampid), id, "runReceive", "stagger", int64(d)/1000000, "ms", qns)
			tmr.Reset(d)
			_ = <-tmr.C
			runtime.Gosched()
		}
	}
	// Unsubscribe
	sngecomm.Unsubscribe(conn, qname, id)

	fmt.Println(sngecomm.ExampIdNow(exampid), "runReceive", "ends", q)
	w.Done()
}

/*
receiverConnection starts individual receivers for this connection.
*/
func receiverConnection(conn *stompngo.Connection, cn, qpc int) {
	fmt.Println(sngecomm.ExampIdNow(exampid), "receiverConnection", "starts", cn)

	cb := cn - 1
	q1 := qpc*cb + 1
	ql := q1 + qpc - 1
	if ql > sngecomm.Nqs() {
		ql = sngecomm.Nqs()
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "receiverConnection", "startq", q1, "endq", ql)
	var wgrconn sync.WaitGroup
	for q := q1; q <= ql; q++ {
		wgrconn.Add(1)
		go runReceive(conn, q, &wgrconn)
	}
	wgrconn.Wait()
	//
	fmt.Println(sngecomm.ExampIdNow(exampid), "receiverConnection", "ends", cn)
	wgrecv.Done()
}

/*
startReceivers creates connections per environment variables, and starts each
connection.
*/
func startReceivers() {
	mr := sngecomm.Recvconns()
	if mr > sngecomm.Nqs() {
		mr = sngecomm.Nqs()
	}
	qpc := sngecomm.Nqs() / mr
	if sngecomm.Nqs()%mr != 0 {
		qpc += 1
	}
	if qpc == 0 {
		mr = 1
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "startReceivers", "mr", mr, "qpc", qpc)
	ncs := make([]net.Conn, 0)
	cs := make([]*stompngo.Connection, 0)
	for c := 1; c <= mr; c++ { // :-)
		n, conn := openSconn()
		ncs = append(ncs, n)
		cs = append(cs, conn)
		wgrecv.Add(1)
		go receiverConnection(conn, c, qpc)
	}
	wgrecv.Wait()
	fmt.Println(sngecomm.ExampIdNow(exampid), "startReceivers", "receives done")
	//
	for i := 0; i < len(ncs); i++ {
		closeSconn(ncs[i], cs[i])
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "startReceivers", "finishes")
	//
	wgall.Done()
}

/*
runSender sends all messages to a specified queue.
*/
func runSender(conn *stompngo.Connection, qns string) {
	qname := sngecomm.Dest() + "." + qns
	id := stompngo.Uuid() // A unique sender id
	fmt.Println(sngecomm.ExampIdNow(exampid), id, "send queue name:", qname)
	h := stompngo.Headers{"destination", qname, "senderId", id,
		"qnum", qns} // basic send Headers
	tmr := time.NewTimer(100 * time.Hour)
	for m := 1; m <= sngecomm.Nmsgs(); m++ {

		ms := fmt.Sprintf("%d", m)
		sh := append(h, "msgnum", ms)
		// Generate a message to send ...............
		fmt.Println(sngecomm.ExampIdNow(exampid), id, "send message", qns, ms)
		e := conn.Send(sh, string(sngecomm.Partial()))
		if e != nil {
			log.Fatalln(sngecomm.ExampIdNow(exampid), id, "send error", e, qns, ms)
		}
		if m == sngecomm.Nmsgs() {
			break
		}
		if sendWait {
			d := time.Duration(sngecomm.ValueBetween(min, max, sendFact))
			fmt.Println(sngecomm.ExampIdNow(exampid), id, "send", "stagger", int64(d)/1000000, "ms", qns, ms)
			tmr.Reset(d)
			_ = <-tmr.C
			runtime.Gosched()
		}
	}
	//
	wgsend.Done()
}

/*
startSender initializes the send connection, and runs senders for each queue.
*/
func startSender() {
	n, conn := openSconn()
	for q := 1; q <= sngecomm.Nqs(); q++ {
		wgsend.Add(1)
		go runSender(conn, fmt.Sprintf("%d", q))
	}
	wgsend.Wait()
	closeSconn(n, conn)
	//
	wgall.Done()
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
	fmt.Println(sngecomm.ExampIdNow(exampid), "main starts")
	fmt.Println(sngecomm.ExampIdNow(exampid), "main profiling", sngecomm.Pprof())
	fmt.Println(sngecomm.ExampIdNow(exampid), "main current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	if sngecomm.SetMAXPROCS() {
		nc := runtime.NumCPU()
		fmt.Println(sngecomm.ExampIdNow(exampid), "main number of CPUs is:", nc)
		c := runtime.GOMAXPROCS(nc)
		fmt.Println(sngecomm.ExampIdNow(exampid), "main previous number of GOMAXPROCS is:", c)
		fmt.Println(sngecomm.ExampIdNow(exampid), "main current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	}
	// Wait flags
	sendWait = sngecomm.SendWait()
	recvWait = sngecomm.RecvWait()
	sendFact = sngecomm.SendFactor()
	recvFact = sngecomm.RecvFactor()
	fmt.Println(sngecomm.ExampIdNow(exampid), "main Sleep Factors", "send", sendFact, "recv", recvFact)

	wgall.Add(1)
	go startReceivers()
	wgall.Add(1)
	go startSender()
	wgall.Wait()

	// The end
	dur := time.Since(start)
	fmt.Println(sngecomm.ExampIdNow(exampid), "main ends", dur)
}
