//
// Copyright © 2013-2016 Guy M. Allard
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

Subscribe and receive messages from a STOMP broker.

	Examples:

		# Subscribe to a broker with all defaults:
		# Host is "localhost"
		# Port is 61613
		# Login is "guest"
		# Passcode is "guest
		# Virtual Host is "localhost"
		# Protocol is 1.2
		go run subscribe.go

		# Subscribe to a broker using STOMP protocol level 1.1:
		STOMP_PROTOCOL=1.1 go run subscribe.go

		# Subscribe to a broker using a custom host and port:
		STOMP_HOST=tjjackson STOMP_PORT=62613 go run subscribe.go

		# Subscribe to a broker using a custom port and virtual host:
		STOMP_PORT=41613 STOMP_VHOST="/" go run subscribe.go

		# Subscribe to a broker using a custom login and passcode:
		STOMP_LOGIN="userid" STOMP_PASSCODE="t0ps3cr3t" go run subscribe.go

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
	exampid = "subscribe: "
	ll      = log.New(os.Stdout, "ESUBS ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

// Connect to a STOMP broker, subscribe and receive some messages and disconnect.
func main() {
	ll.Printf("%s v1:%v\n", exampid, "starts ...")

	// Set up the connection.
	h, p := senv.HostAndPort()
	hap := net.JoinHostPort(h, p)
	n, e := net.Dial("tcp", hap)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v v2:%v\n", exampid, "dial complete ...", hap)
	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v v2:%v\n", exampid, "stomp connect complete ...", conn.Protocol())

	pbc := sngecomm.Pbc() // Print byte count

	ll.Printf("%s v1:%v v2:%v\n", exampid, "connected headers", conn.ConnectResponse.Headers)
	// *NOTE* your application functionaltiy goes here!
	// With Stomp, you must SUBSCRIBE to a destination in order to receive.
	// Subscribe returns a channel of MessageData struct.
	// Here we use a common utility routine to handle the differing subscribe
	// requirements of each protocol level.
	d := senv.Dest()
	id := stompngo.Uuid()
	sc := sngecomm.HandleSubscribe(conn, d, id, "auto")
	ll.Printf("%s v1:%v\n", exampid, "stomp subscribe complete ...")
	// Read data from the returned channel
	var md stompngo.MessageData
	for i := 1; i <= senv.Nmsgs(); i++ {

		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// Frames RECEIPT or ERROR not expected here
			ll.Fatalf("%s v1:%v\n", exampid, md) // Handle this
		}

		ll.Printf("%s v1:%v\n", exampid, "channel read complete ...")
		ll.Printf("%s v1:%v v2:%v\n", exampid, "Message Number:", i)
		// MessageData has two components:
		// a) a Message struct
		// b) an Error value.  Check the error value as usual
		if md.Error != nil {
			ll.Fatalf("%s v1:%v\n", exampid, md.Error) // Handle this
		}
		//
		ll.Printf("Frame Type: %s\n", md.Message.Command) // Will be MESSAGE or ERROR!
		if md.Message.Command != stompngo.MESSAGE {
			ll.Fatalf("%s v1:%v\n", exampid, md) // Handle this ...
		}
		wh := md.Message.Headers
		for j := 0; j < len(wh)-1; j += 2 {
			ll.Printf("Header: %s:%s\n", wh[j], wh[j+1])
		}
		if pbc > 0 {
			maxlen := pbc
			if len(md.Message.Body) < maxlen {
				maxlen = len(md.Message.Body)
			}
			ss := string(md.Message.Body[0:maxlen])
			ll.Printf("Payload: %s\n", ss) // Data payload
		}
	}
	// It is polite to unsubscribe, although unnecessary if a disconnect follows.
	// Again we use a utility routine to handle the different protocol level
	// requirements.
	sngecomm.HandleUnsubscribe(conn, d, id)
	ll.Printf("%s v1:%v\n", exampid, "unsubscribe complete")

	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v\n", exampid, "stomp disconnect complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s v1:%v\n", exampid, "network close complete ...")

	ll.Printf("%s v1:%v\n", exampid, "ends ...")
}
