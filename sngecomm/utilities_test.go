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

package sngecomm

import (
	"os"
	"testing"
	// "github.com/gmallard/stompngo"
)

type headersData struct {
	key   string
	value string
}

var headersTests10 = []headersData{
	{"login", "guest"},
	{"passcode", "guest"},
}

var headersTests1x = []headersData{
	{"login", "guest"},
	{"passcode", "guest"},
	{"accept-version", "1.1"},
	{"host", "testhost"},
	{"heart-beat", "100,500"},
}

/*
	Test utilities defaults.  Must run with no environment variables set.
*/
func TestConnectHeaders10(t *testing.T) {
	h := ConnectHeaders()
	for _, v := range headersTests10 {
		val, ok := h.Contains(v.key)
		if !ok {
			t.Errorf("missing key [%s], expected [%t], got [%t]\n",
				v.key, true, ok)
		}
		//
		if v.value != val {
			t.Errorf("incorrect value, expected [%s], got [%s]\n",
				v.value, val)
		}
	}
}

/*
	Test utilities.
*/
func TestConnectHeaders1x(t *testing.T) {
	//
	_ = os.Setenv("STOMP_HOST", "testhost")
	_ = os.Setenv("STOMP_PROTOCOL", "1.1")
	_ = os.Setenv("STOMP_HEARTBEATS", "100,500")
	//
	h := ConnectHeaders()
	for _, v := range headersTests1x {
		val, ok := h.Contains(v.key)
		if !ok {
			t.Errorf("missing key [%s], expected [%t], got [%t]\n",
				v.key, true, ok)
		}
		//
		if v.value != val {
			t.Errorf("incorrect value, expected [%s], got [%s]\n",
				v.value, val)
		}
	}
}
