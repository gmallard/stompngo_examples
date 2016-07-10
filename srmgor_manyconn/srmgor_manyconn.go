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
	"github.com/davecheney/profile"
	//
	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "srmgor_manyconn:"

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
)

func sendMessages(conn *stompngo.Connection, qnum int, nc net.Conn) {
	qns := fmt.Sprintf("%d", qnum) // queue number
	d := senv.Dest() + "." + qns
	ll.Printf("%s connsess:%s sendMessages_start d:%s qnum:%d\n",
		exampid, conn.Session(), d, qnum)
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

		ll.Printf("%s connsess:%s sendMessages_message mc:%d qnum:%d\n",
			exampid, conn.Session(), mc, qnum)
		e := conn.Send(sh, string(sngecomm.Partial()))
		if e != nil {
			ll.Fatalln(exampid, "send:", e, nc.LocalAddr().String(), qnum)
		}
		if mc == nmsgs {
			break
		}
		if sw {
			runtime.Gosched() // yield for this example
			dt := time.Duration(sngecomm.ValueBetween(min, max, sf))
			ll.Printf("%s connsess:%s sendMessages_stagger dt:%v qnum:%d\n",
				exampid, conn.Session(), dt, qnum)
			tmr.Reset(dt)
			_ = <-tmr.C
		}
	}
}

func receiveMessages(conn *stompngo.Connection, qnum int, nc net.Conn) {
	qns := fmt.Sprintf("%d", qnum) // queue number
	d := senv.Dest() + "." + qns

	ll.Printf("%s connsess:%s receiveMessages_start d:%s qnum:%d nmsgs:%d\n",
		exampid, conn.Session(), d, qnum, nmsgs)
	// Subscribe
	id := stompngo.Uuid() // A unique subscription ID
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
			ll.Fatalln(exampid, md) // Handle this
		}
		if md.Error != nil {
			ll.Fatalln(exampid, "recv read:", md.Error, nc.LocalAddr().String(), qnum)
		}

		// Sanity check the queue and message numbers
		mcs := fmt.Sprintf("%d", mc) // message number
		if !md.Message.Headers.ContainsKV("qnum", qns) || !md.Message.Headers.ContainsKV("msgnum", mcs) {
			ll.Fatalln("Bad Headers", md.Message.Headers, qns, mcs)
		}

		// Process the inbound message .................
		sl := len(md.Message.Body)
		if pbc > 0 {
			sl = pbc
			if len(md.Message.Body) < sl {
				sl = len(md.Message.Body)
			}
		}
		ll.Printf("%s connsess:%s receiveMessages_msg d:%s body:%s qnum:%d msgnum:%s\n",
			exampid, conn.Session(), d, string(md.Message.Body[0:sl]), qnum,
			md.Message.Headers.Value("msgnum"))
		if mc == nmsgs {
			break
		}
		// Handle ACKs if needed
		if sngecomm.AckMode() != "auto" {
			ah := []string{}
			sngecomm.HandleAck(conn, ah, id)
		}
		//
		if rw {
			runtime.Gosched() // yield for this example
			dt := time.Duration(sngecomm.ValueBetween(min, max, rf))
			ll.Printf("%s connsess:%s recvMessages_stagger dt:%v qnum:%d\n",
				exampid, conn.Session(), dt, qnum)
			tmr.Reset(dt)
			_ = <-tmr.C
		}
	}
	ll.Printf("%s connsess:%s receiveMessages_end d:%s qnum:%d nmsgs:%d\n",
		exampid, conn.Session(), d, qnum, nmsgs)

	// Unsubscribe
	sngecomm.HandleUnsubscribe(conn, d, id)
	//
}

func runReceiver(qnum int) {
	ll.Printf("%s runReceiver_start qnum:%d\n", exampid, qnum)
	// Network Open
	h, p := senv.HostAndPort() // host and port
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalln(exampid, "recv nectonnr:", qnum, e) // Handle this ......
	}

	ll.Printf("%s runReceiver_start open_complete qnum:%d\n", exampid, qnum)
	ll.Printf("%s runReceiver_start open_complete local:%s qnum:%d\n",
		exampid, n.LocalAddr().String(), qnum)
	ll.Printf("%s runReceiver_start open_complete remote:%s qnum:%d\n",
		exampid, n.RemoteAddr().String(), qnum)

	// Stomp connect
	ch := sngecomm.ConnectHeaders()

	ll.Printf("%s runReceiver_start ch:%v  vhost:%s protocol:%s qnum:%d\n",
		exampid, ch, senv.Vhost(), senv.Protocol(), qnum)
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalln(exampid, "recv stompconnect:", qnum, e) // Handle this ......
	}

	ll.Printf("%s connsess:%s runReceiver_stompconnect qnum:%d\n",
		exampid, conn.Session(), qnum)

	//
	conn.SetSubChanCap(senv.SubChanCap()) // Experiment with this value, YMMV
	// Receives
	receiveMessages(conn, qnum, n)

	ll.Printf("%s connsess:%s runReceiver_receives_complete qnum:%d\n",
		exampid, conn.Session(), qnum)

	// Disconnect from Stomp server
	eh := stompngo.Headers{"recv_discqueue", fmt.Sprintf("%d", qnum)}
	e = conn.Disconnect(eh)
	if e != nil {
		ll.Fatalln(exampid, "recv disconnects:", qnum, e) // Handle this ......
	}

	ll.Printf("%s connsess:%s runReceiver_disconnect_complete qnum:%d\n",
		exampid, conn.Session(), qnum)

	// Network close
	e = n.Close()
	if e != nil {
		ll.Fatalln(exampid, "recv netcloser", qnum, e) // Handle this ......
	}
	// ll.Println(exampid, "recv network close complete", qnum)
	ll.Printf("%s connsess:%s runReceiver_net_close_complete qnum:%d\n",
		exampid, conn.Session(), qnum)

	sngecomm.ShowStats(exampid, "recv "+fmt.Sprintf("%d", qnum), conn)
	wgr.Done()
}

func runSender(qnum int) {

	ll.Printf("%s runSender_start qnum:%d\n", exampid, qnum)
	// Network Open
	h, p := senv.HostAndPort() // host and port
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalln(exampid, "send nectonnr:", qnum, e) // Handle this ......
	}

	ll.Printf("%s runSender_start open_complete qnum:%d\n", exampid, qnum)
	ll.Printf("%s runSender_start open_complete local:%s qnum:%d\n",
		exampid, n.LocalAddr().String(), qnum)
	ll.Printf("%s runSender_start open_complete remote:%s qnum:%d\n",
		exampid, n.RemoteAddr().String(), qnum)

	// Stomp connect
	ch := sngecomm.ConnectHeaders()

	ll.Printf("%s runSender_start ch:%v  vhost:%s protocol:%s qnum:%d\n",
		exampid, ch, senv.Vhost(), senv.Protocol(), qnum)
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalln(exampid, "send stompconnect:", qnum, e) // Handle this ......
	}

	ll.Printf("%s connsess:%s runSender_stompconnect qnum:%d\n",
		exampid, conn.Session(), qnum)
	//
	sendMessages(conn, qnum, n)

	ll.Printf("%s connsess:%s runSender_sends_complete qnum:%d\n",
		exampid, conn.Session(), qnum)
	// Disconnect from Stomp server
	eh := stompngo.Headers{"send_discqueue", fmt.Sprintf("%d", qnum)}
	e = conn.Disconnect(eh)
	if e != nil {
		ll.Fatalln(exampid, "send disconnects:", qnum, e) // Handle this ......
	}

	ll.Printf("%s connsess:%s runSender_disconnect_complete qnum:%d\n",
		exampid, conn.Session(), qnum)
	// Network close
	e = n.Close()
	if e != nil {
		ll.Fatalln(exampid, "send netcloser", qnum, e) // Handle this ......
	}

	ll.Printf("%s connsess:%s runSender_net_close_complete qnum:%d\n",
		exampid, conn.Session(), qnum)
	sngecomm.ShowStats(exampid, "send "+fmt.Sprintf("%d", qnum), conn)
	wgs.Done()
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
	numq := sngecomm.Nqs()
	nmsgs = senv.Nmsgs() // message count
	//
	ll.Println(exampid, "main starting receivers")
	for q := 1; q <= numq; q++ {
		wgr.Add(1)
		go runReceiver(q)
	}
	ll.Println(exampid, "main started receivers")
	//
	ll.Println(exampid, "main starting senders")
	for q := 1; q <= numq; q++ {
		wgs.Add(1)
		go runSender(q)
	}
	ll.Println(exampid, "main started senders")
	//
	wgs.Wait()
	ll.Println(exampid, "main senders complete")
	wgr.Wait()
	ll.Println(exampid, "main receivers complete")
	//
	ll.Println(exampid, "main ends", time.Since(tn))
}
