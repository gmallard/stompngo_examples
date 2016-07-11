#!/bin/sh
#
set -x
go build
./tlsuc4 -srvCAFile=/ad3/gma/sslwork/2013/TestCA.crt \
	-cliCertFile=/ad3/gma/sslwork/2013/client.crt \
	-cliKeyFile=/ad3/gma/sslwork/2013/client.key
set +x

