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
Publish messages to a STOMP broker.

	Examples:

		# Publish to a broker with all defaults:
		# Host is "localhost"
		# Port is 61613
		# Login is "guest"
		# Passcode is "guest
		# Virtual Host is "localhost"
		# Protocol is 1.1
		go run publish.go

		# Publish to a broker using STOMP protocol level 1.0:
		STOMP_PROTOCOL=1.0 go run publish.go

		# Publish to a broker using a custom host and port:
		STOMP_HOST=tjjackson STOMP_PORT=62613 go run publish.go

		# Publish to a broker using a custom port and virtual host:
		STOMP_PORT=41613 STOMP_VHOST="/" go run publish.go

		# Publish to a broker using a custom login and passcode:
		STOMP_LOGIN="userid" STOMP_PASSCODE="t0ps3cr3t" go run publish.go

		# Important environment variables for this program are:

		# STOMP_NGORS - the number of go routines used to write to the
		# sepcified queues.

		# STOMP_NMSGS - the number of messages each go routine will write.

		# STOMP_NQS - The number of queues to write messages to.  If this
		# variable is absent, the value defaults to the value specified
		# for STOMP_NGORS.  If this value is specified, all go routines
		# are multi-plexed across this number of queues.

*/
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"sync"
	"time"
	//
	"github.com/gmallard/stompngo"
	// senv methods could be used in general by stompngo clients.
	"github.com/gmallard/stompngo/senv"
	// sngecomm methods are used specifically for these example clients.
	"github.com/gmallard/stompngo_examples/sngecomm"
)

var (
	exampid = "publish: "
	ll      = log.New(os.Stdout, "EPUB ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	tag     = "pubmain"
	wg      sync.WaitGroup
	conn    *stompngo.Connection
	n       net.Conn
	e       error
	nqs     int
	ngor    int
	// MNHDR is the message number, in message headers
	MNHDR  = "sng_msgnum"
	gorstr = 1 // Starting destination number
	//
	msfl  = true // Fixed message length
	mslen = 1024 // The fixed length
	msf   []byte
	//
	gorsl   = false
	gorslfb = false
	gorslfx = time.Duration(250 * time.Millisecond)
	gorslms = 250
	//
	max int64 = 1e9      // Max stagger time (nanoseconds)
	min       = max / 10 // Min stagger time (nanoseconds)
	// Sleep multipliers
	sf = 1.0
	rf = 1.0
)

func init() {
	ngor = sngecomm.Ngors()
	nqs = sngecomm.Nqs()

	// Options around message length:
	// 1) fixed length
	// 2) randomly variable length
	if os.Getenv("STOMP_VARMSL") != "" {
		msfl = false // Use randomly variable message lengths
	}
	if msfl {
		if s := os.Getenv("STOMP_FXMSLEN"); s != "" {
			i, e := strconv.ParseInt(s, 10, 32)
			if nil != e {
				log.Printf("v1:%v v2:%v\n", "FXMSLEN conversion error", e)
			} else {
				mslen = int(i) // The fixed length to use
			}
		}
		msf = sngecomm.PartialSubstr(mslen)
	}

	// Options controlling sleeps between message sends.  Options are:
	// 1) Don't sleep
	// 2) Sleep a fixed amount of time
	// 3) Sleep a random variable amount of time
	if os.Getenv("STOMP_DOSLEEP") != "" {
		gorsl = true // Do sleep
	}
	if os.Getenv("STOMP_FIXSLEEP") != "" {
		gorslfb = true // Do a fixed length sleep
	}
	if gorsl {
		var err error
		if s := os.Getenv("STOMP_SLEEPMS"); s != "" { // Fixed length milliseconds
			mss := fmt.Sprintf("%s", s) + "ms"
			gorslfx, err = time.ParseDuration(mss)
			if err != nil {
				log.Printf("v1:%v v2:%v v3:%v\n", "ParseDuration conversion error", mss, e)
			}
		}
	}
	//
	if s := os.Getenv("STOMP_GORNSTR"); s != "" {
		i, e := strconv.ParseInt(s, 10, 32)
		if nil != e {
			log.Printf("v1:%v v2:%v\n", "GORNSTR conversion error", e)
		} else {
			gorstr = int(i) // The fixed length to use
		}
	}
}
func runSends(gr int, qn int) {
	var err error
	qns := fmt.Sprintf("%d", qn)
	qname := sngecomm.Dest() + "." + qns
	sh := stompngo.Headers{"destination", qname}
	ll.Printf("%stag:%s connsess:%s destination dest:%s BEGIN_runSends %d\n",
		exampid, tag, conn.Session(),
		qname, gr)
	if senv.Persistent() {
		sh = sh.Add("persistent", "true")
	}
	sh = sh.Add(MNHDR, "0")
	mnhnum := sh.Index(MNHDR)
	ll.Printf("%stag:%s connsess:%s send headers:%v\n",
		exampid, tag, conn.Session(),
		sh)
	for i := 1; i <= senv.Nmsgs(); i++ {
		is := fmt.Sprintf("%d", i) // Next message number
		sh[mnhnum+1] = is          // Put message number in headers
		// Log send headers
		ll.Printf("%stag:%s connsess:%s main_sending gr:%d hdrs:%v\n",
			exampid, tag, conn.Session(),
			gr, sh)

		// Handle fixed or variable message length
		rml := 0
		if msfl {
			err = conn.SendBytes(sh, msf)
			rml = len(msf)
		} else {
			// ostr := string(sngecomm.Partial())
			// err = conn.Send(sh, ostr)
			oby := sngecomm.Partial()
			err = conn.SendBytes(sh, oby)
			rml = len(oby)
		}
		if err != nil {
			ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
				exampid, tag, conn.Session(),
				err.Error()) // Handle this ......
		}
		ll.Printf("%stag:%s connsess:%s main_send_complete gr:%d msfl:~%t~len:%d\n",
			exampid, tag, conn.Session(),
			gr, msfl, rml)

		// Handle sleep options
		if gorsl {
			if gorslfb {
				// Fixed time to sleep
				ll.Printf("%stag:%s connsess:%s gr:%d main_fixed sleep:~%v\n",
					exampid, tag, conn.Session(), gr, gorslfx)
				time.Sleep(gorslfx)
			} else {
				// Variable time to sleep
				dt := time.Duration(sngecomm.ValueBetween(min, max, sf))
				ll.Printf("%stag:%s connsess:%s gr:%d main_rand sleep:~%v\n",
					exampid, tag, conn.Session(), gr, dt)
				time.Sleep(dt)
			}
		}
	}
	if sngecomm.UseEOF() {
		sh := stompngo.Headers{"destination", qname}
		_ = conn.Send(sh, sngecomm.EOF_MSG)
		ll.Printf("%stag:%s connsess:%s gr:%d sent EOF [%s]\n",
			exampid, tag, conn.Session(), gr, sngecomm.EOF_MSG)
	}
	wg.Done() // signal a goroutine completion
}

// Connect to a STOMP broker, publish some messages and disconnect.
func main() {

	if sngecomm.Pprof() {
		if sngecomm.Cpuprof() != "" {
			ll.Printf("%stag:%s connsess:%s CPUPROF %s\n",
				exampid, tag, sngecomm.Lcs, sngecomm.Cpuprof())
			f, err := os.Create(sngecomm.Cpuprof())
			if err != nil {
				log.Fatal("could not create CPU profile: ", err)
			}
			if err := pprof.StartCPUProfile(f); err != nil {
				log.Fatal("could not start CPU profile: ", err)
			}
			defer pprof.StopCPUProfile()
		}
	}

	st := time.Now()

	// Standard example connect sequence
	n, conn, e = sngecomm.CommonConnect(exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_connect error:%v",
			exampid, tag, sngecomm.Lcs,
			e.Error()) // Handle this ......
	}

	ll.Printf("%stag:%s connsess:%s START gorstr:%d ngor:%d nqs:%d nmsgs:%d\n",
		exampid, tag, conn.Session(), gorstr, ngor, nqs, senv.Nmsgs())

	rqn := gorstr - 1
	for i := gorstr; i <= gorstr+ngor-1; i++ {
		wg.Add(1)
		rqn++
		if nqs > 1 && rqn > nqs {
			rqn = gorstr
		}
		go runSends(i, rqn)
	}
	wg.Wait()

	// Standard example disconnect sequence
	e = sngecomm.CommonDisconnect(n, conn, exampid, tag, ll)
	if e != nil {
		ll.Fatalf("%stag:%s connsess:%s main_on_disconnect error:%v",
			exampid, tag, conn.Session(),
			e.Error()) // Handle this ......
	}

	ll.Printf("%stag:%s connsess:%s main_elapsed:%v\n",
		exampid, tag, conn.Session(),
		time.Now().Sub(st))

	if sngecomm.Pprof() {
		if sngecomm.Memprof() != "" {
			ll.Printf("%stag:%s connsess:%s MEMPROF %s\n",
				exampid, tag, conn.Session(), sngecomm.Memprof())
			f, err := os.Create(sngecomm.Memprof())
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
			f.Close()
		}
	}

}
