//
// Copyright © 2013-2014 Guy M. Allard
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
	// "time"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var exampid = "publish: "

// Connect to a STOMP broker, publish some messages and disconnect.
func main() {
	fmt.Println(sngecomm.ExampIdNow(exampid) + "starts ...")

	// Open a net connection
	h, p := sngecomm.HostAndPort()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid) + "dial complete ...")

	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid)+"stomp connect complete ...", conn.Protocol())

	fmt.Println(sngecomm.ExampIdNow(exampid)+"connected headers", conn.ConnectResponse.Headers)
	// *NOTE* your application functionaltiy goes here!
	s := stompngo.Headers{"destination", sngecomm.Dest(),
		"persistent", "true"} // send headers
	m := exampid + " message: "
	for i := 1; i <= sngecomm.Nmsgs(); i++ {
		t := m + fmt.Sprintf("%d", i)
		fmt.Println(sngecomm.ExampIdNow(exampid), "sending now:", t)
		e := conn.Send(s, t)
		if e != nil {
			log.Fatalln("bad send", e) // Handle this ...
		}
		fmt.Println(sngecomm.ExampIdNow(exampid), "send complete:", t)
		//		time.Sleep(16 * time.Millisecond)
		//		time.Sleep(1 * time.Millisecond) // DB Behind ~ 4 messages
		//		time.Sleep(64 * time.Millisecond) // DB OK
		//		time.Sleep(125 * time.Millisecond) // DB OK
		//		time.Sleep(250 * time.Millisecond) // DB OK
		//		time.Sleep(500 * time.Millisecond) // DB OK
		// time.Sleep(2 * time.Minute)
	}

	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid) + "stomp disconnect complete ...")
	// Close the network connection
	e = n.Close()
	if e != nil {
		log.Fatalln(e) // Handle this ......
	}
	fmt.Println(sngecomm.ExampIdNow(exampid) + "network close complete ...")

	fmt.Println(sngecomm.ExampIdNow(exampid) + "ends ...")
}
