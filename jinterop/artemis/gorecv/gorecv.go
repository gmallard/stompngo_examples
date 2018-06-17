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
A message receiver, to demonstrate JMS interoperability.
*/
package main

import (
	"log"
	"os"
	"time"
	//
	"github.com/gmallard/stompngo"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "gorecv: "
	ll      = log.New(os.Stdout, "GOJRCVRT ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	nmsgs   = 1
	tag     = "jintartgr"
)

// Connect to a STOMP broker, receive some messages and disconnect.
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

	pbc := sngecomm.Pbc() // Print byte count

	// Setup Headers ...
	id := stompngo.Uuid() // Use package convenience function for unique ID
	d := "jms.queue.exampleQueue"
	sc := sngecomm.HandleSubscribe(conn, d, id, "auto")
	ll.Printf("%stag:%s connsess:%s stomp_subscribe_complete\n",
		exampid, tag, conn.Session())

	//
	var md stompngo.MessageData
	// Read data from the returned channel
	for i := 1; i <= nmsgs; i++ {
		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// Frames RECEIPT or ERROR not expected here
			ll.Fatalf("%stag:%s connsess:%s bad_frame error:%v",
				exampid, tag, conn.Session(),
				e.Error()) // Handle this ......
		}
		ll.Printf("%stag:%s connsess:%s channel_read_complete\n",
			exampid, tag, conn.Session())
		// MessageData has two components:
		// a) a Message struct
		// b) an Error value.  Check the error value as usual
		if md.Error != nil {
			ll.Fatalf("%stag:%s connsess:%s error_read error:%v",
				exampid, tag, conn.Session(),
				e.Error()) // Handle this ......
		}
		ll.Printf("%stag:%s connsess:%s frame_type cmd:%s\n",
			exampid, tag, conn.Session(),
			md.Message.Command)

		if md.Message.Command != stompngo.MESSAGE {
			ll.Fatalf("%stag:%s connsess:%s error_frame_type error:%v",
				exampid, tag, conn.Session(),
				e.Error()) // Handle this ......
		}
		wh := md.Message.Headers
		for j := 0; j < len(wh)-1; j += 2 {
			ll.Printf("%stag:%s connsess:%s Header:%s:%s\n",
				exampid, tag, conn.Session(),
				wh[j], wh[j+1])
		}
		if pbc > 0 {
			maxlen := pbc
			if len(md.Message.Body) < maxlen {
				maxlen = len(md.Message.Body)
			}
			ss := string(md.Message.Body[0:maxlen])
			ll.Printf("%stag:%s connsess:%s payload body:%s\n",
				exampid, tag, conn.Session(),
				ss)
		}
	}

	//
	sngecomm.HandleUnsubscribe(conn, d, id)
	ll.Printf("%stag:%s connsess:%s unsubscribe_complete\n",
		exampid, tag, conn.Session())

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_disconnect error:%v",
			exampid, tag, conn.Session(),
			e.Error()) // Handle this ......
	}

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, conn.Session(),
		time.Now().Sub(st))

}
