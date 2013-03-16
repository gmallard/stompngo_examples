//
// Copyright © 2011-2012 Guy M. Allard
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
	"fmt"
	"github.com/gmallard/stompngo"
	"os"
	"strconv"
)

var h10 = "localhost" // default 1.0 host 
var p10 = "61613"     // default 1.0 port (ActiveMQ on the author's machine)

var h11 = "localhost" // default 1.1 host 
var p11 = "62613"     // default 1.1 port (Apollo on the author's machine)

var h12 = "localhost" // default 1.2 host 
var p12 = "62613"     // default 1.2 port (Apollo on the author's machine)

var nmsgs = 1                         // Default number of messages to send
var dest = "/queue/snge.common.queue" // Default destination
var nqs = 1                           // Default number of queues for multi-queue demo(s)

// Override 1.0 Host and port for Dial if requested.
func HostAndPort10() (string, string) {
	he := os.Getenv("STOMP_HOST")
	if he != "" {
		h10 = he
	}
	pe := os.Getenv("STOMP_PORT")
	if pe != "" {
		p10 = pe
	}
	return h10, p10
}

// Override 1.1 Host and port for Dial if requested.
func HostAndPort11() (string, string) {
	he := os.Getenv("STOMP_HOST")
	if he != "" {
		h11 = he
	}
	pe := os.Getenv("STOMP_PORT")
	if pe != "" {
		p11 = pe
	}
	return h11, p11
}

// Override 1.2 Host and port for Dial if requested.
func HostAndPort12() (string, string) {
	he := os.Getenv("STOMP_HOST")
	if he != "" {
		h12 = he
	}
	pe := os.Getenv("STOMP_PORT")
	if pe != "" {
		p12 = pe
	}
	return h12, p12
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
	if d == "" {
		return "localhost"
	}
	return d
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
