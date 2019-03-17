//
// Copyright © 2019 Guy M. Allard
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
package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo/senv"
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "putget: "
	ll      = log.New(os.Stdout, "PUGT ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	tag     = "putgemm"
	conn    *stompngo.Connection
	n       net.Conn
	e       error
	u2      = false
	wns     int64
	wd      time.Duration
)

func main() {
	st := time.Now()
	// Environment variable controls
	if os.Getenv("STOMP_2CONN") != "" {
		u2 = true
	}
	wt := os.Getenv("STOMP_WTIME")
	if wt != "" {
		_, err := strconv.ParseInt(wt, 10, 64)
		if err != nil {
			ll.Fatalf("%stag:%s connsess:%s conversion error:%v",
				exampid, tag, sngecomm.Lcs,
				e.Error()) // Handle this ......
		} else {
			d, err := time.ParseDuration(wt + "ms")
			if err != nil {
				ll.Fatalf("%stag:%s connsess:%s parse error:%v",
					exampid, tag, sngecomm.Lcs,
					e.Error()) // Handle this ......
			}
			wd = d
			wns = wd.Nanoseconds()
		}
	}
	// Standard example connect sequence
	n, conn, e = sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	ll.Printf("%stag:%s connsess:%s START nmsgs:%d\n",
		exampid, tag, conn.Session(), senv.Nmsgs())
	// Put
	qname := sngecomm.Dest()
	sh := stompngo.Headers{"destination", qname}
	amsg := "A message from putget"
	for i := 0; i < senv.Nmsgs(); i++ {
		e = conn.Send(sh, amsg)
		if e != nil {
			ll.Fatalln("Send error:", e)
		}
	}
	// Possible 2nd connection
	if u2 {
		// Standard example disconnect sequence
		e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
		if e != nil {
			ll.Fatalf("%stag:%s connsess:%s main_on_disconnect1 error:%v",
				exampid, tag, conn.Session(),
				e.Error()) // Handle this ......
		}
		if wns > 0 {
			time.Sleep(wd)
		}
		//
		n, conn, e = sngecomm.CommonConnect(exampid, tag, ll)
		if e != nil {
			ll.Fatalf("%stag:%s connsess:%s main_on_connect1 error:%v",
				exampid, tag, sngecomm.Lcs,
				e.Error()) // Handle this ......
		}
	} else if wns > 0 {
		time.Sleep(wd)
	}
	// Get
	id := "putget-subid1"
	sc := sngecomm.HandleSubscribe(conn, qname, id, sngecomm.AckMode())
	ll.Printf("%stag:%s connsess:%s subscribe_complete id:%v dest:%v\n",
		exampid, tag, conn.Session(),
		id, qname)
	//
	for i := 0; i < senv.Nmsgs(); i++ {
		md := <-sc
		if md.Error != nil {
			ll.Fatalf("%stag:%s connsess:%s recv_error dest:%s error:%v",
				exampid, tag, conn.Session(),
				qname, md.Error) // Handle this ......
		}
		ll.Printf("%stag:%s connsess:%s message is:%s\n",
			exampid, tag, conn.Session(),
			md.Message.BodyString())
	}
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
