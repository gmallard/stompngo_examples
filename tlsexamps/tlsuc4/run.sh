#!/bin/sh
#
set -x
go build
./tlsuc4 -srvCAFile=/ad3/gma/sslwork/2016-02/ca.crt \
	-cliCertFile=/ad3/gma/sslwork/2016-02/client.crt \
	-cliKeyFile=/ad3/gma/sslwork/2016-02/client.key
set +x

