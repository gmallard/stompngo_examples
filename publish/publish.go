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
	tag     = "pubmain"
)

// Connect to a STOMP broker, publish some messages and disconnect.
func main() {

	st := time.Now()

	// Standard example connect sequence
	n, conn, e := sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	// *NOTE* application specific functionaltiy starts here!
	sh := stompngo.Headers{"destination", sngecomm.Dest()}
	ll.Printf("%stag:%s connsess:%s destination dest:%s\n",
		exampid, tag, conn.Session(),
		sngecomm.Dest())
	if senv.Persistent() {
		sh = sh.Add("persistent", "true")
	}
	ms := exampid + "message: "
	for i := 1; i <= senv.Nmsgs(); i++ {
		mse := ms + fmt.Sprintf("%d", i)
		ll.Printf("%stag:%s connsess:%s main_sending mse:~%s~\n",
			exampid, tag, conn.Session(),
			mse)
		e := conn.Send(sh, mse)
		if e != nil {
			ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
				exampid, tag, conn.Session(),
				e.Error()) // Handle this ......
		}
		ll.Printf("%stag:%s connsess:%s main_send_complete mse:~%s~\n",
			exampid, tag, conn.Session(),
			mse)
		time.Sleep(100 * time.Millisecond)
	}
	// *NOTE* application specific functionaltiy ends here!

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_disconnect error:%v",
			exampid, tag, conn.Session(),
			e.Error()) // Handle this ......
	}

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, conn.Session(),
		time.Now().Sub(st))
}
