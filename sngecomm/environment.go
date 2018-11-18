//
// Copyright © 2016-2018 Guy M. Allard
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
	"sync"
	//
	// "github.com/gmallard/stompngo"
)

var (
	//
	nqs     = 1               // Default number of queues for multi-queue demo(s)
	nqsLock sync.Mutex        // nqs variable lock
	mdml    = 1024 * 32       // Message data max length of variable message, 32K
	md      = make([]byte, 1) // Additional message data, primed during init()
	pbc     = 64              // Number of bytes to print (used in some examples that receive).

	ngors    = 1  // Number of go routines to use (publish)
	gorsleep = "" // If non-empty, go routines will sleep (publish)

	//
	sendFact = 1.0 // Send sleep time factor
	recvFact = 1.0 // Receive sleep time factor
	//
	ackMode = "auto" // The default ack mode
	//
	pprof = false // Do not do profiling

	// TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256 = 0xC0,0x2F
	// TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256 = 0xC0,0x2B
	// TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 = 0xC0,0x30
	cipherSuites = []uint16{
		0xc02f,
		0xc02b,
		0xc030,
	}

	useCustomCiphers = false // Set Custom Cipher Suite list

	memprof = "" // memory profile file
	cpuprof = "" // cpu profile file
)

const (
	// EOFMsg is the EOF message body
	EOFMsg = "STOMP_EOF"
)

// Initialization
func init() {
	p := "_123456789ABCDEF"
	c := mdml / len(p)
	b := []byte(p)
	md = bytes.Repeat(b, c) // A long string
	//
	memprof = os.Getenv("STOMP_MEMPROF")
	cpuprof = os.Getenv("STOMP_CPUPROF")
}

// Ngors sets the number of go routines
func Ngors() int {
	//
	if s := os.Getenv("STOMP_NGORS"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Printf("v1:%v v2:%v\n", "NGORS conversion error", e)
		} else {
			ngors = int(i)
		}
	}
	return ngors
}

// Nqs sets the number of queues
func Nqs() int {
	//
	nqsLock.Lock()
	defer nqsLock.Unlock()
	if s := os.Getenv("STOMP_NQS"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Printf("v1:%v v2:%v\n", "NQS conversion error", e)
		} else {
			nqs = int(i)
		}
	}
	return nqs
}

// Mdml sets the Max Data Message Length
func Mdml() int {
	if s := os.Getenv("STOMP_MDML"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Printf("v1:%v v2:%v\n", "MDML conversion error", e)
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

// Pprof indicates whether to use profiling or not
func Pprof() bool {
	if am := os.Getenv("STOMP_PPROF"); am != "" {
		pprof = true
	}
	return pprof
}

// Memprof returns the memory profile file name
func Memprof() string {
	return memprof
}

// Cpuprof returns the CPU profile file name
func Cpuprof() string {
	return cpuprof
}

// AckMode returns an ACK mode value for those examples that use it.
func AckMode() string {
	if am := os.Getenv("STOMP_ACKMODE"); am != "" {
		if am == "auto" || am == "client" || am == "client-individual" {
			ackMode = am
		} else {
			log.Printf("v1:%v v2:%v\n", "ACKMODE error", am)
		}
	}
	return ackMode
}

// SendFactor returns the send sleep factor
func SendFactor() float64 {
	if s := os.Getenv("STOMP_SENDFACT"); s != "" {
		f, e := strconv.ParseFloat(s, 64)
		if nil != e {
			log.Printf("v1:%v v2:%v\n", "SENDFACT conversion error", e)
		} else {
			sendFact = float64(f)
		}
	}
	return sendFact
}

// RecvFactor returns the recv sleep factor
func RecvFactor() float64 {
	if s := os.Getenv("STOMP_RECVFACT"); s != "" {
		f, e := strconv.ParseFloat(s, 64)
		if nil != e {
			log.Printf("v1:%v v2:%v\n", "RECVFACT conversion error", e)
		} else {
			recvFact = float64(f)
		}
	}
	return recvFact
}

// Partial returns the partial byte slice for logging, random length
func Partial() []byte {
	r := int(ValueBetween(1, int64(mdml-1), 1.0))
	return md[0:r]
}

// PartialSubstr returns the partial string for logging, fixed length
func PartialSubstr(l int) []byte {
	return md[0:l]
}

// Pbc returns the byte count to log
func Pbc() int {
	if s := os.Getenv("STOMP_PBC"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Printf("v1:%v v2:%v\n", "PBC conversion error", e)
		} else {
			pbc = int(i)
		}
	}
	return pbc
}

// Gorsleep returns an indication of whether go routines will sleep or not
func Gorsleep() string {
	gorsleep = os.Getenv("STOMP_GORSLEEP")
	return gorsleep
}

// RecvWait indicates whether to wait in receives to simulate message processing
func RecvWait() bool {
	f := os.Getenv("STOMP_RECVWAIT")
	if f == "" {
		return true
	}
	return false
}

// SendWait indicates whether to wait in sends to simulate message processing
func SendWait() bool {
	f := os.Getenv("STOMP_SENDWAIT")
	if f == "" {
		return true
	}
	return false
}

// SetMAXPROCS returns true if max procs are to be set
func SetMAXPROCS() bool {
	f := os.Getenv("STOMP_SETMAXPROCS")
	if f == "" {
		return false
	}
	return true
}

// UseCustomCiphers returns true if custom ciphers are to be used
func UseCustomCiphers() bool {
	f := os.Getenv("STOMP_USECUSTOMCIPHERS")
	if f == "" {
		return useCustomCiphers
	}
	useCustomCiphers = true
	return useCustomCiphers
}

// CustomCiphers returns a slice of custom ciphers
func CustomCiphers() []uint16 {
	if UseCustomCiphers() {
		return cipherSuites
	}
	return []uint16{}
}

// Logger returns an indication of whether to do logging
func Logger() string {
	return os.Getenv("STOMP_LOGGER")
}

// UseEOF returns true if en EOF message is to be used.
func UseEOF() bool {
	if os.Getenv("STOMP_USEEOF") != "" {
		return true
	}
	return false
}
