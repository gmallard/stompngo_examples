#!/bin/sh
#
# Copyright © 2012-2016 Guy M. Allard
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -x
#
./clean.sh
#
./compile.sh
#
./jsend.sh
./jrecv.sh
#
./gosend.sh
./gorecv.sh
#
./gosend.sh
./jrecv.sh
#
./jsend.sh
./gorecv.sh
#
STOMP_NOSCL=true ./gosend.sh
./jrecv.sh
#
./clean.sh
set +x

