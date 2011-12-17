//
// Copyright © 2011 Guy M. Allard
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

package sngecomm

import (
  "fmt"
  "os"
  "strconv"
)

var h10 = "localhost" // default 1.0 host 
var p10 = "61613"     // default 1.0 port (ActiveMQ on the author's machine)

var h11 = "localhost" // default 1.1 host 
var p11 = "62613"     // default 1.1 port (Apollo on the author's machine)

var nmsgs = 1         // Default number of messages to send
var dest = "/queue/snge.common.queue" // Default destination

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

// Destination to send to
func Dest() string {
  d := os.Getenv("STOMP_DEST")
  if d == "" {
    return dest
  }
  return d
}

