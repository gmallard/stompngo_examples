//
// Copyright © 2015-2016 Guy M. Allard
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

/*
Show receiving a RECIPT, requested from an ACK.

	Examples:

		# Using a broker with all defaults:
		# Host is "localhost"
		# Port is 61613
		# Login is "guest"
		# Passcode is "guest
		# Virtual Host is "localhost"
		# Protocol is 1.1
		go run onack.go

		# Using a broker using a custom host and port:
		STOMP_HOST=tjjackson STOMP_PORT=62613 go run onack.go

		# Using a broker using a custom port and virtual host:
		STOMP_PORT=41613 STOMP_VHOST="/" go run onack.go

		# Using a broker using a custom login and passcode:
		STOMP_LOGIN="userid" STOMP_PASSCODE="t0ps3cr3t" go run onack.go
*/
package main

import (
	"log"
	"net"
	"os"
	//
	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "onack: "
	ll      = log.New(os.Stdout, "OACK ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

func main() {

	// Make sure that the queue used by this example do not exist, or are
	// empty.

	// Following is a lengthy piece of code.  Read it striaght from top
	// to bottom.  There is zero comlex logic here.

	// Here is what we will do:
	// Phase 1:
	// - Connect to a broker
	// - Verify a connection spec level
	// - Send a single message to the specified queue on that broker
	// - Disconnect from that broker
	//
	// Phase 2:
	// - Reconnect to the same broker
	// - Subscribe to the specified queue, using "ack:client-individual"
	// - Receive a single message
	// - Send an ACK, asking for a receipt
	// - Receive a RECEIPT # The point of this exercise.
	// - Show data from the RECEIPT and verify it
	// - Disconnect from the broker

	ll.Printf("%s starts\n", exampid)

	// **************************************** Phase 1
	// Set up the connection.
	h, p := senv.HostAndPort()
	hap := net.JoinHostPort(h, p)
	n, e := net.Dial("tcp", hap)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s dial1_complete hap:%s\n",
		exampid, hap)
	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}

	if conn.Protocol() == stompngo.SPL_10 {
		ll.Fatalf("%s v1:%v\n", exampid, "STOMP 1.0 not supported for this example")
	}
	ll.Printf("%s connsess:%s stomp_connect1_complete protocol:%s\n",
		exampid, conn.Session(),
		conn.Protocol())

	// ****************************************
	// App logic here .....

	d := senv.Dest()
	// Prep
	ll.Printf("%s connsess:%s d%s\n",
		exampid, conn.Session(),
		d)

	// ****************************************
	// Send exactly one message.
	sh := stompngo.Headers{"destination", senv.Dest()}
	if senv.Persistent() {
		sh = sh.Add("persistent", "true")
	}
	m := exampid + " message: "
	t := m + "1"

	ll.Printf("%s connsess:%s sending_now t:%s\n",
		exampid, conn.Session(),
		t)
	e = conn.Send(sh, t)
	if e != nil {
		ll.Fatalf("%s v1:%v v2:%v\n", exampid, "bad send", e) // Handle this ...
	}

	ll.Printf("%s connsess:%s send_complete t:%s\n",
		exampid, conn.Session(),
		t)

	// ****************************************
	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s connsess:%s stomp_disconnect1_complete t:%s\n",
		exampid, conn.Session(),
		t)
	// Close the network connection
	e = n.Close()
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s connsess:%s net_close1_complete t:%s\n",
		exampid, conn.Session(),
		t)

	// **************************************** Phase 2

	n, e = net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}

	ll.Printf("%s dial2_complete hap:%s\n",
		exampid, net.JoinHostPort(h, p))

	conn, e = stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalf("%s  f4v:%v\n", exampid, 10, e) // Handle this ......
	}
	ll.Printf("%s connsess:%s stomp_connect2_complete protocol:%s\n",
		exampid, conn.Session(),
		conn.Protocol())

	// ****************************************
	// Subscribe here
	id := stompngo.Uuid()

	var md stompngo.MessageData // A message data instance

	// Get the "subscribe channel"
	sc := sngecomm.HandleSubscribe(conn, d, id, "client-individual")
	ll.Printf("%s connsess:%s stomp_subscribe_complete\n",
		exampid, conn.Session())

	// Get data from the broker
	select {
	case md = <-sc:
	case md = <-conn.MessageData:
		// This would be contain an ERROR or RECEIPT frame.  Both are unexpected
		// in this example.
		ll.Fatalf("%s v1:%v\n", exampid, md) // Handle this
	}
	ll.Printf("%s connsess:%s channel_read_complete\n",
		exampid, conn.Session())

	// MessageData has two components:
	// a) a Message struct
	// b) an Error value.  Check the error value as usual
	if md.Error != nil {
		ll.Fatalf("%s v1:%v\n", exampid, md.Error) // Handle this
	}

	ll.Printf("%s connsess:%s read_message_COMMAND command:%s\n",
		exampid, conn.Session(),
		md.Message.Command)
	ll.Printf("%s connsess:%s read_message_HEADERS headers:%v\n",
		exampid, conn.Session(),
		md.Message.Headers)
	ll.Printf("%s connsess:%s read_message_BODY body:%s\n",
		exampid, conn.Session(),
		string(md.Message.Body))

	// Here we need to send an ACK.  Required Headers are different between
	// a 1.1 and a 1.2 connection level.
	var ah stompngo.Headers
	if conn.Protocol() == stompngo.SPL_11 { // 1.1
		ah = ah.Add("subscription", md.Message.Headers.Value("subscription"))
		ah = ah.Add("message-id", md.Message.Headers.Value("message-id"))
	} else { // 1.2
		ah = ah.Add("id", md.Message.Headers.Value("ack"))
	}

	// We are also going to ask for a RECEIPT for the ACK
	rid := "1"
	ah = ah.Add("receipt", rid)
	//
	e = conn.Ack(ah)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}

	// ****************************************
	// Finally get the RECEIPT.  Where is it?  It is *not* on the "subscribe
	// channel".  It is on the connection level MessageData channel.  Why?
	// Because the broker does *not* include a "subscription" header in
	// RECEIPT frames..
	// ****************************************

	// ***IMPORTANT***
	// ***NOTE*** which channel this RECEIPT MessageData comes in on.
	var rd stompngo.MessageData

	ll.Printf("%s connsess:%s start_receipt_read\n",
		exampid, conn.Session())
	select {
	case rd = <-sc:
		// This would contain a MESSAGE frame.  It is unexpected here
		// in this example.
		ll.Fatalf("%s v1:%v\n", exampid, md) // Handle this
	case rd = <-conn.MessageData: // RECEIPT frame s/b in the MessageData
		if rd.Message.Command != stompngo.RECEIPT {
			ll.Fatalf("%s v1:%v\n", exampid, md) // Handle this
		}
	}
	ll.Printf("%s connsess:%s end_receipt_read\n",
		exampid, conn.Session())
	// ****************************************
	// Show stuff about the RECEIPT MessageData struct

	ll.Printf("%s connsess:%s receipt_COMMAND command:%s\n",
		exampid, conn.Session(),
		rd.Message.Command)
	ll.Printf("%s connsess:%s receipt_HEADERS headers:%v\n",
		exampid, conn.Session(),
		rd.Message.Headers)
	ll.Printf("%s connsess:%s receipt_BODY body:%s\n",
		exampid, conn.Session(),
		string(rd.Message.Body))

	// ****************************************
	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}

	ll.Printf("%s connsess:%s stomp_disconnect2_complete\n",
		exampid, conn.Session())
	// Close the network connection
	e = n.Close()
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}

	ll.Printf("%s connsess:%s net_close2_complete\n",
		exampid, conn.Session())
}
