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

package main

import (
  "fmt"
  "log"
  "net"
  "stomp"
  . "stompngo_examples/common"
)

var exampid = "connect_10: "

// Connect to a STOMP 1.0 broker and disconnect.
func main() {
  fmt.Println(exampid + "starts ...")

  // Open a net connection
  h, p := HostAndPort10()
  n, e := net.Dial("tcp", net.JoinHostPort(h, p))
  if e != nil {
    log.Fatalln(e)  // Handle this ......
  }
  fmt.Println(exampid + "dial complete ...")

  // All stomp API methods require 'Headers'.  Stomp headers are key/value 
  // pairs.  The stomp package implements them using a string slice.
  //
  // Empty Headers are useful for a number of API method calls, and we
  // use them to connect to a Stomp 1.0 broker.
  eh := stomp.Headers{}

  // Get a stomp connection.  Parameters are:
  // a) the opened net connection
  // b) the (empty) Headers
  conn, e := stomp.Connect(n, eh)
  if e != nil {
    log.Fatalln(e)  // Handle this ......
  }
  fmt.Println(exampid + "stomp connect complete ...")

  // *NOTE* your application functionaltiy goes here!

  // Polite Stomp disconnects are not required, but highly recommended.
  // Empty headers again in this example.
  e = conn.Disconnect(eh)
  if e != nil {
    log.Fatalln(e)  // Handle this ......
  }
  fmt.Println(exampid + "stomp disconnect complete ...")

  // Close the net connection.  We ignore errors here.
  e = n.Close()
  if e != nil {
    log.Fatalln(e)  // Handle this ......
  }
  fmt.Println(exampid + "network close complete ...")

  fmt.Println(exampid + "ends ...")
}


