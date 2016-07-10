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
	exampid = "onsend: "
	ll      = log.New(os.Stdout, "OSND ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

func main() {

	ll.Printf("%s starts\n", exampid)

	// ****************************************
	// Set up the connection.
	h, p := senv.HostAndPort()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s dial_complete hap:%s\n",
		exampid, net.JoinHostPort(h, p))
	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s connsess:%s stomp_connect_complete protocol:%s\n",
		exampid, conn.Session(),
		conn.Protocol())

	// ****************************************
	// App logic here .....

	d := senv.Dest()
	ll.Printf("%s connsess:%s dest d%s\n",
		exampid, conn.Session(),
		d)

	// Prep

	// ****************************************
	// Send exactly one message.  Ask for a receipt.
	rid := "1" // The receipt ID
	sh := stompngo.Headers{"destination", d,
		"receipt", rid} // send headers
	if senv.Persistent() {
		sh = sh.Add("persistent", "true")
	}
	ms := exampid + " message: "
	t := ms + rid
	ll.Printf("%s connsess:%s sending_now t:%s\n",
		exampid, conn.Session(),
		t)
	e = conn.Send(sh, t)
	if e != nil {
		ll.Fatalln(exampid, "bad send", e) // Handle this ...
	}
	ll.Printf("%s connsess:%s send_complete t:%s\n",
		exampid, conn.Session(),
		t)

	// ****************************************
	// OK, here we are looking for a RECEIPT.
	// The MessageData struct of the receipt will be on the connection
	// level MessageData channel.
	ll.Printf("%s connsess:%s start_receipt_read t:%s\n",
		exampid, conn.Session(),
		t)

	// ****************************************
	// ***IMPORTANT***
	// ***NOTE*** which channel this RECEIPT MessageData comes in on.
	// ***NOTE*** we do not use a select here with a subscribe channel.
	//            because we have not SUBSCRIB'ed.
	rd := <-conn.MessageData

	if rd.Message.Command != stompngo.RECEIPT {
		ll.Fatalln(exampid, rd) // Handle this
	}

	ll.Printf("%s connsess:%s end_receipt_read t:%s\n",
		exampid, conn.Session(),
		t)

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
	// Get the returned ID.
	irid := rd.Message.Headers.Value("receipt-id")

	ll.Printf("%s connsess:%s received_receipt_id irif:%s\n",
		exampid, conn.Session(),
		irid)
	// Check that it matches what we asked for
	if rid != irid {
		ll.Fatalln(exampid, "notsame", rid, irid) // Handle this ......
	}

	ll.Printf("%s connsess:%s validation_complete\n",
		exampid, conn.Session())
	// ****************************************
	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}

	ll.Printf("%s connsess:%s stomp_disconnect_complete\n",
		exampid, conn.Session())
	// Close the network connection
	e = n.Close()
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}

	ll.Printf("%s connsess:%s net_close2_complete\n",
		exampid, conn.Session())

	ll.Printf("%s ends\n", exampid)

}
