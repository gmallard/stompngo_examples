//
// Copyright © 2016-2018 Guy M. Alluard
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
		hb := senv.Heartbeats()
		if hb != "" {
			h = h.Add("heart-beat", hb)
		}
	}
	//

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
	llu.Printf("%stag:%s frame_read_count:%v\n", exampid, tag, r)
	llu.Printf("%stag:%s bytes_read:%v\n", exampid, tag, br)
	llu.Printf("%stag:%s frame_write_count:%v\n", exampid, tag, w)
	llu.Printf("%stag:%s bytes_written:%v\n", exampid, tag, bw)
	llu.Printf("%stag:%s current_duration(ns):%v\n", exampid, tag, n)

	llu.Printf("%stag:%s current_duration(sec):%20.6f\n", exampid, tag, s)
	llu.Printf("%stag:%s frame_reads/sec:%20.6f\n", exampid, tag, float64(r)/s)
	llu.Printf("%stag:%s bytes_read/sec:%20.6f\n", exampid, tag, float64(br)/s)
	llu.Printf("%stag:%s frame_writes/sec:%20.6f\n", exampid, tag, float64(w)/s)
	llu.Printf("%stag:%s bytes_written/sec:%20.6f\n", exampid, tag, float64(bw)/s)
}

// Get a value between min amd max
func ValueBetween(min, max int64, fact float64) int64 {
	rt, _ := rand.Int(rand.Reader, big.NewInt(max-min)) // Ignore errors here
	return int64(fact * float64(min+rt.Int64()))
}

// Dump a TLS Configuration Struct
func DumpTLSConfig(exampid string, c *tls.Config, n *tls.Conn) {
	llu.Printf("%s TLSConfig:\n", exampid)
	llu.Printf("%s Rand:%#v\n", exampid, c.Rand)
	if c.Time != nil {
		llu.Printf("%s Time:%v\n", exampid, c.Time())
	} else {
		llu.Printf("%s Time:%v\n", exampid, nil)
	}
	llu.Printf("%s Certificates:%#v\n", exampid, c.Certificates)
	llu.Printf("%s NameToCertificate:%#v\n", exampid, c.NameToCertificate)
	llu.Printf("%s RootCAs:%#v\n", exampid, c.RootCAs)
	llu.Printf("%s NextProtos:%v\n", exampid, c.NextProtos)
	llu.Printf("%s ServerName:%v\n", exampid, c.ServerName)
	llu.Printf("%s ClientAuth:%v\n", exampid, c.ClientAuth)
	llu.Printf("%s ClientCAs:%v#\n", exampid, c.ClientCAs)
	llu.Printf("%s CipherSuites:%#v\n", exampid, c.CipherSuites)
	llu.Printf("%s PreferServerCipherSuites:%v\n", exampid, c.PreferServerCipherSuites)
	llu.Printf("%s SessionTicketsDisabled:%v\n", exampid, c.SessionTicketsDisabled)
	llu.Printf("%s SessionTicketKey:%#v\n", exampid, c.SessionTicketKey)

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
	if len(certs) == 1 {
		llu.Printf("%s There is %d Server Cert:\n", exampid, len(certs))
	} else {
		llu.Printf("%s There are %d Server Certs:\n", exampid, len(certs))
	}

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
	if cv, ok := h.Contains(stompngo.HK_RECEIPT); ok {
		ah = ah.Add(stompngo.HK_RECEIPT, cv)
	}
	e := c.Ack(ah)
	if e != nil {
		llu.Fatalf("v1:%v v2:%v v3:%v\n", "ack failed", e, c.Protocol())
	}
	return
}

func ShowRunParms(exampid string) {
	llu.Printf("%sHOST:%v\n", exampid, os.Getenv("STOMP_HOST"))
	llu.Printf("%sPORT:%v\n", exampid, os.Getenv("STOMP_PORT"))
	llu.Printf("%sPROTOCOL:%v\n", exampid, senv.Protocol())
	llu.Printf("%sVHOST:%v\n", exampid, senv.Vhost())
	llu.Printf("%sNQS:%v\n", exampid, Nqs())
	llu.Printf("%sNMSGS:%v\n", exampid, senv.Nmsgs())
	llu.Printf("%sSUBCHANCAP:%v\n", exampid, senv.SubChanCap())
	llu.Printf("%sRECVFACT:%v\n", exampid, RecvFactor())
	llu.Printf("%sSENDFACT:%v\n", exampid, SendFactor())
	llu.Printf("%sRECVWAIT:%t\n", exampid, RecvWait())
	llu.Printf("%sSENDWAIT:%t\n", exampid, SendWait())
	llu.Printf("%sACKMODE:%v\n", exampid, AckMode())
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
		return nil, conn, e
	}
	SetLogger(conn) // Maybe set a connection logger
	l.Printf("%stag:%s connsess:%s common_connect_complete host:%s port:%s vhost:%s protocol:%s server:%s\n",
		exampid, tag, conn.Session(),
		h, p, senv.Vhost(), conn.Protocol(), ServerIdent(conn))

	// Show connect response
	l.Printf("%stag:%s connsess:%s common_connect_response connresp:%v\n",
		exampid, tag, conn.Session(),
		conn.ConnectResponse)

	// Heartbeat Data
	l.Printf("%stag:%s connsess:%s common_connect_heart_beat_send hbsend:%d\n",
		exampid, tag, conn.Session(),
		conn.SendTickerInterval())
	l.Printf("%stag:%s connsess:%s common_connect_heart_beat_recv hbrecv:%d\n",
		exampid, tag, conn.Session(),
		conn.ReceiveTickerInterval())

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
	SetLogger(conn)
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

// Example destination
func Dest() string {
	d := senv.Dest()
	if os.Getenv("STOMP_ARTEMIS") == "" {
		return d
	}
	pref := "jms.queue"
	if strings.Index(d, "topic") >= 0 {
		pref = "jms.topic"
	}
	return pref + strings.Replace(d, "/", ".", -1)
}

// Set Logger
func SetLogger(conn *stompngo.Connection) {
	if Logger() != "" {
		ul := log.New(os.Stdout, Logger()+" ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
		conn.SetLogger(ul)
	}
}
