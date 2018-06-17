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

package main

/*
	Provide package version information.  A nod to the concept of semver.

	Example:
		go run $GOPATH/src/github.com/gmallard/stompngo_examples/version.go

*/

import (
	"fmt"
	//
	"github.com/gmallard/stompngo_examples"
)

func main() {
	fmt.Printf(stompngo_examples.Version())
}
