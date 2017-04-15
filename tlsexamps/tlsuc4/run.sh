#!/bin/sh
#
set -x
CERTBASE=${CERTBASE:-/ad3/gma/ad3/sslwork/2017-01}
CACERT=${CACERT:-ca.crt}
CLICERT=${CLICERT:-client.crt}
CLIKEY=${CLIKEY:-client.key}
go build
./tlsuc4 -srvCAFile=$CERTBASE/$CACERT \
	-cliCertFile=$CERTBASE/$CLICERT \
	-cliKeyFile=$CERTBASE/$CLIKEY
set +x

