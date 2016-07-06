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
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var exampid = "onack: "

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

	log.Println(exampid + "starts ...")

	// **************************************** Phase 1
	// Set up the connection.
	h, p := sngecomm.HostAndPort()
	log.Println(exampid+"host", h, "port", p)
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "dial 1 complete ...")
	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}

	if conn.Protocol() == stompngo.SPL_10 {
		panic("STOMP 1.0 not supported for this example")
	}
	log.Println(exampid+"stomp connect 1 complete ...", conn.Protocol())

	// ****************************************
	// App logic here .....

	// Prep
	log.Println(sngecomm.ExampIdNow(exampid), "dest:", sngecomm.Dest())

	// ****************************************
	// Send exactly one message.
	s := stompngo.Headers{"destination", sngecomm.Dest()}
	m := exampid + " message: "
	t := m + "1"
	log.Println(sngecomm.ExampIdNow(exampid), "sending now:", t)
	e = conn.Send(s, t)
	if e != nil {
		log.Fatalln("bad send", e) // Handle this ...
	}
	log.Println(sngecomm.ExampIdNow(exampid), "send complete:", t)

	// ****************************************
	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "stomp disconnect 1 complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "network close 1 complete ...")

	// **************************************** Phase 2

	n, e = net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "dial 2 complete ...")

	conn, e = stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(10, e) // Handle this ......
	}
	log.Println(exampid+"stomp connect 2 complete ...", conn.Protocol())

	// ****************************************
	// Subscribe here
	d := sngecomm.Dest()
	i := stompngo.Uuid()

	// Get the "subscribe channel"
	sc := sngecomm.Subscribe(conn, d, i, "client-individual")
	log.Println(exampid + "stomp subscribe complete ...")
	// Get what is on the subscribe channel
	md := <-sc
	log.Println(exampid + "channel read complete ...")
	// MessageData has two components:
	// a) a Message struct
	// b) an Error value.  Check the error value as usual
	if md.Error != nil {
		log.Fatalln(md.Error) // Handle this
	}
	log.Println(exampid+"read message COMMAND", md.Message.Command)
	log.Println(exampid+"read message HEADERS", md.Message.Headers)
	log.Println(exampid+"read message BODY", string(md.Message.Body))

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
		log.Fatalln(e) // Handle this ......
	}

	// Finally get the RECEIPT.  Where is it?  It is *not* on the "subscribe
	// channel".  It is on the connection level MessageData channel.  Why?
	// Because it does *not* have a "subscription" header.
	// ****************************************
	// ***IMPORTANT***
	// ***NOTE*** which channel this RECEIPT MessageData comes in on.
	log.Println(sngecomm.ExampIdNow(exampid), "start receipt read")
	r := <-conn.MessageData
	log.Println(sngecomm.ExampIdNow(exampid), "end receipt read")

	// ****************************************
	// Show stuff about the RECEIPT MessageData struct
	log.Println(sngecomm.ExampIdNow(exampid), "receipt COMMAND", r.Message.Command)
	log.Println(sngecomm.ExampIdNow(exampid), "receipt HEADERS", r.Message.Headers)
	log.Println(sngecomm.ExampIdNow(exampid), "receipt BODY", string(r.Message.Body))

	// ****************************************
	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "stomp disconnect 2 complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	log.Println(exampid + "network close 2 complete ...")

}
