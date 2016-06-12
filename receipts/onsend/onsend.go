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
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var exampid = "onsend: "

func main() {
	log.Println(exampid + "starts ...")

	// ****************************************
	// Set up the connection.
	h, p := sngecomm.HostAndPort()
	log.Println(exampid+"host", h, "port", p)
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "dial complete ...")
	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid+"stomp connect complete ...", conn.Protocol())

	// ****************************************
	// App logic here .....

	// Prep
	log.Println(sngecomm.ExampIdNow(exampid), "dest:", sngecomm.Dest())

	// ****************************************
	// Send exactly one message.  Ask for a receipt.
	s := stompngo.Headers{"destination", sngecomm.Dest(),
		"receipt", "1"} // send headers
	rid := "1" // The receipt ID
	m := exampid + " message: "
	t := m + rid
	log.Println(sngecomm.ExampIdNow(exampid), "sending now:", t)
	e = conn.Send(s, t)
	if e != nil {
		log.Fatalln("bad send", e) // Handle this ...
	}
	log.Println(sngecomm.ExampIdNow(exampid), "send complete:", t)

	// ****************************************
	// OK, here we are looking for a RECEIPT.
	// The MessageData struct of the receipt will be on the connection
	// level MessageData channel.
	log.Println(sngecomm.ExampIdNow(exampid), "start receipt read")

	// ****************************************
	// ***IMPORTANT***
	// ***NOTE*** which channel this RECEIPT MessageData comes in on.
	r := <-conn.MessageData
	log.Println(sngecomm.ExampIdNow(exampid), "end receipt read")

	// ****************************************
	// Show stuff about the RECEIPT MessageData struct
	log.Println(sngecomm.ExampIdNow(exampid), "COMMAND", r.Message.Command)
	log.Println(sngecomm.ExampIdNow(exampid), "HEADERS", r.Message.Headers)
	log.Println(sngecomm.ExampIdNow(exampid), "BODY", string(r.Message.Body))

	// ****************************************
	// Get the returned ID.
	irid := r.Message.Headers.Value("receipt-id")
	log.Println(sngecomm.ExampIdNow(exampid), "irid", irid)
	// Check that it matches what we asked for
	if rid != irid {
		log.Fatalln("notsame", rid, irid) // Handle this ......
	}
	log.Println(sngecomm.ExampIdNow(exampid), "validation complete")

	// ****************************************
	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "stomp disconnect complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "network close complete ...")

	log.Println(exampid + "ends ...")

}
