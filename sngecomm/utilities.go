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
	"net"
	"os"
	"strings"
	//
	"github.com/gmallard/stompngo"
	"github.com/gmallard/stompngo/senv"
)

var (
	llu = log.New(os.Stdout, "UTIL ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	Lcs = "NotAvailable"
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
	llu.Printf("%s v1:%v v2:%v v3:%v\n", exampid, tag, "frame read count", r)
	llu.Printf("%s v1:%v v2:%v v3:%v\n", exampid, tag, "bytes read", br)
	llu.Printf("%s v1:%v v2:%v v3:%v\n", exampid, tag, "frame write count", w)
	llu.Printf("%s v1:%v v2:%v v3:%v\n", exampid, tag, "bytes written", bw)
	llu.Printf("%s v1:%v v2:%v v3:%v\n", exampid, tag, "current duration(ns)", n)
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
	llu.Printf("%s TLSConfig:\n", exampid)
	llu.Printf("%s Rand:%v\n", exampid, c.Rand)
	llu.Printf("%s Time:%v\n", exampid, c.Time)
	llu.Printf("%s Certificates:%v\n", exampid, c.Certificates)
	llu.Printf("%s NameToCertificate:%v\n", exampid, c.NameToCertificate)
	llu.Printf("%s RootCAs:%v\n", exampid, c.RootCAs)
	llu.Printf("%s NextProtos:%v\n", exampid, c.NextProtos)
	llu.Printf("%s ServerName:%v\n", exampid, c.ServerName)
	llu.Printf("%s ClientAuth:%v\n", exampid, c.ClientAuth)
	llu.Printf("%s ClientCAs:%v\n", exampid, c.ClientCAs)
	llu.Printf("%s CipherSuites:%v\n", exampid, c.CipherSuites)
	llu.Printf("%s PreferServerCipherSuites:%v\n", exampid, c.PreferServerCipherSuites)
	llu.Printf("%s SessionTicketsDisabled:%v\n", exampid, c.SessionTicketsDisabled)
	llu.Printf("%s SessionTicketKey:%v\n", exampid, c.SessionTicketKey)

	// Idea Embelluished From:
	// https://groups.google.com/forum/#!topic/golang-nuts/TMNdOxugbTY
	cs := n.ConnectionState()
	llu.Printf("%s HandshakeComplete:%v\n", exampid, cs.HandshakeComplete)
	llu.Printf("%s DidResume:%v\n", exampid, cs.DidResume)
	llu.Printf("%s CipherSuite:%d(0x%X)\n", exampid, cs.CipherSuite, cs.CipherSuite)
	llu.Printf("%s NegotiatedProtocol:%v\n", exampid, cs.NegotiatedProtocol)
	llu.Printf("%s NegotiatedProtocolIsMutual:%v\n", exampid, cs.NegotiatedProtocolIsMutual)
	// llu.Printf("%s ServerName:%v\n", exampid, cs.ServerName) // Server side only
	// Portions of any Peer Certificates present
	certs := cs.PeerCertificates
	if certs == nil || len(certs) < 1 {
		llu.Printf("Could not get server's certificate from the TLS connection.\n")
		return
	}
	llu.Printf("%s Server Certs:\n", exampid)
	for i, cert := range certs {
		llu.Printf("%s Certificate chain:%d\n", exampid, i)
		llu.Printf("%s Common Name:%s\n", exampid, cert.Subject.CommonName)
		//
		llu.Printf("%s Subject Alternative Names (DNSNames):\n", exampid)
		for idx, dnsn := range cert.DNSNames {
			llu.Printf("%s \tNumber:%d, DNS Name:%s\n", exampid, idx+1, dnsn)
		}
		//
		llu.Printf("%s Subject Alternative Names (Emailaddresses):\n", exampid)
		for idx, enn := range cert.EmailAddresses {
			llu.Printf("%s \tNumber:%d, DNS Name:%s\n", exampid, idx+1, enn)
		}
		//
		llu.Printf("%s Subject Alternative Names (IPAddresses):\n", exampid)
		for idx, ipadn := range cert.IPAddresses {
			llu.Printf("%s \tNumber:%d, DNS Name:%v\n", exampid, idx+1, ipadn)
		}
		//
		llu.Printf("%s Valid Not Before:%s\n", exampid, cert.NotBefore.Local().String())
		llu.Printf("%s Valid Not After:%s\n", exampid, cert.NotAfter.Local().String())
		llu.Println(strings.Repeat("=", 80))
	}

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
		llu.Fatalf("v1:%v v2:%v\n", "subscribe invalid protocol level, should not happen")
	}
	//
	r, e := c.Subscribe(h)
	if e != nil {
		llu.Fatalf("v1:%v v2:%v\n", "subscribe failed", e)
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
		llu.Fatalf("v1:%v v2:%v\n", "unsubscribe invalid protocol level, should not happen")
	}
	e := c.Unsubscribe(sbh)
	if e != nil {
		llu.Fatalf("v1:%v v2:%v d:%v\n", "unsubscribe failed", e, d)
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
		llu.Fatalf("v1:%v v2:%v\n", "ack invalid protocol level, should not happen")
	}
	e := c.Ack(ah)
	if e != nil {
		llu.Fatalf("v1:%v v2:%v v3:%v\n", "ack failed", e, c.Protocol())
	}
	return
}

func ShowRunParms(exampid string) {
	llu.Printf("%s v1:%v v2:%v\n", exampid, "HOST", os.Getenv("STOMP_HOST"))
	llu.Printf("%s v1:%v v2:%v\n", exampid, "PORT", os.Getenv("STOMP_PORT"))
	llu.Printf("%s v1:%v v2:%v\n", exampid, "PROTOCOL", senv.Protocol())
	llu.Printf("%s v1:%v v2:%v\n", exampid, "VHOST", senv.Vhost())
	llu.Printf("%s v1:%v v2:%v\n", exampid, "NQS", Nqs())
	llu.Printf("%s v1:%v v2:%v\n", exampid, "NMSGS", senv.Nmsgs())
	llu.Printf("%s v1:%v v2:%v\n", exampid, "SUBCHANCAP", senv.SubChanCap())
	llu.Printf("%s v1:%v v2:%v\n", exampid, "RECVFACT", RecvFactor())
	llu.Printf("%s v1:%v v2:%v\n", exampid, "SENDFACT", SendFactor())
	llu.Printf("%s v1:%v v2:%v\n", exampid, "ACKMODE", AckMode())
}

// Return broker identity
func ServerIdent(c *stompngo.Connection) string {
	cdh := c.ConnectResponse
	sr, ok := cdh.Headers.Contains("server")
	if !ok {
		return "N/A"
	}
	return sr
}

// Common example connect logic
func CommonConnect(exampid, tag string, l *log.Logger) (net.Conn,
	*stompngo.Connection,
	error) {

	l.Printf("%stag:%s consess:%v common_connect_starts\n",
		exampid, tag, Lcs)

	// Set up the connection.
	h, p := senv.HostAndPort()
	hap := net.JoinHostPort(h, p)
	n, e := net.Dial("tcp", hap)
	if e != nil {
		return nil, nil, e
	}

	l.Printf("%stag:%s connsess:%s common_connect_host_and_port:%v\n",
		exampid, tag, Lcs,
		hap)

	// Create connect headers and connect to stompngo
	ch := ConnectHeaders()
	l.Printf("%stag:%s connsess:%s common_connect_headers headers:%v\n",
		exampid, tag, Lcs,
		ch)
	conn, e := stompngo.Connect(n, ch)
	if e != nil {
		return nil, nil, e
	}
	l.Printf("%stag:%s connsess:%s common_connect_complete host:%s vhost:%s protocol:%s server:%s\n",
		exampid, tag, conn.Session(),
		h, senv.Vhost(), conn.Protocol(), ServerIdent(conn))

	// Show connect response
	l.Printf("%stag:%s connsess:%s common_connect_response connresp:%v\n",
		exampid, tag, conn.Session(),
		conn.ConnectResponse)

	// Show heartbeat data (if heart beats are being used)
	if senv.Heartbeats() != "" {
		l.Printf("%stag:%s connsess:%s common_connect_heart_beat_send hbsend:%v\n",
			exampid, tag, conn.Session(),
			conn.SendTickerInterval())
		l.Printf("%stag:%s connsess:%s common_connect_heart_beat_recv hbrecv:%v\n",
			exampid, tag, conn.Session(),
			conn.ReceiveTickerInterval())
	}

	l.Printf("%stag:%s connsess:%s common_connect_local_addr:%s\n",
		exampid, tag, conn.Session(),
		n.LocalAddr().String())
	l.Printf("%stag:%s connsess:%s common_connect_remote_addr:%s\n",
		exampid, tag, conn.Session(),
		n.RemoteAddr().String())

	//
	return n, conn, nil
}

// Common example disconnect logic
func CommonDisconnect(n net.Conn, conn *stompngo.Connection,
	exampid, tag string,
	l *log.Logger) error {

	// Disconnect from the Stomp server
	e := conn.Disconnect(stompngo.Headers{})
	if e != nil {
		return e
	}
	l.Printf("%stag:%s consess:%v common_disconnect_complete local_addr:%s remote_addr:%s\n",
		exampid, tag, conn.Session(),
		n.LocalAddr().String(), n.RemoteAddr().String())

	// Close the network connection
	e = n.Close()
	if e != nil {
		return e
	}

	// Parting messages
	l.Printf("%stag:%s consess:%v common_disconnect_network_close_complete\n",
		exampid, tag, conn.Session())
	l.Printf("%stag:%s consess:%v common_disconnect_ends\n",
		exampid, tag, conn.Session())

	//
	return nil
}

// Common example TLS connect logic
func CommonTLSConnect(exampid, tag string, l *log.Logger,
	c *tls.Config) (net.Conn, *stompngo.Connection, error) {

	l.Printf("%stag:%s consess:%s common_tls_connect_starts\n",
		exampid, tag, Lcs)

	// Set up the connection.
	h, p := senv.HostAndPort()
	hap := net.JoinHostPort(h, p)
	n, e := net.Dial("tcp", hap)
	if e != nil {
		return nil, nil, e
	}

	c.ServerName = h // SNI

	nc := tls.Client(n, c) // Returns: *tls.Conn : implements net.Conn
	e = nc.Handshake()
	if e != nil {
		if e.Error() == "EOF" {
			l.Printf("%stag:%s consess:%s common_tls_handshake_EOF_Is_the_broker_port_TLS_enabled? port:%s\n",
				exampid, tag, Lcs,
				p)
		}
		l.Fatalf("%stag:%s consess:%s common_tls_handshake_failed error:%v\n",
			exampid, tag, Lcs,
			e.Error())
	}
	l.Printf("%stag:%s consess:%s common_tls_handshake_complete\n",
		exampid, tag, Lcs)

	l.Printf("%stag:%s connsess:%s common_tls_connect_host_and_port:%v\n",
		exampid, tag, Lcs,
		hap)

	// Create connect headers and connect to stompngo
	ch := ConnectHeaders()
	l.Printf("%stag:%s connsess:%s common_tls_connect_headers headers:%v\n",
		exampid, tag, Lcs,
		ch)
	conn, e := stompngo.Connect(nc, ch)
	if e != nil {
		return nil, nil, e
	}
	l.Printf("%stag:%s connsess:%s common_tls_connect_complete host:%s vhost:%s protocol:%s server:%s\n",
		exampid, tag, conn.Session(),
		h, senv.Vhost(), conn.Protocol(), ServerIdent(conn))

	// Show connect response
	l.Printf("%stag:%s connsess:%s common_tls_connect_response connresp:%v\n",
		exampid, tag, conn.Session(),
		conn.ConnectResponse)

	// Show heartbeat data (if heart beats are being used)
	if senv.Heartbeats() != "" {
		l.Printf("%stag:%s connsess:%s common_tls_connect_heart_beat_send hbsend:%v\n",
			exampid, tag, conn.Session(),
			conn.SendTickerInterval())
		l.Printf("%stag:%s connsess:%s common_tls_connect_heart_beat_recv hbrecv:%v\n",
			exampid, tag, conn.Session(),
			conn.ReceiveTickerInterval())
	}

	l.Printf("%stag:%s connsess:%s common_tls_connect_local_addr:%s\n",
		exampid, tag, conn.Session(),
		n.LocalAddr().String())
	l.Printf("%stag:%s connsess:%s common_tls_connect_remote_addr:%s\n",
		exampid, tag, conn.Session(),
		n.RemoteAddr().String())

	//
	return nc, conn, nil
}
