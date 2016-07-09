#!/bin/sh
#
set -x
go build
./tlsuc2 -srvCAFile=/ad3/gma/sslwork/2013/TestCA.crt
set +x

