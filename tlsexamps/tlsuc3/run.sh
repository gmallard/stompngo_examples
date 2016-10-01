#!/bin/sh
#
set -x
go build
./tlsuc3 -cliCertFile=/ad3/gma/sslwork/2016-02/client.crt \
	-cliKeyFile=/ad3/gma/sslwork/2016-02/client.key
set +x

