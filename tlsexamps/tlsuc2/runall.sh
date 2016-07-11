#!/bin/sh
#
set -x
go build
# Default AMQ port, not TSL enabled
./tlsuc2 -srvCAFile=/ad3/gma/sslwork/2013/TestCA.crt
#
STOMP_PORT=61612 ./tlsuc2 -srvCAFile=/ad3/gma/sslwork/2013/TestCA.crt
#
STOMP_PORT=62614 ./tlsuc2 -srvCAFile=/ad3/gma/sslwork/2013/TestCA.crt
set +x

