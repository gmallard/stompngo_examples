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

package main

import (
	//"fmt"
	"log"
	"net"
	"os"
	"time"
	//
	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	//"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	ll = log.New(os.Stdout, "4EVRD ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	exampid = "reader: "

	//
	n    net.Conn             // Network Connection
	conn *stompngo.Connection // Stomp Connection

	lhl = 44

	tag = "reader"
)

// A forever reader.
func main() {

	st := time.Now()

	sngecomm.ShowRunParms(exampid)

	// Standard example connect sequence
	var e error
	n, conn, e = sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		if conn != nil {
			ll.Printf("%stag:%s  connsess:%s Connect Response headers:%v body%s\n",
				exampid, tag, conn.Session(), conn.ConnectResponse.Headers,
				string(conn.ConnectResponse.Body))
		}
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	//**
	qn := 1
	ltag := tag + "-receiver"

	id := stompngo.Uuid() // A unique subscription ID
	d := sngecomm.Dest()

	ll.Printf("%stag:%s connsess:%s queue_info id:%v d:%v qnum:%v\n",
		exampid, ltag, conn.Session(),
		id, d, qn)
	// Subscribe
	sc := sngecomm.HandleSubscribe(conn, d, id, sngecomm.AckMode())
	ll.Printf("%stag:%s connsess:%s subscribe_complete id:%v d:%v qnum:%v\n",
		exampid, ltag, conn.Session(),
		id, d, qn)
	//
	var md stompngo.MessageData
	// Receive loop
	mc := 1
	for {
		ll.Println("========================== ", "Expecting Message:", mc,
			"==========================")
		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			// A RECEIPT or ERROR frame is unexpected here
			ll.Fatalf("%stag:%s connsess:%s bad_frame qnum:%v headers:%v body:%s",
				exampid, tag, conn.Session(),
				qn, md.Message.Headers, md.Message.Body) // Handle this ......
		}
		if md.Error != nil {
			ll.Fatalf("%stag:%s connsess:%s recv_error qnum:%v error:%v",
				exampid, tag, conn.Session(),
				qn, md.Error) // Handle this ......
		}

		mbs := string(md.Message.Body)
		ll.Printf("%stag:%s connsess:%s next_frame mc:%v command:%v headers:%v body:%v\n",
			exampid, tag, conn.Session(),
			mc, md.Message.Command, md.Message.Headers, mbs) // Handle this ......

		mc++
		if sngecomm.UseEOF() && mbs == sngecomm.EOF_MSG {
			ll.Printf("%stag:%s connsess:%s received EOF\n",
				exampid, tag, conn.Session())
			break
		}
	}

	//**

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_disconnect error:%v",
			exampid, tag, conn.Session(),
			e.Error()) // Handle this ......
	}

	sngecomm.ShowStats(exampid, tag, conn)

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, conn.Session(),
		time.Now().Sub(st))

}
