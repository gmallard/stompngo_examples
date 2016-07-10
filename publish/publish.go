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
Publish messages to a STOMP broker.

	Examples:

		# Publish to a broker with all defaults:
		# Host is "localhost"
		# Port is 61613
		# Login is "guest"
		# Passcode is "guest
		# Virtual Host is "localhost"
		# Protocol is 1.1
		go run publish.go

		# Publish to a broker using STOMP protocol level 1.0:
		STOMP_PROTOCOL=1.0 go run publish.go

		# Publish to a broker using a custom host and port:
		STOMP_HOST=tjjackson STOMP_PORT=62613 go run publish.go

		# Publish to a broker using a custom port and virtual host:
		STOMP_PORT=41613 STOMP_VHOST="/" go run publish.go

		# Publish to a broker using a custom login and passcode:
		STOMP_LOGIN="userid" STOMP_PASSCODE="t0ps3cr3t" go run publish.go

*/
package main

import (
	"fmt"
	"log"
	"net"
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
	exampid = "publish: "
	ll      = log.New(os.Stdout, "EPUB ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

// Connect to a STOMP broker, publish some messages and disconnect.
func main() {
	ll.Println(exampid + "starts ...")

	// Open a net connection
	h, p := senv.HostAndPort()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalln(e) // Handle this ......
	}
	ll.Println(exampid+"dial complete ...",
		net.JoinHostPort(h, p))

	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalln(e) // Handle this ......
	}
	ll.Println(exampid+"stomp connect complete ...", conn.Protocol())

	ll.Println(exampid+"connected headers", conn.ConnectResponse.Headers)
	// *NOTE* your application functionaltiy goes here!
	sh := stompngo.Headers{"destination", senv.Dest()}
	if senv.Persistent() {
		sh = sh.Add("persistent", "true")
	}
	ms := exampid + " message: "
	for i := 1; i <= senv.Nmsgs(); i++ {
		t := ms + fmt.Sprintf("%d", i)
		ll.Println(exampid, "sending now:", t)
		e := conn.Send(sh, t)
		if e != nil {
			ll.Fatalln("bad send", e) // Handle this ...
		}
		ll.Println(exampid, "send complete:", t)
		//		time.Sleep(16 * time.Millisecond)
		//		time.Sleep(1 * time.Millisecond) // DB Behind ~ 4 messages
		//		time.Sleep(64 * time.Millisecond) // DB OK
		//		time.Sleep(125 * time.Millisecond) // DB OK
		//		time.Sleep(250 * time.Millisecond) // DB OK
		//		time.Sleep(500 * time.Millisecond) // DB OK
		//		time.Sleep(2 * time.Minute)
		time.Sleep(100 * time.Millisecond)
	}

	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalln(e) // Handle this ......
	}
	ll.Println(exampid + "stomp disconnect complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		ll.Fatalln(e) // Handle this ......
	}
	ll.Println(exampid + "network close complete ...")

	ll.Println(exampid + "ends ...")
}
