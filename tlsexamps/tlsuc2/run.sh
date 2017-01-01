#!/bin/sh
#
set -x
CERTBASE=${CERTBASE:-/ad3/gma/ad3/sslwork/2016-02}
CACERT=${CACERT:-ca.crt}
go build
./tlsuc2 -srvCAFile=$CERTBASE/$CACERT
set +x

