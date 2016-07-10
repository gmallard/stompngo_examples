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
	ll.Printf("%s starts\n", exampid)

	// Open a net connection
	h, p := senv.HostAndPort()
	n, e := net.Dial("tcp", net.JoinHostPort(h, p))
	if e != nil {
		ll.Fatalf("%s Dial %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s dial_complete hap:%s\n", exampid, net.JoinHostPort(h, p))

	ch := sngecomm.ConnectHeaders()
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalf("%s Connect %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s connsess:%s protocol:%s connectresponse:%v\n",
		exampid, conn.Session(), conn.Protocol(), conn.ConnectResponse)

	// *NOTE* your application functionaltiy goes here!
	sh := stompngo.Headers{"destination", senv.Dest()}
	if senv.Persistent() {
		sh = sh.Add("persistent", "true")
	}
	ms := exampid + " message: "
	for i := 1; i <= senv.Nmsgs(); i++ {
		mse := ms + fmt.Sprintf("%d", i)
		ll.Printf("%s connsess:%s sending mse:%s\n",
			exampid, conn.Session(),
			mse)
		e := conn.Send(sh, mse)
		if e != nil {
			ll.Fatalf("%s Send %s\n",
				exampid, conn.Session(),
				e.Error())
		}
		ll.Printf("%s connsess:%s send_complete mse:%s\n",
			exampid, conn.Session(),
			mse)
		time.Sleep(100 * time.Millisecond)
	}

	// Disconnect from the Stomp server
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalf("%s Disconnect %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Printf("%s connsess:%s disconnect_complete\n",
		exampid, conn.Session())
	// Close the network connection
	e = n.Close()
	if e != nil {
		ll.Fatalf("%s Close %s\n", exampid, e.Error()) // Handle this ......
	}

	ll.Printf("%s connsess:%s net_close_complete\n",
		exampid, conn.Session())

	ll.Printf("%s connsess:%s ends\n",
		exampid, conn.Session())
}
