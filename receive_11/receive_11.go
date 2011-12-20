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
  "stompngo_examples/common"
)

var exampid = "receive_11: "

// Connect to a STOMP 1.1 broker, receive some messages and disconnect.
func main() {
  fmt.Println(exampid + "starts ...")

  // Set up the connection.
  h, p := sngecomm.HostAndPort11() // A 1.1 connection
  n, e := net.Dial("tcp", net.JoinHostPort(h, p))
  if e != nil {
    log.Fatalln(e)  // Handle this ......
  }
  fmt.Println(exampid + "dial complete ...")
  eh := stomp.Headers{}
  conn, e := stomp.Connect(n, eh)
  if e != nil {
    log.Fatalln(e)  // Handle this ......
  }
  fmt.Println(exampid + "stomp connect complete ...")

  // *NOTE* your application functionaltiy goes here!
  // With Stomp, you must SUBSCRIBE to a destination in order to receive.
  // Stomp 1.1 demands subscribing with a unique subscription id, and we
  // do that here.
  u := stomp.Uuid() // Use package convenience function for unique ID
  s := stomp.Headers{"destination", sngecomm.Dest(),
    "id", u} // subscribe/unsubscribe headers
  // Subscribe returns:
  // a) A channel of MessageData struct.  Note that with this package, the
  // channel is unique (based on the unique subscription id).
  // b) A possible error.  Always check for errors.  They can be logical
  // errors detected by the stomp package, or even hard network errors, for
  // example the broker just crashed.
  r, e := conn.Subscribe(s)
  if e != nil {
    log.Fatalln(e) // Handle this ...
  }
  fmt.Println(exampid + "stomp subscribe complete ...")
  // Read data from the returned channel
  for i := 1; i <= sngecomm.Nmsgs(); i++ {
    m := <-r
    fmt.Println(exampid + "channel read complete ...")
    // MessageData has two components:
    // a) a Message struct
    // b) an Error value.  Check the error value as usual
    if m.Error != nil {
      log.Fatalln(m.Error)  // Handle this
    }
    //
    fmt.Printf("Frame Type: %s\n", m.Message.Command)   // Will be MESSAGE or ERROR!
    h := m.Message.Headers
    for j := 0; j < len(h) -1; j+= 2 {
      fmt.Printf("Header: %s:%s\n", h[j], h[j+1])
    }
    fmt.Printf("Payload: %s\n", string(m.Message.Body)) // Data payload
    // One item to note:  Stomp 1.1 servers _must_ return a 'subscription'
    // header in each message, where the value of the 'subscription' header
    // matches the unique subscription id supplied on subscribe.  With _some_
    // stomp client libraries, this allows you to categorize messages by
    // 'subscription'.  That is not required with this package!!!  This
    // because of the unique MessageData channels returned by Subscribe.
    // But check that this is the case for demonstration purposes.
    if h.Value("subscription") != u {
      fmt.Printf("Error condition, expected [%s], received [%s]\n", u, h.Value("subscription"))
      log.Fatalln("Bad subscription header")
    }
  }
  // It is polite to unsubscribe, although unnecessary if a disconnect follows.
  // With Stomp 1.1, the same unique ID is required on UNSUBSCRIBE.  Failure
  // to provide it will result in an error return.
  e = conn.Unsubscribe(s)
  if e != nil {
    log.Fatalln(e)  // Handle this ...
  }
  fmt.Println(exampid + "stomp unsubscribe complete ...")

  // Disconnect and Close
  e = conn.Disconnect(eh)
  if e != nil {
    log.Fatalln(e)  // Handle this ......
  }
  fmt.Println(exampid + "stomp disconnect complete ...")
  e = n.Close()
  if e != nil {
    log.Fatalln(e)  // Handle this ......
  }
  fmt.Println(exampid + "network close complete ...")

  fmt.Println(exampid + "ends ...")
}


