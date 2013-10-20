//
// Copyright © 2011-2013 Guy M. Allard
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
	"github.com/gmallard/stompngo"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

var (
	host     = "localhost" // default host
	port     = "61613"     // default port
	protocol = "1.1"       // Default protocol level
	login    = "guest"     // Default login
	passcode = "guest"     // Default passcode
	vhost    = "localhost"
	//
	nmsgs = 1                          // Default number of messages to send
	dest  = "/queue/snge.common.queue" // Default destination
	nqs   = 1                          // Default number of queues for multi-queue demo(s)
	scc   = 1                          // Default subscribe channel capacity
	md    = ""                         // Additional message data, primed during init()
	mdml  = 1024 * 32                  // Message data max length of variable message, 32K
)

// Initialization
func init() {
	p := "_123456789ABCDEF"
	c := mdml / len(p)
	md = strings.Repeat(p, c) // A long string
}

// Get partial string, random length
func Partial() string {
	r := int(ValueBetween(1, int64(mdml-1)))
	return string(md[0:r])
}

// Override default protocol level
func Protocol() string {
	p := os.Getenv("STOMP_PROTOCOL")
	if p != "" {
		protocol = p
	}
	return protocol
}

// Override Host and port for Dial if requested.
func HostAndPort() (string, string) {
	he := os.Getenv("STOMP_HOST")
	if he != "" {
		host = he
	}
	pe := os.Getenv("STOMP_PORT")
	if pe != "" {
		port = pe
	}
	return host, port
}

// Override login
func Login() string {
	l := os.Getenv("STOMP_LOGIN")
	if l != "" {
		login = l
	}
	if l == "NONE" {
		login = ""
	}
	return login
}

// Override passcode
func Passcode() string {
	p := os.Getenv("STOMP_PASSCODE")
	if p != "" {
		passcode = p
	}
	if p == "NONE" {
		passcode = ""
	}
	return passcode
}

// Provide connect headers
func ConnectHeaders() stompngo.Headers {
	h := stompngo.Headers{}
	l := Login()
	if l != "" {
		h = h.Add("login", l)
	}
	pc := Passcode()
	if pc != "" {
		h = h.Add("passcode", pc)
	}
	//
	p := Protocol()
	if p != stompngo.SPL_10 { // 1.1 and 1.2
		h = h.Add("accept-version", p).Add("host", Vhost())
	}
	return h
}

// Number of messages to send
func Nmsgs() int {
	c := os.Getenv("STOMP_NMSGS")
	if c == "" {
		return nmsgs
	}
	n, e := strconv.ParseInt(c, 10, 0)
	if e != nil {
		fmt.Printf("NMSGS Conversion error: %v\n", e)
		return nmsgs
	}
	return int(n)
}

// Number of queues to use
func Nqs() int {
	c := os.Getenv("STOMP_NQS")
	if c == "" {
		return nqs
	}
	n, e := strconv.ParseInt(c, 10, 0)
	if e != nil {
		fmt.Printf("NQS Conversion error: %v\n", e)
		return nqs
	}
	return int(n)
}

// Subscribe Channel Capacity
func SubChanCap() int {
	c := os.Getenv("STOMP_SUBCHANCAP")
	if c == "" {
		return scc
	}
	n, e := strconv.ParseInt(c, 10, 0)
	if e != nil {
		fmt.Printf("SUBCHANCAP Conversion error: %v\n", e)
		return scc
	}
	return int(n)
}

// Destination to send to
func Dest() string {
	d := os.Getenv("STOMP_DEST")
	if d == "" {
		return dest
	}
	return d
}

// Does receive wait to simulate message processing
func RecvWait() bool {
	f := os.Getenv("STOMP_NORECVW")
	if f == "" {
		return true
	}
	return false
}

// Does send wait to simulate message building
func SendWait() bool {
	f := os.Getenv("STOMP_NOSENDW")
	if f == "" {
		return true
	}
	return false
}

// True if persistent messages are desired.
func Persistent() bool {
	f := os.Getenv("STOMP_PERSISTENT")
	if f == "" {
		return false
	}
	return true
}

// True if max procs are to be set
func SetMAXPROCS() bool {
	f := os.Getenv("STOMP_SETMAXPROCS")
	if f == "" {
		return false
	}
	return true
}

// Virtual Host Name to use
func Vhost() string {
	d := os.Getenv("STOMP_VHOST")
	if d != "" {
		vhost = d
	}
	return vhost
}

// Show connection metrics.
func ShowStats(exampid, tag string, conn *stompngo.Connection) {
	r := conn.FramesRead()
	br := conn.BytesRead()
	w := conn.FramesWritten()
	bw := conn.BytesWritten()
	s := conn.Running().Seconds()
	n := conn.Running().Nanoseconds()
	fmt.Println(exampid, tag, "frame read count", r)
	fmt.Println(exampid, tag, "bytes read", br)
	fmt.Println(exampid, tag, "frame write count", w)
	fmt.Println(exampid, tag, "bytes written", bw)
	fmt.Println(exampid, tag, "current duration(ns)", n)
	fmt.Printf("%s %s %s %20.6f\n", exampid, tag, "current duration(sec)", s)
	fmt.Printf("%s %s %s %20.6f\n", exampid, tag, "frame reads/sec", float64(r)/s)
	fmt.Printf("%s %s %s %20.6f\n", exampid, tag, "bytes read/sec", float64(br)/s)
	fmt.Printf("%s %s %s %20.6f\n", exampid, tag, "frame writes/sec", float64(w)/s)
	fmt.Printf("%s %s %s %20.6f\n", exampid, tag, "bytes written/sec", float64(bw)/s)
}

// Get a value between min amd max
func ValueBetween(min, max int64) int64 {
	br, _ := rand.Int(rand.Reader, big.NewInt(max-min)) // Ignore errors here
	return br.Add(big.NewInt(min), br).Int64()
}

// Dump a TLS Configuration Struct
func DumpTLSConfig(c *tls.Config, n *tls.Conn) {
	fmt.Println()
	fmt.Printf("Rand: %v\n", c.Rand)
	fmt.Printf("Time: %v\n", c.Time)
	fmt.Printf("Certificates: %v\n", c.Certificates)
	fmt.Printf("NameToCertificate: %v\n", c.NameToCertificate)
	fmt.Printf("RootCAs: %v\n", c.RootCAs)
	fmt.Printf("NextProtos: %v\n", c.NextProtos)
	fmt.Printf("ServerName: %v\n", c.ServerName)
	fmt.Printf("ClientAuth: %v\n", c.ClientAuth)
	fmt.Printf("ClientCAs: %v\n", c.ClientCAs)
	fmt.Printf("CipherSuites: %v\n", c.CipherSuites)
	fmt.Printf("PreferServerCipherSuites: %v\n", c.PreferServerCipherSuites)
	fmt.Printf("SessionTicketsDisabled: %v\n", c.SessionTicketsDisabled)
	fmt.Printf("SessionTicketKey: %v\n", c.SessionTicketKey)

	// Idea Embellished From:
	// https://groups.google.com/forum/#!topic/golang-nuts/TMNdOxugbTY
	cs := n.ConnectionState()
	fmt.Println("HandshakeComplete:", cs.HandshakeComplete)
	fmt.Println("DidResume:", cs.DidResume)
	fmt.Println("CipherSuite:", cs.CipherSuite)
	fmt.Println("NegotiatedProtocol:", cs.NegotiatedProtocol)
	fmt.Println("NegotiatedProtocolIsMutual:", cs.NegotiatedProtocolIsMutual)
	fmt.Println("ServerName:", cs.ServerName)
	// Portions of any Peer Certificates present
	certs := cs.PeerCertificates
	if certs == nil || len(certs) < 1 {
		fmt.Println("Could not get server's certificate from the TLS connection.")
		fmt.Println()
		return
	}
	fmt.Println("Server Certs:")
	for i, cert := range certs {
		fmt.Printf("Certificate chain:%d\n", i)
		fmt.Printf("Common Name:%s\n", cert.Subject.CommonName)
		fmt.Printf("Alternate Name:%v\n", cert.DNSNames)
		fmt.Printf("Valid Not Before:%s\n", cert.NotBefore.Local().String())
		fmt.Println("" + strings.Repeat("=", 80) + "\n")
	}

	fmt.Println()
}

// Handle a subscribe for the different protocol levels
func Subscribe(c *stompngo.Connection, d, i, a string) <-chan stompngo.MessageData {
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

// Handle a unsubscribe for the different protocol levels
func Unsubscribe(c *stompngo.Connection, d, i string) {
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
