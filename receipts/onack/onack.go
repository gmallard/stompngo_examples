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
	"os"
	"time"
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

	tag = "onackmain"
)

func main() {

	// Make sure that the queue used by this example do not exist, or are
	// empty.

	// Following is a lengthy piece of code.  Read it striaght from top
	// to bottom.  There is zero complex logic here.

	// What this code will do:
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
	//**************************************************************************
	// - Receive a RECEIPT // The point of this exercise.
	// - Show data from the RECEIPT and verify it // The point of this exercise.
	//**************************************************************************
	// - Disconnect from the broker

	// Start

	st := time.Now()

	// **************************************** Phase 1
	// Set up the connection.
	// Standard example connect sequence
	n, conn, e := sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	// ****************************************
	// App logic here .....

	d := senv.Dest()
	ll.Printf("%stag:%s connsess:%s destination:%v\n",
		exampid, tag, conn.Session(),
		d)

	// ****************************************
	// Send exactly one message.
	sh := stompngo.Headers{"destination", senv.Dest()}
	if senv.Persistent() {
		sh = sh.Add("persistent", "true")
	}
	m := exampid + " message: "
	t := m + "1"
	ll.Printf("%stag:%s connsess:%s sending_now body:%v\n",
		exampid, tag, conn.Session(),
		t)
	e = conn.Send(sh, t)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_bad_send error:%v",
			exampid, tag, conn.Session(),
			e.Error()) // Handle this ......
	}
	ll.Printf("%stag:%s connsess:%s send_complete body:%v\n",
		exampid, tag, conn.Session(),
		t)

	// ****************************************
	// Disconnect from the Stomp server
	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_disconnect error:%v",
			exampid, tag, conn.Session(),
			e.Error()) // Handle this ......
	}

	// **************************************** Phase 2

	// Standard example connect sequence
	n, conn, e = sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	// ****************************************
	// Subscribe here
	id := stompngo.Uuid()
	// Get the "subscribe channel"
	sc := sngecomm.HandleSubscribe(conn, d, id, "client-individual")
	ll.Printf("%stag:%s connsess:%s stomp_subscribe_complete\n",
		exampid, tag, conn.Session())

	// Get data from the broker
	var md stompngo.MessageData // A message data instance
	select {
	case md = <-sc:
	case md = <-conn.MessageData:
		// This would be contain an ERROR or RECEIPT frame.  Both are unexpected
		// in this example.
		ll.Fatalf("%stag:%s connsess:%s bad_frame md:%v",
			exampid, tag, conn.Session(),
			md) // Handle this ......
	}
	ll.Printf("%stag:%s connsess:%s channel_read_complete\n",
		exampid, tag, conn.Session())

	// MessageData has two components:
	// a) a Message struct
	// b) an Error value.  Check the error value as usual
	if md.Error != nil {
		ll.Fatalf("%stag:%s connsess:%s message_error md:%v",
			exampid, tag, conn.Session(),
			md.Error) // Handle this ......
	}

	ll.Printf("%stag:%s connsess:%s read_message_COMMAND command:%s\n",
		exampid, tag, conn.Session(),
		md.Message.Command)
	ll.Printf("%stag:%s connsess:%s read_message_HEADERS headers:%s\n",
		exampid, tag, conn.Session(),
		md.Message.Headers)
	ll.Printf("%stag:%s connsess:%s read_message_BODY body:%s\n",
		exampid, tag, conn.Session(),
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
	rid := "receipt-001"
	ah = ah.Add("receipt", rid)
	//
	ll.Printf("%stag:%s connsess:%s ACK_receipt_headers headers:%v\n",
		exampid, tag, conn.Session(),
		ah)
	e = conn.Ack(ah)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s ack_error error:%v",
			exampid, tag, conn.Session(),
			e.Error()) // Handle this ......
	}

	// ****************************************
	// Finally get the RECEIPT.  Where is it?  It is *not* on the "subscribe
	// channel".  It is on the connection level MessageData channel.  Why?
	// Because the broker does *not* include a "subscription" header in
	// RECEIPT frames..
	// ****************************************

	// ***IMPORTANT***
	// ***NOTE*** which channel this RECEIPT MessageData comes in on.

	ll.Printf("%stag:%s connsess:%s start_receipt_read\n",
		exampid, tag, conn.Session())
	var rd stompngo.MessageData
	select {
	case rd = <-sc:
		// This would contain a MESSAGE frame.  It is unexpected here
		// in this example.
		ll.Fatalf("%stag:%s connsess:%s bad_frame_channel rd:%v\n",
			exampid, tag, conn.Session(),
			rd) // Handle this ......
	case rd = <-conn.MessageData: // RECEIPT frame s/b in the MessageData
		// Step 1 of Verify
		if rd.Message.Command != stompngo.RECEIPT {
			ll.Fatalf("%stag:%s connsess:%s bad_frame_command rd:%v\n",
				exampid, tag, conn.Session(),
				rd) // Handle this ......
		}
	}
	ll.Printf("%stag:%s connsess:%s end_receipt_read\n",
		exampid, tag, conn.Session())

	// ****************************************
	// Show details about the RECEIPT MessageData struct
	ll.Printf("%stag:%s connsess:%s receipt_COMMAND command:%s\n",
		exampid, tag, conn.Session(),
		rd.Message.Command)
	ll.Printf("%stag:%s connsess:%s receipt_HEADERS headers:%v\n",
		exampid, tag, conn.Session(),
		rd.Message.Headers)
	ll.Printf("%stag:%s connsess:%s receipt_BODY body:%s\n",
		exampid, tag, conn.Session(),
		string(rd.Message.Body))

	// Step 2 of Verify
	// Verify that the receipt has the id we asked for
	if rd.Message.Headers.Value("receipt-id") != rid {
		ll.Fatalf("%stag:%s connsess:%s bad_receipt_id wanted:%v got:%v\n",
			exampid, tag, conn.Session(),
			rid, rd.Message.Headers.Value("receipt-id")) // Handle this ......
	}
	ll.Printf("%stag:%s connsess:%s receipt_id_verified rid:%s\n",
		exampid, tag, conn.Session(),
		rid)

	// ****************************************
	// Disconnect from the Stomp server
	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_disconnect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, conn.Session(),
		time.Now().Sub(st))

}
