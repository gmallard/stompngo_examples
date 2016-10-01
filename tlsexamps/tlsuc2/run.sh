#!/bin/sh
#
set -x
go build
./tlsuc2 -srvCAFile=/ad3/gma/sslwork/2016-02/ca.crt
set +x

