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
Show receiving a RECIPT, requested from a SEND.

	Examples:

		# Publish to a broker with all defaults:
		# Host is "localhost"
		# Port is 61613
		# Login is "guest"
		# Passcode is "guest
		# Virtual Host is "localhost"
		# Protocol is 1.1
		go run onsend.go

		# Publish to a broker using a custom host and port:
		STOMP_HOST=tjjackson STOMP_PORT=62613 go run onsend.go

		# Publish to a broker using a custom port and virtual host:
		STOMP_PORT=41613 STOMP_VHOST="/" go run onsend.go

		# Publish to a broker using a custom login and passcode:
		STOMP_LOGIN="userid" STOMP_PASSCODE="t0ps3cr3t" go run onsend.go

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
	exampid = "onsend: "
	ll      = log.New(os.Stdout, "OSND ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	tag = "onsendmain"
)

func main() {

	st := time.Now()

	// ****************************************
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
	// Send exactly one message.  Ask for a receipt.

	d := senv.Dest()
	ll.Printf("%stag:%s connsess:%s destination:%v\n",
		exampid, tag, conn.Session(),
		d)

	rid := "recipt-002" // The receipt ID
	sh := stompngo.Headers{"destination", d,
		"receipt", rid} // send headers
	if senv.Persistent() {
		sh = sh.Add("persistent", "true")
	}
	ms := exampid + " message: "
	t := ms + rid
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

	// Look for the receipt
	ll.Printf("%stag:%s connsess:%s start_receipt_read\n",
		exampid, tag, conn.Session())
	rd := <-conn.MessageData // The RECEIPT frame should be on this channel
	if rd.Message.Command != stompngo.RECEIPT {
		ll.Fatalf("%stag:%s connsess:%s bad_frame_command rd:%v\n",
			exampid, tag, conn.Session(),
			rd) // Handle this ......
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
