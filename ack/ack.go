//
// Copyright © 2011-2018 Guy M. Allard
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
Receive messages from a STOMP broker, and ACK them.

	Examples:

		# ACK messages from a broker with all defaults:
		# Host is "localhost"
		# Port is 61613
		# Login is "guest"
		# Passcode is "guest
		# Virtual Host is "localhost"
		# Protocol is 1.2
		go run ack.go

		# ACK messages from a broker using STOMP protocol level 1.0:
		STOMP_PROTOCOL=1.0 go run ack.go

		# ACK messages from a broker using a custom host and port:
		STOMP_HOST=tjjackson STOMP_PORT=62613 go run ack.go

		# ACK messages from a broker using a custom port and virtual host:
		STOMP_PORT=41613 STOMP_VHOST="/" go run ack.go

		# ACK messages from a broker using a custom login and passcode:
		STOMP_LOGIN="userid" STOMP_PASSCODE="t0ps3cr3t" go run ack.go

*/
package main

import (
	"log"
	"os"
	"time"

	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "ack: "
	ll      = log.New(os.Stdout, "EACK ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	tag     = "ackmain"
)

// Connect to a STOMP broker, receive some messages, ACK them, and disconnect.
func main() {

	st := time.Now()

	// Standard example connect sequence
	n, conn, e := sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	pbc := sngecomm.Pbc() // Print byte count

	// *NOTE* your application functionaltiy goes here!
	// With Stomp, you must SUBSCRIBE to a destination in order to receive.
	// Subscribe returns a channel of MessageData struct.
	// Here we use a common utility routine to handle the differing subscribe
	// requirements of each protocol level.
	d := senv.Dest()
	id := stompngo.Uuid()
	sc := sngecomm.HandleSubscribe(conn, d, id, "client")
	ll.Printf("%stag:%s connsess:%s main_subscribe_complete\n",
		exampid, tag, conn.Session())
	// Read data from the returned channel
	var md stompngo.MessageData
	for i := 1; i <= senv.Nmsgs(); i++ {

		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// Frames RECEIPT or ERROR not expected here
			ll.Fatalf("%stag:%s connsess:%s main_channel_frame error:%v",
				exampid, tag, conn.Session(),
				e.Error()) // Handle this ......
		}

		ll.Printf("%stag:%s connsess:%s main_channel_read_complete\n",
			exampid, tag, conn.Session())
		// MessageData has two components:
		// a) a Message struct
		// b) an Error value.  Check the error value as usual
		if md.Error != nil {
			ll.Fatalf("%stag:%s connsess:%s main_channel_read error:%v",
				exampid, tag, conn.Session(),
				e.Error()) // Handle this ......
		}
		//
		ll.Printf("%stag:%s connsess:%s frame_type:%v\n",
			exampid, tag, conn.Session(),
			md.Message.Command)
		if md.Message.Command != stompngo.MESSAGE {
			ll.Fatalf("%stag:%s connsess:%s bad_frame frame:%v",
				exampid, tag, conn.Session(),
				md.Message.Command) // Handle this ......
		}
		wh := md.Message.Headers
		for j := 0; j < len(wh)-1; j += 2 {
			ll.Printf("%stag:%s connsess:%s header:%s:%s\n",
				exampid, tag, conn.Session(),
				wh[j], wh[j+1])
		}
		if pbc > 0 {
			maxlen := pbc
			if len(md.Message.Body) < maxlen {
				maxlen = len(md.Message.Body)
			}
			ss := string(md.Message.Body[0:maxlen])
			ll.Printf("%stag:%s connsess:%s message_payload body:%s\n",
				exampid, tag, conn.Session(),
				ss)

		}
		// ACK the message just received.
		// Agiain we use a utility routine to handle the different requirements
		// of the protocol versions.
		sngecomm.HandleAck(conn, md.Message.Headers, id)
		ll.Printf("%stag:%s connsess:%s  ack_complete\n",
			exampid, tag, conn.Session())
	}
	// It is polite to unsubscribe, although unnecessary if a disconnect follows.
	// Again we use a utility routine to handle the different protocol level
	// requirements.
	sngecomm.HandleUnsubscribe(conn, d, id)
	ll.Printf("%stag:%s connsess:%s stomp_unsubscribe_complete\n",
		exampid, tag, conn.Session())

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s disconnect error:%v",
			exampid, tag, conn.Session(),
			e.Error()) // Handle this ......
	}

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, conn.Session(),
		time.Now().Sub(st))

}
