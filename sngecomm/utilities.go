//
// Copyright © 2016 Guy M. Alluard
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
	"log"
	"math/big"
	"os"
	"strings"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo/senv"
)

var (
	llu = log.New(os.Stdout, "UTIL ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
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
	llu.Println(exampid, tag, "frame read count", r)
	llu.Println(exampid, tag, "bytes read", br)
	llu.Println(exampid, tag, "frame write count", w)
	llu.Println(exampid, tag, "bytes written", bw)
	llu.Println(exampid, tag, "current duration(ns)", n)
	llu.Printf("%s %s %s %20.6f\n", exampid, tag, "current duration(sec)", s)
	llu.Printf("%s %s %s %20.6f\n", exampid, tag, "frame reads/sec", float64(r)/s)
	llu.Printf("%s %s %s %20.6f\n", exampid, tag, "bytes read/sec", float64(br)/s)
	llu.Printf("%s %s %s %20.6f\n", exampid, tag, "frame writes/sec", float64(w)/s)
	llu.Printf("%s %s %s %20.6f\n", exampid, tag, "bytes written/sec", float64(bw)/s)
}

// Get a value between min amd max
func ValueBetween(min, max int64, fact float64) int64 {
	rt, _ := rand.Int(rand.Reader, big.NewInt(max-min)) // Ignore errors here
	return int64(fact * float64(min+rt.Int64()))
}

// Dump a TLS Configuration Struct
func DumpTLSConfig(exampid string, c *tls.Config, n *tls.Conn) {
	llu.Println()
	llu.Printf("%s Rand: %v\n", exampid, c.Rand)
	llu.Printf("%s Time: %v\n", exampid, c.Time)
	llu.Printf("%s Certificates: %v\n", exampid, c.Certificates)
	llu.Printf("%s NameToCertificate: %v\n", exampid, c.NameToCertificate)
	llu.Printf("%s RootCAs: %v\n", exampid, c.RootCAs)
	llu.Printf("%s NextProtos: %v\n", exampid, c.NextProtos)
	llu.Printf("%s ServerName: %v\n", exampid, c.ServerName)
	llu.Printf("%s ClientAuth: %v\n", exampid, c.ClientAuth)
	llu.Printf("%s ClientCAs: %v\n", exampid, c.ClientCAs)
	llu.Printf("%s CipherSuites: %v\n", exampid, c.CipherSuites)
	llu.Printf("%s PreferServerCipherSuites: %v\n", exampid, c.PreferServerCipherSuites)
	llu.Printf("%s SessionTicketsDisabled: %v\n", exampid, c.SessionTicketsDisabled)
	llu.Printf("%s SessionTicketKey: %v\n", exampid, c.SessionTicketKey)

	// Idea Embelluished From:
	// https://groups.google.com/forum/#!topic/golang-nuts/TMNdOxugbTY
	cs := n.ConnectionState()
	llu.Println(exampid, "HandshakeComplete:", cs.HandshakeComplete)
	llu.Println(exampid, "DidResume:", cs.DidResume)
	llu.Printf("%s %s %d(0x%X)\n", exampid, "CipherSuite:", cs.CipherSuite, cs.CipherSuite)
	llu.Println(exampid, "NegotiatedProtocol:", cs.NegotiatedProtocol)
	llu.Println(exampid, "NegotiatedProtocolIsMutual:", cs.NegotiatedProtocolIsMutual)
	llu.Println(exampid, "ServerName:", cs.ServerName)
	// Portions of any Peer Certificates present
	certs := cs.PeerCertificates
	if certs == nil || len(certs) < 1 {
		llu.Println("Could not get server's certificate from the TLS connection.")
		llu.Println()
		return
	}
	llu.Println(exampid, "Server Certs:")
	for i, cert := range certs {
		llu.Printf("Certificate chain: %d\n", i)
		llu.Printf("Common Name:%s \n", cert.Subject.CommonName)
		//
		llu.Printf("Subject Alternative Names (DNSNames):\n")
		for idx, dnsn := range cert.DNSNames {
			llu.Printf("\tNumber: %d, DNS Name: %s\n", idx+1, dnsn)
		}
		//
		llu.Printf("Subject Alternative Names (Emailaddresses):\n")
		for idx, enn := range cert.EmailAddresses {
			llu.Printf("\tNumber: %d, DNS Name: %s\n", idx+1, enn)
		}
		//
		llu.Printf("Subject Alternative Names (IPAddresses):\n")
		for idx, ipadn := range cert.IPAddresses {
			llu.Printf("\tNumber: %d, DNS Name: %v\n", idx+1, ipadn)
		}
		//
		llu.Printf("Valid Not Before: %s\n", cert.NotBefore.Local().String())
		llu.Printf("Valid Not After: %s\n", cert.NotAfter.Local().String())
		llu.Println("" + strings.Repeat("=", 80) + "\n")
	}

	llu.Println()
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
		llu.Fatalln("subscribe invalid protocol level, should not happen")
	}
	//
	r, e := c.Subscribe(h)
	if e != nil {
		llu.Fatalln("subscribe failed", e)
	}
	return r
}

// Handle a unsubscribe for the different protocol levels.
func HandleUnsubscribe(c *stompngo.Connection, d, i string) {
	sbh := stompngo.Headers{}
	//
	switch c.Protocol() {
	case stompngo.SPL_12:
		sbh = sbh.Add("id", i)
	case stompngo.SPL_11:
		sbh = sbh.Add("id", i)
	case stompngo.SPL_10:
		sbh = sbh.Add("destination", d)
	default:
		llu.Fatalln("unsubscribe invalid protocol level, should not happen")
	}
	e := c.Unsubscribe(sbh)
	if e != nil {
		llu.Fatalln("unsubscribe failed", e)
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
		llu.Fatalln("ack invalid protocol level, should not happen")
	}
	e := c.Ack(ah)
	if e != nil {
		llu.Fatalln("ack failed", e, c.Protocol())
	}
	return
}

func ShowRunParms(exampid string) {
	llu.Println(exampid, "HOST", os.Getenv("STOMP_HOST"))
	llu.Println(exampid, "PORT", os.Getenv("STOMP_PORT"))
	llu.Println(exampid, "PROTOCOL", senv.Protocol())
	llu.Println(exampid, "VHOST", senv.Vhost())
	llu.Println(exampid, "NQS", Nqs())
	llu.Println(exampid, "NMSGS", senv.Nmsgs())
	llu.Println(exampid, "SUBCHANCAP", senv.SubChanCap())
	llu.Println(exampid, "RECVFACT", RecvFactor())
	llu.Println(exampid, "SENDFACT", SendFactor())
	llu.Println(exampid, "CON2BUFFER", Conn2Buffer())
	llu.Println(exampid, "ACKMODE", AckMode())
}
