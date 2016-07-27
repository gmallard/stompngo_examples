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
	"os"
	"time"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "gosend: "
	ll      = log.New(os.Stdout, "GOJSND ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	nmsgs   = 1
	tag     = "jintamqjs"
)

// Connect to a STOMP broker, send some messages and disconnect.
func main() {

	st := time.Now()

	// Standard example connect sequence
	// Use AMQ port here
	n, conn, e := sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	// Suppress content length here, so JMS will treat this as a 'text' message.
	sh := stompngo.Headers{"destination", "/queue/allards.queue"}
	if os.Getenv("STOMP_NOSCL") != "true" {
		sh = sh.Add("suppress-content-length", "true")
	}
	if senv.Persistent() {
		sh = sh.Add("persistent", "true")
	}
	ms := exampid + " message: "
	for i := 1; i <= nmsgs; i++ {
		mse := ms + fmt.Sprintf("%d", i)
		e := conn.Send(sh, mse)
		if e != nil {
			ll.Fatalf("%stag:%s connsess:%s send_error error:%v\n",
				exampid, tag, conn.Session(),
				e.Error()) // Handle this ......
		}
		ll.Printf("%stag:%s connsess:%s send_complete\n",
			exampid, tag, conn.Session())
	}

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_disconnect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, conn.Session(),
		time.Now().Sub(st))

}
