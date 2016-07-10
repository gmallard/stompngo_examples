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
	exampid = "srmgor_2conn:"

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
)

// Send messages to a particular queue
func sender(conn *stompngo.Connection, qn, nmsgs int) {
	qns := fmt.Sprintf("%d", qn) // queue number
	ll.Printf("%s connsess:%s sender_starts qn:%d nmsgs:%d\n",
		exampid, conn.Session(), qn, nmsgs)
	//
	// qp := senv.Dest() // queue name prefix
	d := senv.Dest() + "." + qns
	ll.Printf("%s connsess:%s sender_starts d:%s\n",
		exampid, conn.Session(), d)
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
		ll.Printf("%s connsess:%s sender_message qns:%s si:%s\n",
			exampid, conn.Session(), qns, si)
		e := conn.Send(sh, string(sngecomm.Partial()))

		if e != nil {
			ll.Fatalln(exampid, "send error", e, qn)
			break
		}
		if sw {
			runtime.Gosched() // yield for this example
			dt := time.Duration(sngecomm.ValueBetween(min, max, sf))
			ll.Printf("%s connsess:%s send_stagger dt:%v qns:%s\n",
				exampid, conn.Session(),
				dt, qns)
			tmr.Reset(dt)
			_ = <-tmr.C
		}
	}
	// Sending is done
	ll.Printf("%s connsess:%s sender_ends qn:%d nmsgs:%d\n",
		exampid, conn.Session(), qn, nmsgs)
	wgs.Done()
}

// Asynchronously process all messages for a given subscription.
func receiveWorker(sc <-chan stompngo.MessageData, qns string, nmsgs int,
	qc chan<- bool, conn *stompngo.Connection, id string) {
	//
	tmr := time.NewTimer(100 * time.Hour)

	pbc := sngecomm.Pbc() // Print byte count

	// Receive loop
	var md stompngo.MessageData
	for i := 1; i <= nmsgs; i++ {

		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// Frames RECEIPT or ERROR not expected here
			ll.Fatalln(exampid, md) // Handle this
		}
		if md.Error != nil {
			ll.Fatalln(exampid, "recv read error", md.Error, qns)
		}

		// Sanity check the queue and message numbers
		mns := fmt.Sprintf("%d", i) // message number
		if !md.Message.Headers.ContainsKV("qnum", qns) || !md.Message.Headers.ContainsKV("msgnum", mns) {
			ll.Fatalln("Bad Headers", md.Message.Headers, qns, mns)
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
		// ll.Println(exampid, "recv message", string(d.Message.Body[0:sl]), qns, d.Message.Headers.Value("msgnum"))
		ll.Printf("%s connsess:%s recv_message body:%s qns:%s msgnum:%s\n",
			exampid, conn.Session(),
			string(md.Message.Body[0:sl]),
			qns,
			md.Message.Headers.Value("msgnum"))
		if i == nmsgs {
			break
		}
		if rw {
			runtime.Gosched() // yield for this example
			dt := time.Duration(sngecomm.ValueBetween(min, max, rf))
			ll.Printf("%s connsess:%s recv_stagger dt:%v qns:%s\n",
				exampid, conn.Session(),
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
	qns := fmt.Sprintf("%d", qn) // queue number
	ll.Printf("%s connsess:%s recveiver_starts qns:%d nmsgs:%d\n",
		exampid, conn.Session(), qn, nmsgs)
	//
	qp := senv.Dest() // queue name prefix
	q := qp + "." + qns
	// ll.Println(exampid, "recv queue name", q, qn)
	ll.Printf("%s connsess:%s recveiver_names q:%s qn:%d\n",
		exampid, conn.Session(), q, qn)
	id := stompngo.Uuid() // A unique subscription ID
	sc := sngecomm.HandleSubscribe(conn, q, id, sngecomm.AckMode())
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
			log.Fatalln("CONN2BUFFER conversion error", e)
		} else {
			bs = int(i)
		}
	}
	if bs < 1 {
		bs = nmsgs
	}
	ll.Printf("%s connsess:%s recveiver_mdbuffersize bs:%d qn:%d\n",
		exampid, conn.Session(), bs, qn)

	// Process all inputs async .......
	// var mc chan stompngo.MessageData
	mdc := make(chan stompngo.MessageData, bs)      // MessageData Buffer size
	dc := make(chan bool)                           // Receive processing done channel
	go receiveWorker(mdc, qns, nmsgs, dc, conn, id) // Start async processor
	for i := 1; i <= nmsgs; i++ {
		mdc <- <-sc // Receive message data as soon as possible, and internally queue it
	}
	ll.Printf("%s connsess:%s recveiver_waitforWorkersBegin qns:%s\n",
		exampid, conn.Session(), qns)
	<-dc // Wait until receive processing is done for this queue
	ll.Printf("%s connsess:%s recveiver_waitforWorkersEnd qns:%s\n",
		exampid, conn.Session(), qns)

	// Unsubscribe
	sngecomm.HandleUnsubscribe(conn, q, id)

	// Receiving is done
	ll.Printf("%s connsess:%s recveiver_ends qns:%s\n",
		exampid, conn.Session(), qns)
	wgr.Done()
}

func startSenders(qn int) {
	ll.Printf("%s startSenders_starts qn:%d\n",
		exampid, qn)

	// Open
	h, p := senv.HostAndPort() // host and port
	hap := net.JoinHostPort(h, p)
	n, e := net.Dial("tcp", hap)
	if e != nil {
		ll.Fatalln(exampid, "startSenders netconnect error", e, qn) // Handle this ......
	}

	// Stomp connect
	ch := sngecomm.ConnectHeaders()
	ll.Printf("%s startSenders_sdata vhost:%s protocol:%s qn:%d\n",
		exampid, senv.Vhost(), senv.Protocol(), qn)
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalln(exampid, "startSenders stompconnect error", e, qn) // Handle this ......
	}
	ll.Printf("%s connsess:%s startSenders_connection qn:%d\n",
		exampid, conn.Session(), qn)
	nmsgs := senv.Nmsgs() // message count
	ll.Printf("%s connsess:%s startSenders_message_count nmsgs:%d qn:%d\n",
		exampid, conn.Session(), nmsgs, qn)
	for i := 1; i <= qn; i++ { // all queues
		wgs.Add(1)
		go sender(conn, i, nmsgs)
	}
	wgs.Wait()

	// Disconnect from Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Println(exampid, "startSenders disconnect error", e, qn) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		ll.Fatalln(exampid, "startSenders netclose error", e, qn) // Handle this ......
	}

	ll.Printf("%s startSenders_ends qn:%d\n",
		exampid, qn)
	sngecomm.ShowStats(exampid, "startSenders", conn)
	wga.Done()
}

func startReceivers(qn int) {
	ll.Printf("%s startReceivers_starts qn:%d\n",
		exampid, qn)

	// Open
	h, p := senv.HostAndPort() // host and port
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalln(exampid, "startReceivers nectonnr:", e, qn) // Handle this ......
	}
	ch := sngecomm.ConnectHeaders()
	ll.Printf("%s startReceivers_sdata vhost:%s protocol:%s qn:%dn",
		exampid, senv.Vhost(), senv.Protocol(), qn)
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalln(exampid, "startReceivers stompconnectr:", e, qn) // Handle this ......
	}
	ll.Printf("%s  connsess:%s startReceivers_conndata qn:%d\n",
		exampid, conn.Session(), qn)
	nmsgs := senv.Nmsgs() // get message count
	ll.Printf("%s  connsess:%s startReceivers_message_count nmsgs:%d qn:%d\n",
		exampid, conn.Session(), nmsgs, qn)
	for i := 1; i <= qn; i++ { // all queues
		wgr.Add(1)
		go receiver(conn, i, nmsgs)
	}
	wgr.Wait()

	// Disconnect from Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Println(exampid, "startReceivers disconnect error", e, qn) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		ll.Println(exampid, "startReceivers netclose error", e, qn) // Handle this ......
	}

	ll.Printf("%s startReceivers_ends qn:%d\n",
		exampid, qn)
	sngecomm.ShowStats(exampid, "startReceivers", conn)
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

	tn := time.Now()
	ll.Println(exampid, "main starts")

	if sngecomm.SetMAXPROCS() {
		nc := runtime.NumCPU()
		ll.Println(exampid, "main number of CPUs is:", nc)
		c := runtime.GOMAXPROCS(nc)
		ll.Println(exampid, "main previous number of GOMAXPROCS is:", c)
		ll.Println(exampid, "main current number of GOMAXPROCS is:", runtime.GOMAXPROCS(-1))
	}
	//
	sw = sngecomm.SendWait()
	rw = sngecomm.RecvWait()
	sf = sngecomm.SendFactor()
	rf = sngecomm.RecvFactor()
	ll.Println(exampid, "main Sleep Factors", "send", sf, "recv", rf)
	//
	q := sngecomm.Nqs()
	//
	wga.Add(2)
	go startReceivers(q)
	go startSenders(q)
	wga.Wait()

	ll.Println(exampid, "main ends", time.Since(tn))
}
