//
// Copyright © 2012-2014 Guy M. Allard
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
	"runtime"
	"sync"
	"time"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var exampid = "srmgor_manyconn:"

var wgsend sync.WaitGroup
var wgrecv sync.WaitGroup

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

// Number of messages
var nmsgs = 1

func sendMessages(conn *stompngo.Connection, qnum int, nc net.Conn) {
	qns := fmt.Sprintf("%d", qnum) // queue number
	qp := sngecomm.Dest()          // queue name prefix
	q := qp + "." + qns
	fmt.Println(sngecomm.ExampIdNow(exampid), "send queue name:", q, qnum)
	h := stompngo.Headers{"destination", q,
		"qnum", qns} // send Headers
	if sngecomm.Persistent() {
		h = h.Add("persistent", "true")
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "send starts", nmsgs, qnum)
	//
	tmr := time.NewTimer(100 * time.Hour)
	// Send messages
	for n := 1; n <= nmsgs; n++ {
		si := fmt.Sprintf("%d", n)
		sh := append(h, "msgnum", si)
		// Generate a message to send ...............
		fmt.Println(sngecomm.ExampIdNow(exampid), "send message", qnum, si)
		e := conn.Send(sh, string(sngecomm.Partial()))
		if e != nil {
			log.Fatalln(sngecomm.ExampIdNow(exampid), "send:", e, nc.LocalAddr().String(), qnum)
		}
		if n == nmsgs {
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
}

func receiveMessages(conn *stompngo.Connection, qnum int, nc net.Conn) {
	qns := fmt.Sprintf("%d", qnum) // queue number
	qp := sngecomm.Dest()          // queue name prefix
	q := qp + "." + qns
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv queue name:", q, qnum)
	// Subscribe
	id := stompngo.Uuid() // A unique subscription ID
	r := sngecomm.Subscribe(conn, q, id, "auto")
	// Receive messages
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv starts", nmsgs, qnum)
	//
	tmr := time.NewTimer(100 * time.Hour)
	for n := 1; n <= nmsgs; n++ {
		d := <-r
		if d.Error != nil {
			log.Fatalln(sngecomm.ExampIdNow(exampid), "recv read:", d.Error, nc.LocalAddr().String(), qnum)
		}

		// Sanity check the queue and message numbers
		mns := fmt.Sprintf("%d", n) // message number
		if !d.Message.Headers.ContainsKV("qnum", qns) || !d.Message.Headers.ContainsKV("msgnum", mns) {
			log.Fatalln("Bad Headers", d.Message.Headers, qns, mns)
		}

		// Process the inbound message .................
		sl := 16
		if len(d.Message.Body) < sl {
			sl = len(d.Message.Body)
		}
		fmt.Println(sngecomm.ExampIdNow(exampid), "recv message", string(d.Message.Body[0:sl]), qnum, d.Message.Headers.Value("msgnum"))
		if n == nmsgs {
			break
		}
		//
		if recv_wait {
			runtime.Gosched() // yield for this example
			d := time.Duration(sngecomm.ValueBetween(min, max, recvFact))
			fmt.Println(sngecomm.ExampIdNow(exampid), "recv", "stagger", int64(d)/1000000, "ms")
			tmr.Reset(d)
			_ = <-tmr.C
		}
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv done:", q)
	// Unsubscribe
	sngecomm.Unsubscribe(conn, q, id)
	//
}

func runReceiver(qnum int) {
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv start for queue number", qnum)
	// Network Open
	h, p := sngecomm.HostAndPort() // host and port
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "recv nectonnr:", qnum, e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv network open complete", qnum)
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv network local", n.LocalAddr().String(), qnum)
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv network remote", n.RemoteAddr().String(), qnum)
	// Stomp connect
	ch := sngecomm.ConnectHeaders()
	log.Println(sngecomm.ExampIdNow(exampid), "recv", "vhost:", sngecomm.Vhost(), "protocol:", sngecomm.Protocol())
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "recv stompconnect:", qnum, e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv connection complete:", qnum)
	//
	conn.SetSubChanCap(sngecomm.SubChanCap()) // Experiment with this value, YMMV
	// Receives
	receiveMessages(conn, qnum, n)
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv receives complete:", qnum)
	// Disconnect from Stomp server
	eh := stompngo.Headers{"recv_discqueue", fmt.Sprintf("%d", qnum)}
	e = conn.Disconnect(eh)
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "recv disconnects:", qnum, e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv disconnected:", qnum)
	// Network close
	e = n.Close()
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "recv netcloser", qnum, e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv network close complete", qnum)
	fmt.Println(sngecomm.ExampIdNow(exampid), "recv end for queue number", qnum)
	sngecomm.ShowStats(exampid, "recv "+fmt.Sprintf("%d", qnum), conn)
	wgrecv.Done()
}

func runSender(qnum int) {
	fmt.Println(sngecomm.ExampIdNow(exampid), "send start for queue number", qnum)
	// Network Open
	h, p := sngecomm.HostAndPort() // host and port
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "send nectonnr:", qnum, e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "send network open complete", qnum)
	fmt.Println(sngecomm.ExampIdNow(exampid), "send network local", n.LocalAddr().String(), qnum)
	fmt.Println(sngecomm.ExampIdNow(exampid), "send network remote", n.RemoteAddr().String(), qnum)
	// Stomp connect
	ch := sngecomm.ConnectHeaders()
	log.Println(sngecomm.ExampIdNow(exampid), "send", "vhost:", sngecomm.Vhost(), "protocol:", sngecomm.Protocol())
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "send stompconnect:", qnum, e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "send connection complete:", qnum)
	//
	sendMessages(conn, qnum, n)
	fmt.Println(sngecomm.ExampIdNow(exampid), "send sends complete:", qnum)
	// Disconnect from Stomp server
	eh := stompngo.Headers{"send_discqueue", fmt.Sprintf("%d", qnum)}
	e = conn.Disconnect(eh)
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "send disconnects:", qnum, e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "send disconnected:", qnum)
	// Network close
	e = n.Close()
	if e != nil {
		log.Fatalln(sngecomm.ExampIdNow(exampid), "send netcloser", qnum, e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "send network close complete", qnum)
	fmt.Println(sngecomm.ExampIdNow(exampid), "send end for queue number", qnum)
	sngecomm.ShowStats(exampid, "send "+fmt.Sprintf("%d", qnum), conn)
	wgsend.Done()
}

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
	numq := sngecomm.Nqs()
	nmsgs = sngecomm.Nmsgs() // message count
	//
	fmt.Println(sngecomm.ExampIdNow(exampid), "main starting receivers")
	for q := 1; q <= numq; q++ {
		wgrecv.Add(1)
		go runReceiver(q)
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "main started receivers")
	//
	fmt.Println(sngecomm.ExampIdNow(exampid), "main starting senders")
	for q := 1; q <= numq; q++ {
		wgsend.Add(1)
		go runSender(q)
	}
	fmt.Println(sngecomm.ExampIdNow(exampid), "main started senders")
	//
	wgsend.Wait()
	fmt.Println(sngecomm.ExampIdNow(exampid), "main senders complete")
	wgrecv.Wait()
	fmt.Println(sngecomm.ExampIdNow(exampid), "main receivers complete")
	//
	fmt.Println(sngecomm.ExampIdNow(exampid), "main ends", time.Since(tn))
}
