//
// Copyright © 2013-2018 Guy M. Allard
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
Connect and Disconnect from a STOMP broker using TCP.

	Examples:

		# Connect to a broker with all defaults:
		# Host is "localhost"
		# Port is 61613
		# Login is "guest"
		# Passcode is "guest
		# Virtual Host is "localhost"
		# Protocol is 1.1
		go run conndisc.go

		# Connect to a broker using STOMP protocol level 1.0:
		STOMP_PROTOCOL=1.0 go run conndisc.go

		# Connect to a broker using a custom host and port:
		STOMP_HOST=tjjackson STOMP_PORT=62613 go run conndisc.go

		# Connect to a broker using a custom port and virtual host:
		STOMP_PORT=41613 STOMP_VHOST="/" go run conndisc.go

		# Connect to a broker using a custom login and passcode:
		STOMP_LOGIN="userid" STOMP_PASSCODE="t0ps3cr3t" go run conndisc.go

*/
package main

import (
	"log"
	"os"
	"time"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "conndisc: "
	ll      = log.New(os.Stdout, "ECNDS ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	tag     = "condiscmain"
)

// Connect to a STOMP broker and disconnect.
func main() {

	st := time.Now()

	// Standard example connect sequence
	n, conn, e := sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		if conn != nil {
			ll.Printf("%stag%s connsess:%s main_on_connect error Connect Response headers:%v body%s\n",
				exampid, tag, conn.Session(), conn.ConnectResponse.Headers,
				string(conn.ConnectResponse.Body))
		}
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	ll.Printf("%stag%s connsess:%s Connect Response headers:%v body%s\n",
		exampid, tag, conn.Session(), conn.ConnectResponse.Headers,
		string(conn.ConnectResponse.Body))

	// *NOTE* application specific functionaltiy starts here!
	// For you to add.
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
