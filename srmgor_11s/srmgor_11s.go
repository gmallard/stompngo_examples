//
// Copyright © 2011-2012 Guy M. Allard
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
	"crypto/rand"
	"fmt"
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
	"log"
	"math/big"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
)

var exampid = "srmgor_11: "

var wgsend sync.WaitGroup
var wgrecv sync.WaitGroup
var wgall sync.WaitGroup

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

//
var n net.Conn                // Network Connection
var conn *stompngo.Connection // Stomp Connection

// Get a duration between min amd max
func timeBetween(min, max int64) int64 {
	br, _ := rand.Int(rand.Reader, big.NewInt(max-min)) // Ignore errors here
	return (br.Add(big.NewInt(min), br).Int64()) / 2
}

// Send messages to a particular queue
func sender(qn, c int) {
	fmt.Println(exampid + "sender starts ...")
	//
	qp := sngecomm.Dest()        // queue name prefix
	qns := fmt.Sprintf("%d", qn) // queue number
	q := qp + "." + qns
	fmt.Println(exampid + "sender queue name: " + q)
	h := stompngo.Headers{"destination", q} // send Headers

	// Send loop
	for i := 1; i <= c; i++ {
		si := fmt.Sprintf("%d", i)
		// Generate a message to send ...............
		m := exampid + "|" + "payload" + "|qnum:" + qns + "|msgnum:" + si
		fmt.Println("sender", m)
		e := conn.Send(h, m)
		if e != nil {
			log.Fatalln(e)
		}
		if send_wait {
			runtime.Gosched()                                              // yield for this example
			time.Sleep(time.Duration(send_factor * timeBetween(min, max))) // Time to build next message
		}
	}
	// Sending is done 
	fmt.Println(exampid + "sender ends ...")
	wgsend.Done()
}

// Receive messages from a particular queue
func receiver(qn, c int) {
	fmt.Println(exampid + "receiver starts ...")
	//
	qp := sngecomm.Dest()        // queue name prefix
	qns := fmt.Sprintf("%d", qn) // queue number
	q := qp + "." + qns
	fmt.Println(exampid + "receiver queue name: " + q)
	u := stompngo.Uuid() // A unique subscription ID
	h := stompngo.Headers{"destination", q, "id", u}
	// Subscribe
	r, e := conn.Subscribe(h)
	if e != nil {
		log.Fatalln(e)
	}
	// Receive loop
	for i := 1; i <= c; i++ {
		d := <-r
		if d.Error != nil {
			log.Fatalln(d.Error)
		}

		// Process the inbound message .................
		m := d.Message.BodyString()
		fmt.Println("receiver", m)

		// Sanity check the queue and message numbers
		mns := fmt.Sprintf("%d", i) // message number
		t := "|qnum:" + qns + "|msgnum:" + mns
		if !strings.Contains(m, t) {
			log.Fatalln("bad message", m, t)
		}
		if recv_wait {
			runtime.Gosched()                                              // yield for this example
			time.Sleep(time.Duration(recv_factor * timeBetween(min, max))) // Time to process this message
		}
	}
	// Unsubscribe
	e = conn.Unsubscribe(h)
	if e != nil {
		log.Fatalln(e)
	}
	// Receiving is done 
	fmt.Println(exampid + "receiver ends ...")
	wgrecv.Done()
}

func startSenders(q int) {
	fmt.Println(exampid + "startSenders starts ...")

	c := sngecomm.Nmsgs() // message count
	fmt.Printf(exampid+"startSenders message count: %d\n", c)
	for i := 1; i <= q; i++ { // all queues
		wgsend.Add(1)
		go sender(i, c)
	}
	wgsend.Wait()

	fmt.Println(exampid + "startSenders ends ...")
	wgall.Done()
}

func startReceivers(q int) {
	fmt.Println(exampid + "startReceivers starts ...")

	c := sngecomm.Nmsgs() // get message count
	fmt.Printf(exampid+"startReceivers message count: %d\n", c)
	for i := 1; i <= q; i++ { // all queues
		wgrecv.Add(1)
		go receiver(i, c)
	}
	wgrecv.Wait()

	fmt.Println(exampid + "startReceivers ends ...")
	wgall.Done()
}

// Show a number of writers and readers operating concurrently from unique
// destinations.
func main() {
	fmt.Println(exampid + "starts ...")
	// Wait flags
	send_wait = sngecomm.SendWait()
	recv_wait = sngecomm.RecvWait()
	// Number of queues
	q := sngecomm.Nqs()
	fmt.Printf(exampid+"Nqs: %d\n", q)
	// Open net and stomp connections
	h, p := sngecomm.HostAndPort11() // a 1.1 connect
	var e error
	// Network open
	n, e = net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	// Stomp connect
	eh := stompngo.Headers{}
	conn, e = stompngo.Connect(n, eh)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	// Run everything
	wgall.Add(2)
	go startReceivers(q)
	go startSenders(q)
	wgall.Wait()
	// Disconnect from Stomp server
	e = conn.Disconnect(eh)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	// Network close
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(exampid + "ends ...")
}
