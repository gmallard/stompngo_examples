//
// Copyright © 2016-2018 Guy M. Allard
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

Function: read the code :-).

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
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "varmGetter: "
	ll      = log.New(os.Stdout, "VRMGSC ", log.Ldate|log.Lmicroseconds)
	tag     = "vrmgmain"
	unsub   = true
	dodisc  = true
	ar      = false // Want ACK RECEIPT
	session = ""
	wlp     = "publish: message: " // Left part wanted
)

func init() {
	if os.Getenv("VMG_NOUNSUB") != "" {
		unsub = false
	}
	if os.Getenv("VMG_NODISC") != "" {
		dodisc = false
	}
	if os.Getenv("VMG_GETAR") != "" {
		ar = true
	}
}

// Connect to a STOMP broker, subscribe and receive some messages and disconnect.
func main() {

	st := time.Now()

	// Standard example connect sequence
	n, conn, e := sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}
	session = conn.Session()
	//******************
	nqs := sngecomm.Nqs()
	for qn := 1; qn <= nqs; qn++ {
		runNextQueue(qn, conn)
	}
	//******************

	// Standard example disconnect sequence
	if dodisc {
		e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
		if e != nil {
			ll.Fatalf("%stag:%s connsess:%s main_on_disconnect error:%v",
				exampid, tag, session,
				e.Error()) // Handle this ......
		}
		ll.Printf("%stag:%s connsess:%s disconnect_receipt:%v\n",
			exampid, tag, session,
			conn.DisconnectReceipt)
	} else {
		ll.Printf("%stag:%s connsess:%s skipping_disconnect\n",
			exampid, tag, session)
	}

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, session,
		time.Now().Sub(st))
}

func runNextQueue(qn int, conn *stompngo.Connection) {

	qns := fmt.Sprintf("%d", qn)    // string number of the queue
	conn.SetLogger(ll)              // stompngo logging
	pbc := sngecomm.Pbc()           // Print byte count
	d := senv.Dest() + qns          // Destination
	id := stompngo.Uuid()           // A unique name/id
	nmsgs := qn                     // int number of messages to get, same as queue number
	mns := fmt.Sprintf("%d", nmsgs) // string number of messages to get
	am := sngecomm.AckMode()        // ACK mode to use on SUBSCRIBE
	nfa := true                     // Need "final" ACK (possiby reset below)
	wh := stompngo.Headers{         // Starting SUBSCRIBE headers
		stompngo.StompPlusDrainAfter,
		mns} // Need a string here

	// Sanity check ACK mode
	if conn.Protocol() == stompngo.SPL_10 &&
		am == stompngo.AckModeClientIndividual {
		ll.Fatalf("%stag:%s connsess:%s invalid_ack_mode am:%v proto:%v\n",
			exampid, tag, session,
			am, conn.Protocol()) //
	}
	// Do not do final ACK if running ACKs are issued
	if am == stompngo.AckModeClientIndividual ||
		am == stompngo.AckModeAuto {
		nfa = false
	}

	// Show run parameters
	ll.Printf("%stag:%s connsess:%s run_parms\n\tqns:%v\n\tpbc:%v\n\td:%v\n\tid:%v\n\tnmsgs:%v\n\tam:%v\n\tnfa:%v\n\twh:%v\n",
		exampid, tag, session,
		qns, pbc, d, id, nmsgs, am, nfa, wh)

	// Run SUBSCRIBE
	sc := doSubscribe(conn, d, id, am, wh)
	ll.Printf("%stag:%s connsess:%s stomp_subscribe_complete\n",
		exampid, tag, session)

	var md stompngo.MessageData  // Message data from basic read
	var lmd stompngo.MessageData // Possible save (copy) of received data
	mc := 1                      // Initial message number

	// Loop for the requested number of messages
GetLoop:
	for {
		ll.Printf("%stag:%s connsess:%s start_of_read_loop mc:%v nmsgs:%v\n",
			exampid, tag, session, mc, nmsgs)

		mcs := fmt.Sprintf("%d", mc) // string number message count

		// Get something from the stompngo read routine
		select {
		case md = <-sc:
		case md = <-conn.MessageData:
			//
			if md.Message.Command == stompngo.RECEIPT {
				ll.Printf("%stag:%s connsess:%s have_receipt md:%v\n",
					exampid, tag, session,
					md)
				continue GetLoop
			}
			ll.Fatalf("%stag:%s connsess:%s ERROR_frame hdrs:%v body:%v\n",
				exampid, tag, session,
				md.Message.Headers, string(md.Message.Body)) // Handle this ......
		}

		// Save message data for possible use in the final ACK
		if mc == nmsgs && nfa {
			lmd = md // Save last message
		}

		// Basic loop logging
		ll.Printf("%stag:%s connsess:%s channel_read_complete qn:%d mc:%d\n",
			exampid, tag, session,
			qn, mc)
		ll.Printf("%stag:%s connsess:%s message_number:%v\n",
			exampid, tag, session,
			mc)

		// Check if reader returned any error
		if md.Error != nil {
			ll.Fatalf("%stag:%s connsess:%s error_read error:%v",
				exampid, tag, session,
				md.Error) // Handle this ......
		}

		// Show frame type
		ll.Printf("%stag:%s connsess:%s frame_type cmd:%s\n",
			exampid, tag, session,
			md.Message.Command)

		// Pure sanity check:  this should *never* happen based on logic
		// above.
		if md.Message.Command != stompngo.MESSAGE {
			ll.Fatalf("%stag:%s connsess:%s error_frame_type md:%v",
				exampid, tag, session,
				md) // Handle this ......
		}

		// Show Message Headers
		wh := md.Message.Headers
		for j := 0; j < len(wh)-1; j += 2 {
			ll.Printf("%stag:%s connsess:%s Header:%s:%s\n",
				exampid, tag, session,
				wh[j], wh[j+1])
		}
		// Show (part of) Message Body
		if pbc > 0 {
			maxlen := pbc
			if len(md.Message.Body) < maxlen {
				maxlen = len(md.Message.Body)
			}
			ss := string(md.Message.Body[0:maxlen])
			ll.Printf("%stag:%s connsess:%s payload body:%s\n",
				exampid, tag, session,
				ss)
		}

		// Sanity check this message payload
		wm := wlp + mcs // The left part plus the (string) meassage number]
		bm := string(md.Message.Body)
		if bm != wm {
			ll.Fatalf("%stag:%s connsess:%s error_message_payload\n\tGot %s\n\tWant%s\n",
				exampid, tag, session,
				bm, wm) // Handle this ......
		} else {
			ll.Printf("%stag:%s connsess:%s  matched_body_string\n%s\n%s\n",
				exampid, tag, session,
				bm, wm) // Handle this ......)
		}

		// Run individual ACK if required
		if am == stompngo.AckModeClientIndividual {
			wh := md.Message.Headers // Copy Headers
			if ar {                  // ACK receipt wanted
				wh = wh.Add(stompngo.HK_RECEIPT, "rwanted-"+mcs)
			}
			sngecomm.HandleAck(conn, wh, id)
			ll.Printf("%stag:%s connsess:%s  individual_ack_complete mc:%v headers:%v\n",
				exampid, tag, session,
				mc, md.Message.Headers)

		}

		// Check for end of loop condition
		if mc == nmsgs {
			break
		}

		// Increment loop/message counter
		mc++
	}

	// Issue the final ACK if needed
	if nfa {
		wh := lmd.Message.Headers // Copy Headers
		if ar {                   // ACK receipt wanted
			wh = wh.Add(stompngo.HK_RECEIPT, "rwanted-fin")
		}
		sngecomm.HandleAck(conn, wh, id)
		ll.Printf("%stag:%s connsess:%s  final_ack_complete\n",
			exampid, tag, session)
		if ar {
			getReceipt(conn)
		}
	}

	// Unsubscribe (may be skipped if requested)
	if unsub {
		sngecomm.HandleUnsubscribe(conn, d, id)
		ll.Printf("%stag:%s connsess:%s stomp_unsubscribe_complete\n",
			exampid, tag, session)
	} else {
		ll.Printf("%stag:%s connsess:%s skipping_unsubscribe\n",
			exampid, tag, session)
	}
}

// Handle a subscribe for the different protocol levels.
func doSubscribe(c *stompngo.Connection, d, id, a string, h stompngo.Headers) <-chan stompngo.MessageData {
	h = h.Add("destination", d).Add("ack", a)
	//
	switch c.Protocol() {
	case stompngo.SPL_12:
		// Add required id header
		h = h.Add("id", id)
	case stompngo.SPL_11:
		// Add required id header
		h = h.Add("id", id)
	case stompngo.SPL_10:
		// Nothing else to do here
	default:
		ll.Fatalf("v1:%v\n", "subscribe invalid protocol level, should not happen")
	}
	//
	r, e := c.Subscribe(h)
	if e != nil {
		ll.Fatalf("subscribe failed err:[%v]\n", e)
	}
	return r
}

// Get receipt
func getReceipt(conn *stompngo.Connection) {
	rd := <-conn.MessageData
	ll.Printf("%stag:%s connsess:%s have_receipt_sub md:%v\n",
		exampid, tag, session,
		rd)
}
