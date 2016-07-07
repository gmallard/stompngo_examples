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
	"testing"
)

/*
	Test sngecomm defaults.  Must run with no environment variables set.
*/
func TestEnvironmentDefaults(t *testing.T) {
	n := Nqs()
	if n != 1 {
		t.Errorf("sngecomm Nqs, expected [%d], got [%d]\n", 1, n)
	}
	//
	n = Mdml()
	if n != 1024*32 {
		t.Errorf("sngecomm Mdml, expected [%d], got [%d]\n", 1024*32, n)
	}
	//
	b := Pprof()
	if b {
		t.Errorf("sngecomm Pprof, expected [%t], got [%t]\n", false, b)
	}
	//
	s := AckMode()
	if s != "auto" {
		t.Errorf("sngecomm AckMode, expected [%s], got [%s]\n", "auto", s)
	}
	//
	n = Conn2Buffer()
	if n != -1 {
		t.Errorf("sngecomm Conn2Buffer, expected [%d], got [%d]\n", -1, n)
	}
	//
	f := SendFactor()
	if f != 1.0 {
		t.Errorf("sngecomm SendFactor, expected [%v], got [%v]\n", 1.0, n)
	}
	//
	f = RecvFactor()
	if f != 1.0 {
		t.Errorf("sngecomm RecvFactor, expected [%v], got [%v]\n", 1.0, n)
	}
	//
	n = Pbc()
	if n != 64 {
		t.Errorf("sngecomm Pbc, expected [%v], got [%v]\n", 64, n)
	}
	//
	b = RecvWait()
	if !b {
		t.Errorf("sngecomm RecvWait, expected [%t], got [%t]\n", true, b)
	}
	//
	b = SendWait()
	if !b {
		t.Errorf("sngecomm SendWait, expected [%t], got [%t]\n", true, b)
	}
	//
	b = Persistent()
	if b {
		t.Errorf("sngecomm Persistent, expected [%t], got [%t]\n", false, b)
	}
	//
	b = SetMAXPROCS()
	if b {
		t.Errorf("sngecomm SetMAXPROCS, expected [%t], got [%t]\n", false, b)
	}

}
