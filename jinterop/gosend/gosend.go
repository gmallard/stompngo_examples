//
// Copyright © 2011-2016 Guy M. Allard
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
A message sender, to demonstrate JMS interoperability.
*/
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	//
	"github.com/gmallard/stompngo"
)

var (
	exampid = "gosend: "
	ll      = log.New(os.Stdout, "GOJSND ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	nmsgs   = 1
)

// Connect to a STOMP 1.2 broker, send some messages and disconnect.
func main() {
	ll.Println(exampid + "starts ...")

	// Open a net connection
	n, e := net.Dial("tcp", "localhost:61613")
	if e != nil {
		ll.Fatalln(e) // Handle this ......
	}
	ll.Println(exampid + "dial complete ...")

	// Connect to broker
	ch := stompngo.Headers{"login", "userr", "passcode", "passw0rd",
		"host", "localhost", "accept-version", "1.2"}
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalln(e) // Handle this ......
	}
	ll.Println(exampid + "stomp connect complete ...")

	// Suppress content length here, so JMS will treat this as a 'text' message.
	sh := stompngo.Headers{"destination", "/queue/allards.queue",
		"suppress-content-length", "true"} // send headers, suppress content-length
	ms := exampid + " message: "
	for i := 1; i <= nmsgs; i++ {
		t := ms + fmt.Sprintf("%d", i)
		e := conn.Send(sh, t)
		if e != nil {
			ll.Fatalln(e) // Handle this ...
		}
		ll.Println(exampid, "send complete:", t)
	}

	// Disconnect from the Stomp server
	dh := stompngo.Headers{}
	e = conn.Disconnect(dh)
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
