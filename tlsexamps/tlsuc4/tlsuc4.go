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
Connect and Disconnect from a STOMP broker with a TLS connection, use case 4.

	TLS Use Case 4 - broker *does* authenticate client, client *does* authenticate broker

	Subcase 4.A - Message broker configuration does *not* require client authentication

		- Expect connection success

	Subcase 4.B - Message broker configuration *does* require client authentication

		- Expect connection success if the broker can authenticate the client certificate

	Example use might be:

		go build
		STOMP_PORT=61611 ./tlsuc4 -srvCAFile=/ad3/gma/sslwork/2016-02/ca.crt -cliCertFile=/ad3/gma/sslwork/2016-02/client.crt -cliKeyFile=/ad3/gma/sslwork/2016-02/client.key

*/
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid     = "tlsuc4:"
	tc          *tls.Config
	srvCAFile   string // Name of file with broker's CA certificate, PEM format
	cliCertFile string
	cliKeyFile  string

	ll = log.New(os.Stdout, "TLSU4 ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	// TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 = 0xC0,0x2F
	// TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 = 0xC0,0x2B
	// TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 = 0xC0,0x30
	cipherSuites = []uint16{
		0xc02f,
		0xc02b,
		0xc030,
	}

	setCS = false // Set Custom Cipher Suite list

	tag = "tuc4main"
)

func init() {
	flag.StringVar(&srvCAFile, "srvCAFile", "DUMMY", "Name of file with broker CA certificate")
	flag.StringVar(&cliCertFile, "cliCertFile", "DUMMY_CERT", "Name of client cert file")
	flag.StringVar(&cliKeyFile, "cliKeyFile", "DUMMY_KEY", "Name of client key file")
}

// Connect to a STOMP broker using TLS and disconnect.
func main() {

	st := time.Now()

	ll.Printf("%stag:%s connsess:%s starts\n",
		exampid, tag, sngecomm.Lcs)

	flag.Parse() // Parse flags
	ll.Printf("%stag:%s connsess:%s main_using_srvCAFile:%s\n",
		exampid, tag, sngecomm.Lcs,
		srvCAFile)
	ll.Printf("%stag:%s connsess:%s main_using_cliCertFile:%s\n",
		exampid, tag, sngecomm.Lcs,
		cliCertFile)
	ll.Printf("%stag:%s connsess:%s main_using_cliKeyFile:%s\n",
		exampid, tag, sngecomm.Lcs,
		cliKeyFile)

	// TLS Configuration.
	tc = new(tls.Config)
	tc.InsecureSkipVerify = false // *Do* check the broker's certificate
	// Be polite, allow SNI (Server Virtual Hosting)
	tc.ServerName = senv.Host()
	// Finish TLS Config initialization, so client can authenticate broker,
	// and broker can authenticate client.

	// Usually one will use the default cipher suites that go provides.
	// However, if a custom cipher squite list is needed/required this
	// is how it is accomplished.
	if setCS { // Set custom cipher suite list
		tc.CipherSuites = append(tc.CipherSuites, cipherSuites...)
	}

	b, e := ioutil.ReadFile(srvCAFile) // Read broker's CA cert (PEM)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_read_file error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}
	k, _ := pem.Decode(b) // Decode PEM format
	if k == nil {
		ll.Fatalf("%stag:%s connsess:%s main_decode error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}
	//
	c, e := x509.ParseCertificate(k.Bytes) // Create *x509.Certificate
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_parse_cert error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}
	tc.RootCAs = x509.NewCertPool() // Create a cert "pool"
	tc.RootCAs.AddCert(c)           // Add the CA cert to the pool
	// Finish TLS Config initialization, so broker can authenticate client.
	cc, e := tls.LoadX509KeyPair(cliCertFile, cliKeyFile)
	if e != nil {
		ll.Fatalf("%s %s\n", exampid, e.Error()) // Handle this ......
	}
	tc.Certificates = append(tc.Certificates, cc) // Add cert
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
