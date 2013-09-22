//
// Copyright © 2012-2013 Guy M. Allard
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
*/
package main

import (
	"fmt"
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
	"log"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
)

var exampid = "srmgor_11c:"

var wgsend sync.WaitGroup
var wgrecv sync.WaitGroup

// We 'stagger' between each message send and message receive for a random
// amount of time.
// Vary these for experimental purposes.  YMMV.
var max int64 = 1e9      // Max stagger time (nanoseconds)
var min int64 = max / 10 // Min stagger time (nanoseconds)
// Vary these for experimental purposes.  YMMV.
var send_factor int64 = 1 // Send factor time
var recv_factor int64 = 1 // Receive factor time

// Wait flags
var send_wait = true
var recv_wait = true

// Number of messages
var nmsgs = 1

func sendMessages(conn *stompngo.Connection, qnum int, nc net.Conn) {
	qns := fmt.Sprintf("%d", qnum) // queue number
	qp := sngecomm.Dest()          // queue name prefix
	q := qp + "." + qns
	fmt.Println(exampid, "send queue name:", q, qnum)
	h := stompngo.Headers{"destination", q} // send Headers
	if sngecomm.Persistent() {
		h = h.Add("persistent", "true")
	}
	fmt.Println(exampid, "send starts", nmsgs, qnum)
	// Send messages
	for n := 1; n <= nmsgs; n++ {
		si := fmt.Sprintf("%d", n)
		// Generate a message to send ...............
		mp := exampid + "|" + "payload" + "|qnum:" + qns + "|msgnum:" + si + " :"
		m := mp + sngecomm.Partial() // Variable length
		fmt.Println(exampid, "send message", mp, qnum)
		e := conn.Send(h, m)
		if e != nil {
			log.Fatalln(exampid, "send:", e, nc.LocalAddr().String(), qnum)
		}
		if send_wait {
			runtime.Gosched()                                                              // yield for this example
			time.Sleep(time.Duration(send_factor * (sngecomm.ValueBetween(min, max) / 2))) // Time to build next message
		}
	}
}

func receiveMessages(conn *stompngo.Connection, qnum int, nc net.Conn) {
	qns := fmt.Sprintf("%d", qnum) // queue number
	qp := sngecomm.Dest()          // queue name prefix
	q := qp + "." + qns
	fmt.Println(exampid, "recv queue name:", q, qnum)
	// Subscribe
	u := stompngo.Uuid() // A unique subscription ID
	h := stompngo.Headers{"destination", q, "id", u}
	r, e := conn.Subscribe(h)
	if e != nil {
		log.Fatalln(exampid, "recv subscribe", e, nc.LocalAddr().String(), qnum)
	}
	// Receive messages
	fmt.Println(exampid, "recv starts", nmsgs, qnum)
	for n := 1; n <= nmsgs; n++ {
		d := <-r
		if d.Error != nil {
			log.Fatalln(exampid, "recv read:", d.Error, nc.LocalAddr().String(), qnum)
		}

		// Process the inbound message .................
		m := d.Message.BodyString()
		li := strings.LastIndex(m, ":")
		fmt.Println(exampid, "recv message", string(m[0:li]), qnum)

		// Sanity check the queue and message numbers
		mns := fmt.Sprintf("%d", n) // message number
		t := "|qnum:" + qns + "|msgnum:" + mns
		if !strings.Contains(m, t) {
			log.Fatalln(exampid, "recv bad message", m, t, nc.LocalAddr().String(), qnum)
		}
		if recv_wait {
			runtime.Gosched()                                                              // yield for this example
			time.Sleep(time.Duration(send_factor * (sngecomm.ValueBetween(min, max) / 2))) // Time to build next message
		}
	}
	// Unsubscribe
	e = conn.Unsubscribe(h)
	if e != nil {
		log.Fatalln(exampid, "recv unsubscribe", e, nc.LocalAddr().String(), qnum)
	}
	//
}

func runReceiver(qnum int) {
	fmt.Println(exampid, "recv start for queue number", qnum)
	// Network Open
	h, p := sngecomm.HostAndPort11() // a 1.1(+) connect
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(exampid, "recv nectonnr:", qnum, e) // Handle this ......
	}
	fmt.Println(exampid, "recv network open complete", qnum)
	fmt.Println(exampid, "recv network local", n.LocalAddr().String(), qnum)
	fmt.Println(exampid, "recv network remote", n.RemoteAddr().String(), qnum)
	// Stomp connect, 1.1(+)
	ch := stompngo.Headers{"host", sngecomm.Vhost(),
		"accept-version", sngecomm.Protocol()}
	log.Println(exampid, "recv", "vhost:", sngecomm.Vhost(), "protocol:", sngecomm.Protocol())
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(exampid, "recv stompconnect:", qnum, e) // Handle this ......
	}
	fmt.Println(exampid, "recv connection complete:", qnum)
	// Receives
	receiveMessages(conn, qnum, n)
	fmt.Println(exampid, "recv receives complete:", qnum)
	// Disconnect from Stomp server
	eh := stompngo.Headers{"recv_discqueue", fmt.Sprintf("%d", qnum)}
	e = conn.Disconnect(eh)
	if e != nil {
		log.Fatalln(exampid, "recv disconnects:", qnum, e) // Handle this ......
	}
	fmt.Println(exampid, "recv disconnected:", qnum)
	// Network close
	e = n.Close()
	if e != nil {
		log.Fatalln(exampid, "recv netcloser", qnum, e) // Handle this ......
	}
	fmt.Println(exampid, "recv network close complete", qnum)
	fmt.Println(exampid, "recv end for queue number", qnum)
	sngecomm.ShowStats(exampid, "recv "+fmt.Sprintf("%d", qnum), conn)
	wgrecv.Done()
}

func runSender(qnum int) {
	fmt.Println(exampid, "send start for queue number", qnum)
	// Network Open
	h, p := sngecomm.HostAndPort11() // a 1.1(+) connect
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(exampid, "send nectonnr:", qnum, e) // Handle this ......
	}
	fmt.Println(exampid, "send network open complete", qnum)
	fmt.Println(exampid, "send network local", n.LocalAddr().String(), qnum)
	fmt.Println(exampid, "send network remote", n.RemoteAddr().String(), qnum)
	// Stomp connect, 1.1(+)
	ch := stompngo.Headers{"host", sngecomm.Vhost(),
		"accept-version", sngecomm.Protocol()}
	log.Println(exampid, "send", "vhost:", sngecomm.Vhost(), "protocol:", sngecomm.Protocol())
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(exampid, "send stompconnect:", qnum, e) // Handle this ......
	}
	fmt.Println(exampid, "send connection complete:", qnum)
	//
	sendMessages(conn, qnum, n)
	fmt.Println(exampid, "send sends complete:", qnum)
	// Disconnect from Stomp server
	eh := stompngo.Headers{"send_discqueue", fmt.Sprintf("%d", qnum)}
	e = conn.Disconnect(eh)
	if e != nil {
		log.Fatalln(exampid, "send disconnects:", qnum, e) // Handle this ......
	}
	fmt.Println(exampid, "send disconnected:", qnum)
	// Network close
	e = n.Close()
	if e != nil {
		log.Fatalln(exampid, "send netcloser", qnum, e) // Handle this ......
	}
	fmt.Println(exampid, "send network close complete", qnum)
	fmt.Println(exampid, "send end for queue number", qnum)
	sngecomm.ShowStats(exampid, "send "+fmt.Sprintf("%d", qnum), conn)
	wgsend.Done()
}

func main() {
	fmt.Println(exampid, "starts")
	if sngecomm.SetMAXPROCS() {
		nc := runtime.NumCPU()
		fmt.Println(exampid, "number of CPUs is:", nc)
		c := runtime.GOMAXPROCS(nc)
		fmt.Println(exampid, "previous number of GOMAXPROCS is:", c)
		fmt.Println(exampid, "current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	}
	//
	send_wait = sngecomm.SendWait()
	recv_wait = sngecomm.RecvWait()
	//
	numq := sngecomm.Nqs()
	fmt.Println(exampid, "numq:", numq)
	nmsgs = sngecomm.Nmsgs() // message count
	fmt.Println(exampid, "nmsgs:", nmsgs)
	//
	fmt.Println(exampid, "starting receivers")
	for q := 1; q <= numq; q++ {
		wgrecv.Add(1)
		go runReceiver(q)
	}
	fmt.Println(exampid, "started receivers")
	//
	fmt.Println(exampid, "starting senders")
	for q := 1; q <= numq; q++ {
		wgsend.Add(1)
		go runSender(q)
	}
	fmt.Println(exampid, "started senders")
	//
	wgsend.Wait()
	fmt.Println(exampid, "senders complete")
	wgrecv.Wait()
	fmt.Println(exampid, "receivers complete")
	//
	fmt.Println(exampid, "ends")
}
