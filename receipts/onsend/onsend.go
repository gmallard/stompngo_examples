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
	ll.Println(exampid + "starts ...")

	// ****************************************
	// Set up the connection.
	h, p := senv.HostAndPort()
	ll.Println(exampid+"host", h, "port", p)
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalln(exampid, e) // Handle this ......
	}
	ll.Println(exampid + "dial complete ...")
	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalln(exampid, e) // Handle this ......
	}
	ll.Println(exampid+"stomp connect complete ...", conn.Protocol())

	// ****************************************
	// App logic here .....

	// Prep
	ll.Println(exampid, "dest:", senv.Dest())

	// ****************************************
	// Send exactly one message.  Ask for a receipt.
	rid := "1" // The receipt ID
	sh := stompngo.Headers{"destination", senv.Dest(),
		"receipt", rid} // send headers
	ms := exampid + " message: "
	t := ms + rid
	ll.Println(exampid, "sending now:", t)
	e = conn.Send(sh, t)
	if e != nil {
		ll.Fatalln(exampid, "bad send", e) // Handle this ...
	}
	ll.Println(exampid, "send complete:", t)

	// ****************************************
	// OK, here we are looking for a RECEIPT.
	// The MessageData struct of the receipt will be on the connection
	// level MessageData channel.
	ll.Println(exampid, "start receipt read")

	// ****************************************
	// ***IMPORTANT***
	// ***NOTE*** which channel this RECEIPT MessageData comes in on.
	rd := <-conn.MessageData
	ll.Println(exampid, "end receipt read")

	// ****************************************
	// Show stuff about the RECEIPT MessageData struct
	ll.Println(exampid, "COMMAND", rd.Message.Command)
	ll.Println(exampid, "HEADERS", rd.Message.Headers)
	ll.Println(exampid, "BODY", string(rd.Message.Body))

	// ****************************************
	// Get the returned ID.
	irid := rd.Message.Headers.Value("receipt-id")
	ll.Println(exampid, "irid", irid)
	// Check that it matches what we asked for
	if rid != irid {
		ll.Fatalln(exampid, "notsame", rid, irid) // Handle this ......
	}
	ll.Println(exampid, "validation complete")

	// ****************************************
	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalln(exampid, e) // Handle this ......
	}
	ll.Println(exampid + "stomp disconnect complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		ll.Fatalln(exampid, e) // Handle this ......
	}
	ll.Println(exampid + "network close complete ...")

	ll.Println(exampid + "ends ...")

}
