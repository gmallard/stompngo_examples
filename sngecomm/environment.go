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
	"bytes"
	"log"
	"os"
	"strconv"
	//
	// "github.com/gmallard/stompngo"
)

var (
	//
	nqs  = 1               // Default number of queues for multi-queue demo(s)
	mdml = 1024 * 32       // Message data max length of variable message, 32K
	md   = make([]byte, 1) // Additional message data, primed during init()
	rc   = 1               // Receiver connection count, srmgor_1smrconn
	pbc  = 64              // Number of bytes to print (used in some
	// 																 // examples that receive).
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
}

// Number of queues
func Nqs() int {
	//
	if s := os.Getenv("STOMP_NQS"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Println("NQS conversion error", e)
		} else {
			nqs = int(i)
		}
	}
	return nqs
}

// Receiver connection count
func Recvconns() int {
	//
	if s := os.Getenv("STOMP_RECVCONNS"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Println("RECVCONNS conversion error", e)
		} else {
			rc = int(i)
		}
	}
	return rc
}

// Max Data Message Length
func Mdml() int {
	if s := os.Getenv("STOMP_MDML"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Println("MDML conversion error", e)
		} else {
			mdml = int(i)
			p := "_123456789ABCDEF"
			c := mdml / len(p)
			b := []byte(p)
			md = bytes.Repeat(b, c) // A long string
		}
	}
	return mdml
}

// Use profiling or not
func Pprof() bool {
	if am := os.Getenv("STOMP_PPROF"); am != "" {
		pprof = true
	}
	return pprof
}

// ACK mode for those examples that use it.
func AckMode() string {
	if am := os.Getenv("STOMP_ACKMODE"); am != "" {
		if am == "auto" || am == "client" || am == "client-individual" {
			ackMode = am
		} else {
			log.Println("ACKMODE error", am)
		}
	}
	return ackMode
}

// 2 Connection Buffer Size
func Conn2Buffer() int {
	if s := os.Getenv("STOMP_CONN2BUFFER"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Println("CONN2BUFFER conversion error", e)
		} else {
			conn2Buffer = int(i)
		}
	}
	return conn2Buffer
}

// Get Send Sleep Factor
func SendFactor() float64 {
	if s := os.Getenv("STOMP_SENDFACT"); s != "" {
		f, e := strconv.ParseFloat(s, 64)
		if nil != e {
			log.Println("SENDFACT conversion error", e)
		} else {
			sendFact = float64(f)
		}
	}
	return sendFact
}

// Get Recv Sleep Factor
func RecvFactor() float64 {
	if s := os.Getenv("STOMP_RECVFACT"); s != "" {
		f, e := strconv.ParseFloat(s, 64)
		if nil != e {
			log.Println("RECVFACT conversion error", e)
		} else {
			recvFact = float64(f)
		}
	}
	return recvFact
}

// Get partial string, random length
func Partial() []byte {
	r := int(ValueBetween(1, int64(mdml-1), 1.0))
	return md[0:r]
}

// Print Byte Count
func Pbc() int {
	if s := os.Getenv("STOMP_PBC"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Println("PBC conversion error", e)
		} else {
			pbc = int(i)
		}
	}
	return pbc
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
