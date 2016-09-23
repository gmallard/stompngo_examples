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
Connect and Disconnect from a STOMP broker with a TLS connection, use case 1.

	TLS Use Case 1 - client does *not* authenticate broker.

	Subcase 1.A - Message broker configuration does *not* require client authentication

	- Expect connection success

	Subcase 1.B - Message broker configuration *does* require client authentication

	- Expect connection failure (broker must be sent a valid client certificate)

	Example use might be:

		go build
		STOMP_PORT=61611 ./tlsuc1

*/
package main

import (
	"crypto/tls"
	"log"
	"os"
	"time"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "tlsuc1: "
	tc      *tls.Config

	ll = log.New(os.Stdout, "TLSU1 ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	tag = "tuc1main"
)

// Connect to a STOMP broker using TLS and disconnect.
func main() {

	st := time.Now()

	ll.Printf("%stag:%s connsess:%s starts\n",
		exampid, tag, sngecomm.Lcs)

	// TLS Configuration.
	tc = new(tls.Config)
	tc.InsecureSkipVerify = true // Do *not* check the server's certificate

	// Standard example TLS connect sequence
	n, conn, e := sngecomm.CommonTLSConnect(exampid, tag, ll, tc)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	nc := n.(*tls.Conn)
	sngecomm.DumpTLSConfig(exampid, tc, nc)

	// *NOTE* application specific functionaltiy starts here!
	// For you to add.
	// *NOTE* application specific functionaltiy ends here!

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, conn.Session(),
		time.Now().Sub(st))

}
