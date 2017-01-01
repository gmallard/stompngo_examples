#!/bin/sh
#
set -x
go build
CERTBASE=${CERTBASE:-/ad3/gma/ad3/sslwork/2016-02}
CLICERT=${CLICERT:-client.crt}
CLIKEY=${CLIKEY:-client.key}
#
./tlsuc3 -cliCertFile=$CERTBASE/$CLICERT \
	-cliKeyFile=$CERTBASE/$CLIKEY
set +x

