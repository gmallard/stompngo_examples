//
// Copyright © 2016 Guy M. Allard
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
Package sngecomm provides common functionality used in the stompngo_examples
project.
*/
package sngecomm

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo/senv"
)

// Provide connect headers
func ConnectHeaders() stompngo.Headers {
	h := stompngo.Headers{}
	l := senv.Login()
	if l != "" {
		h = h.Add("login", l)
	}
	pc := senv.Passcode()
	if pc != "" {
		h = h.Add("passcode", pc)
	}
	//
	p := senv.Protocol()
	if p != stompngo.SPL_10 { // 1.1 and 1.2
		h = h.Add("accept-version", p).Add("host", senv.Vhost())
	}
	//
	hb := senv.Heartbeats()
	if hb != "" {
		h = h.Add("heart-beat", hb)
	}

	return h
}

// Show connection metrics.
func ShowStats(exampid, tag string, conn *stompngo.Connection) {
	r := conn.FramesRead()
	br := conn.BytesRead()
	w := conn.FramesWritten()
	bw := conn.BytesWritten()
	s := conn.Running().Seconds()
	n := conn.Running().Nanoseconds()
	fmt.Println(ExampIdNow(exampid), tag, "frame read count", r)
	fmt.Println(ExampIdNow(exampid), tag, "bytes read", br)
	fmt.Println(ExampIdNow(exampid), tag, "frame write count", w)
	fmt.Println(ExampIdNow(exampid), tag, "bytes written", bw)
	fmt.Println(ExampIdNow(exampid), tag, "current duration(ns)", n)
	fmt.Printf("%s %s %s %20.6f\n", ExampIdNow(exampid), tag, "current duration(sec)", s)
	fmt.Printf("%s %s %s %20.6f\n", ExampIdNow(exampid), tag, "frame reads/sec", float64(r)/s)
	fmt.Printf("%s %s %s %20.6f\n", ExampIdNow(exampid), tag, "bytes read/sec", float64(br)/s)
	fmt.Printf("%s %s %s %20.6f\n", ExampIdNow(exampid), tag, "frame writes/sec", float64(w)/s)
	fmt.Printf("%s %s %s %20.6f\n", ExampIdNow(exampid), tag, "bytes written/sec", float64(bw)/s)
}

// Get a value between min amd max
func ValueBetween(min, max int64, fact float64) int64 {
	rt, _ := rand.Int(rand.Reader, big.NewInt(max-min)) // Ignore errors here
	return int64(fact * float64(min+rt.Int64()))
}

// Dump a TLS Configuration Struct
func DumpTLSConfig(exampid string, c *tls.Config, n *tls.Conn) {
	fmt.Println()
	fmt.Printf("%s Rand: %v\n", ExampIdNow(exampid), c.Rand)
	fmt.Printf("%s Time: %v\n", ExampIdNow(exampid), c.Time)
	fmt.Printf("%s Certificates: %v\n", ExampIdNow(exampid), c.Certificates)
	fmt.Printf("%s NameToCertificate: %v\n", ExampIdNow(exampid), c.NameToCertificate)
	fmt.Printf("%s RootCAs: %v\n", ExampIdNow(exampid), c.RootCAs)
	fmt.Printf("%s NextProtos: %v\n", ExampIdNow(exampid), c.NextProtos)
	fmt.Printf("%s ServerName: %v\n", ExampIdNow(exampid), c.ServerName)
	fmt.Printf("%s ClientAuth: %v\n", ExampIdNow(exampid), c.ClientAuth)
	fmt.Printf("%s ClientCAs: %v\n", ExampIdNow(exampid), c.ClientCAs)
	fmt.Printf("%s CipherSuites: %v\n", ExampIdNow(exampid), c.CipherSuites)
	fmt.Printf("%s PreferServerCipherSuites: %v\n", ExampIdNow(exampid), c.PreferServerCipherSuites)
	fmt.Printf("%s SessionTicketsDisabled: %v\n", ExampIdNow(exampid), c.SessionTicketsDisabled)
	fmt.Printf("%s SessionTicketKey: %v\n", ExampIdNow(exampid), c.SessionTicketKey)

	// Idea Embellished From:
	// https://groups.google.com/forum/#!topic/golang-nuts/TMNdOxugbTY
	cs := n.ConnectionState()
	fmt.Println(ExampIdNow(exampid), "HandshakeComplete:", cs.HandshakeComplete)
	fmt.Println(ExampIdNow(exampid), "DidResume:", cs.DidResume)
	fmt.Printf("%s %s %d(0x%X)\n", ExampIdNow(exampid), "CipherSuite:", cs.CipherSuite, cs.CipherSuite)
	fmt.Println(ExampIdNow(exampid), "NegotiatedProtocol:", cs.NegotiatedProtocol)
	fmt.Println(ExampIdNow(exampid), "NegotiatedProtocolIsMutual:", cs.NegotiatedProtocolIsMutual)
	fmt.Println(ExampIdNow(exampid), "ServerName:", cs.ServerName)
	// Portions of any Peer Certificates present
	certs := cs.PeerCertificates
	if certs == nil || len(certs) < 1 {
		fmt.Println("Could not get server's certificate from the TLS connection.")
		fmt.Println()
		return
	}
	fmt.Println(ExampIdNow(exampid), "Server Certs:")
	for i, cert := range certs {
		fmt.Printf("Certificate chain: %d\n", i)
		fmt.Printf("Common Name:%s \n", cert.Subject.CommonName)
		//
		fmt.Printf("Subject Alternative Names (DNSNames):\n")
		for idx, dnsn := range cert.DNSNames {
			fmt.Printf("\tNumber: %d, DNS Name: %s\n", idx+1, dnsn)
		}
		//
		fmt.Printf("Subject Alternative Names (Emailaddresses):\n")
		for idx, enn := range cert.EmailAddresses {
			fmt.Printf("\tNumber: %d, DNS Name: %s\n", idx+1, enn)
		}
		//
		fmt.Printf("Subject Alternative Names (IPAddresses):\n")
		for idx, ipadn := range cert.IPAddresses {
			fmt.Printf("\tNumber: %d, DNS Name: %v\n", idx+1, ipadn)
		}
		//
		fmt.Printf("Valid Not Before: %s\n", cert.NotBefore.Local().String())
		fmt.Printf("Valid Not After: %s\n", cert.NotAfter.Local().String())
		fmt.Println("" + strings.Repeat("=", 80) + "\n")
	}

	fmt.Println()
}

// Handle a subscribe for the different protocol levels.
func HandleSubscribe(c *stompngo.Connection, d, i, a string) <-chan stompngo.MessageData {
	h := stompngo.Headers{"destination", d, "ack", a}
	//
	switch c.Protocol() {
	case stompngo.SPL_12:
		// Add required id header
		h = h.Add("id", i)
	case stompngo.SPL_11:
		// Add required id header
		h = h.Add("id", i)
	case stompngo.SPL_10:
		// Nothing else to do here
	default:
		log.Fatalln("subscribe invalid protocol level, should not happen")
	}
	//
	r, e := c.Subscribe(h)
	if e != nil {
		log.Fatalln("subscribe failed", e)
	}
	return r
}

// Handle a unsubscribe for the different protocol levels.
func HandleUnsubscribe(c *stompngo.Connection, d, i string) {
	h := stompngo.Headers{}
	//
	switch c.Protocol() {
	case stompngo.SPL_12:
		h = h.Add("id", i)
	case stompngo.SPL_11:
		h = h.Add("id", i)
	case stompngo.SPL_10:
		h = h.Add("destination", d)
	default:
		log.Fatalln("unsubscribe invalid protocol level, should not happen")
	}
	e := c.Unsubscribe(h)
	if e != nil {
		log.Fatalln("unsubscribe failed", e)
	}
	return
}

// Handle ACKs for the different protocol levels.
func HandleAck(c *stompngo.Connection, h stompngo.Headers, id string) {
	ah := stompngo.Headers{}
	//
	switch c.Protocol() {
	case stompngo.SPL_12:
		ah = ah.Add("id", h.Value("ack"))
	case stompngo.SPL_11:
		ah = ah.Add("message-id", h.Value("message-id")).Add("subscription", id)
	case stompngo.SPL_10:
		ah = ah.Add("message-id", h.Value("message-id"))
	default:
		log.Fatalln("unsubscribe invalid protocol level, should not happen")
	}
	e := c.Ack(ah)
	if e != nil {
		log.Fatalln("ack failed", e, c.Protocol())
	}
	return
}

func ShowRunParms(exampid string) {
	fmt.Println(ExampIdNow(exampid), "HOST", os.Getenv("STOMP_HOST"))
	fmt.Println(ExampIdNow(exampid), "PORT", os.Getenv("STOMP_PORT"))
	fmt.Println(ExampIdNow(exampid), "PROTOCOL", senv.Protocol())
	fmt.Println(ExampIdNow(exampid), "VHOST", senv.Vhost())
	fmt.Println(ExampIdNow(exampid), "NQS", Nqs())
	fmt.Println(ExampIdNow(exampid), "NMSGS", senv.Nmsgs())
	fmt.Println(ExampIdNow(exampid), "SUBCHANCAP", senv.SubChanCap())
	fmt.Println(ExampIdNow(exampid), "RECVFACT", RecvFactor())
	fmt.Println(ExampIdNow(exampid), "SENDFACT", SendFactor())
	fmt.Println(ExampIdNow(exampid), "CON2BUFFER", Conn2Buffer())
	fmt.Println(ExampIdNow(exampid), "ACKMODE", AckMode())
	fmt.Println(ExampIdNow(exampid), "RECVCONNS", Recvconns())
}
