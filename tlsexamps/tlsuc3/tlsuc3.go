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
Connect and Disconnect from a STOMP broker with a TLS connection, use case 3.

	TLS Use Case 3 - broker *does* authenticate client, client does *not* authenticate broker

	Subcase 3.A - Message broker configuration does *not* require client authentication

	- Expect connection success

	Subcase 3.B - Message broker configuration *does* require client authentication

	- Expect connection success if the broker can authenticate the client certificate

	Example use might be:

		go build
		STOMP_PORT=61611 ./tlsuc3 -cliCertFile=/ad3/gma/sslwork/2016-02/client.crt -cliKeyFile=/ad3/gma/sslwork/2016-02/client.key

*/
package main

import (
	"crypto/tls"
	"flag"
	"log"
	"os"
	"time"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid     = "tlsuc3:"
	tc          *tls.Config
	cliCertFile string
	cliKeyFile  string
	ll          = log.New(os.Stdout, "TLSU3 ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	tag = "tuc3main"
)

func init() {
	flag.StringVar(&cliCertFile, "cliCertFile", "DUMMY_CERT", "Name of client cert file")
	flag.StringVar(&cliKeyFile, "cliKeyFile", "DUMMY_KEY", "Name of client key file")
}

// Connect to a STOMP broker using TLS and disconnect.
func main() {

	st := time.Now()

	ll.Printf("%stag:%s connsess:%s starts\n",
		exampid, tag, sngecomm.Lcs)

	flag.Parse() // Parse flags
	ll.Printf("%stag:%s connsess:%s main_using_cliCertFile:%s\n",
		exampid, tag, sngecomm.Lcs,
		cliCertFile)
	ll.Printf("%stag:%s connsess:%s main_using_cliKeyFile:%s\n",
		exampid, tag, sngecomm.Lcs,
		cliKeyFile)

	// TLS Configuration.
	tc = new(tls.Config)
	tc.InsecureSkipVerify = true // Do *not* check the broker's certificate
	// Be polite, allow SNI (Server Virtual Hosting)
	tc.ServerName = senv.Host()
	// Finish TLS Config initialization, so broker can authenticate client.
	// cc -> tls.Certificate
	cc, e := tls.LoadX509KeyPair(cliCertFile, cliKeyFile)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_load_pair error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}
	// Add cert to config
	tc.Certificates = append(tc.Certificates, cc)
	// This is OK, but does not seem to be required
	tc.BuildNameToCertificate() // Build names map

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
