//
// Copyright © 2011-2014 Guy M. Allard
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
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
	//
	"github.com/gmallard/stompngo"
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
	mdml  = 1024 * 32                  // Message data max length of variable message, 32K
	md    = make([]byte, 1)            // Additional message data, primed during init()
	rc    = 1                          // Receiver connection count, srmgor_1smrconn
	//
	sendFact float64 = 1.0 // Send sleep time factor
	recvFact float64 = 1.0 // Receive sleep time factor
	//
	conn2Buffer int = -1 // 2 connection buffer. < 0 means use queue size.
	//
	ackMode = "auto" // The default ack mode
	//
	pprof = false // Do not do profiling
)

// Initialization
func init() {
	p := "_123456789ABCDEF"
	c := mdml / len(p)
	b := []byte(p)
	md = bytes.Repeat(b, c) // A long string
	//
	if s := os.Getenv("STOMP_SENDFACT"); s != "" {
		f, e := strconv.ParseFloat(s, 64)
		if e == nil {
			sendFact = f
		}
	}
	//
	if s := os.Getenv("STOMP_RECVFACT"); s != "" {
		f, e := strconv.ParseFloat(s, 64)
		if e == nil {
			recvFact = f
		}
	}
	//
	if s := os.Getenv("STOMP_CONN2BUFFER"); s != "" {
		s, e := strconv.ParseInt(s, 10, 32)
		if e == nil {
			conn2Buffer = int(s)
		}
	}
	//
	if am := os.Getenv("STOMP_ACKMODE"); am != "" {
		if am == "auto" || am == "client" || am == "client-individual" {
			ackMode = am
		}
	}
	//
	if am := os.Getenv("STOMP_PPROF"); am != "" {
		pprof = true
	}
	//
	if s := os.Getenv("STOMP_MDML"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if e == nil {
			mdml = int(i)
		}
	}
	//
	if s := os.Getenv("STOMP_RECVCONNS"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if e == nil {
			rc = int(i)
		}
	}
}

// Receiver connection count
func Recvconns() int {
	return int(rc)
}

// Max Data Message Length
func Mdml() int {
	return int(mdml)
}

// Use profiling or not
func Pprof() bool {
	return pprof
}

// ACK mode for those examples that use it.
func AckMode() string {
	return ackMode
}

// 2 Connection Buffer Size
func Conn2Buffer() int {
	return conn2Buffer
}

// Timestamp example ids
func ExampIdNow(s string) string {
	return time.Now().String() + " " + s
}

// Get Send Sleep Factor
func SendFactor() float64 {
	return sendFact
}

// Get Recv Sleep Factor
func RecvFactor() float64 {
	return recvFact
}

// Get partial string, random length
func Partial() []byte {
	r := int(ValueBetween(1, int64(mdml-1), 1.0))
	return md[0:r]
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
		fmt.Printf("Certificate chain:%d\n", i)
		fmt.Printf("Common Name:%s\n", cert.Subject.CommonName)
		fmt.Printf("Alternate Name:%v\n", cert.DNSNames)
		fmt.Printf("Valid Not Before:%s\n", cert.NotBefore.Local().String())
		fmt.Println("" + strings.Repeat("=", 80) + "\n")
	}

	fmt.Println()
}

// Handle a subscribe for the different protocol levels.
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

// Handle a unsubscribe for the different protocol levels.
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

// Handle ACKs for the different protocol levels.
func Ack(c *stompngo.Connection, h stompngo.Headers, id string) {
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
	fmt.Println(ExampIdNow(exampid), "HOST", os.Getenv("STOMP_HOST"), "alt", host)
	fmt.Println(ExampIdNow(exampid), "PORT", os.Getenv("STOMP_PORT"), "alt", port)
	fmt.Println(ExampIdNow(exampid), "PROTOCOL", Protocol())
	fmt.Println(ExampIdNow(exampid), "VHOST", Vhost())
	fmt.Println(ExampIdNow(exampid), "NQS", Nqs())
	fmt.Println(ExampIdNow(exampid), "NMSGS", Nmsgs())
	fmt.Println(ExampIdNow(exampid), "SUBCHANCAP", SubChanCap())
	fmt.Println(ExampIdNow(exampid), "RECVFACT", RecvFactor())
	fmt.Println(ExampIdNow(exampid), "SENDFACT", SendFactor())
	fmt.Println(ExampIdNow(exampid), "CON2BUFFER", Conn2Buffer())
	fmt.Println(ExampIdNow(exampid), "ACKMODE", AckMode())
	fmt.Println(ExampIdNow(exampid), "RECVCONNS", Recvconns())
}
