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
	"net"
	"os"
	//
	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "conndisc: "
	ll      = log.New(os.Stdout, "ECNDS ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
)

// Connect to a STOMP broker and disconnect.
func main() {
	ll.Println(exampid + "starts ...")

	// Open a net connection
	h, p := senv.HostAndPort()
	hap := net.JoinHostPort(h, p)
	ll.Println(exampid+"will dial ...", hap)
	n, e := net.Dial("tcp", hap)
	if e != nil {
		ll.Fatalln("netdial", e) // Handle this ......
	}
	ll.Println(exampid+"dial complete ...", hap)

	// All stomp API methods require 'Headers'.  Stomp headers are key/value
	// pairs.  The stompngo package implements them using a string slice.
	ch := sngecomm.ConnectHeaders()

	// Get a stomp connection.  Parameters are:
	// a) the opened net connection
	// b) the Connect Headers
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		ll.Fatalln("sngConnect", e) // Handle this ......
	}
	ll.Println(exampid+"stomp connect complete, protocol level is:", conn.Protocol())

	// Show connect response
	ll.Println(exampid+"connect response:", conn.ConnectResponse)

	if senv.Heartbeats() != "" {
		ll.Println(exampid+"heart-beat send:", conn.SendTickerInterval())
		ll.Println(exampid+"heart-beat receive:", conn.ReceiveTickerInterval())
	}

	// *NOTE* your application functionaltiy goes here!

	// Stomp disconnects are not required, but highly recommended.
	// We use empty headers here for the DISCONNECT.  This will trigger the
	// disconnect receipt described below.
	//
	// If you do *not* want a disconnect receipt use:
	// stompngo.Headers{"noreceipt", "true"}
	//
	e = conn.Disconnect(stompngo.Headers{})
	if e != nil {
		ll.Fatalln("sngDisconnect", e) // Handle this ......
	}
	ll.Println(exampid + "stomp disconnect complete ...")

	// After a DISCONNECT there is (by default) a 'Disconnect
	// Receipt' available.  The disconnect receipt is an instance of
	// stompngo.MessageData.
	ll.Println(exampid+"disconnect receipt", conn.DisconnectReceipt)

	// Close the net connection.
	e = n.Close()
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	ll.Println(exampid + "network close complete ...")

	ll.Println(exampid + "ends ...")
}
