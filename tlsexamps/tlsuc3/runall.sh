#!/bin/sh
#
set -x
go build
# AMQ detault, not TLS enabled
./tlsuc3 -cliCertFile=/ad3/gma/sslwork/2013/client.crt \
	-cliKeyFile=/ad3/gma/sslwork/2013/client.key
# AMQ TLS
STOMP_PORT=61612 ./tlsuc3 -cliCertFile=/ad3/gma/sslwork/2013/client.crt \
	-cliKeyFile=/ad3/gma/sslwork/2013/client.key
# Apollo TLS
STOMP_PORT=62614 ./tlsuc3 -cliCertFile=/ad3/gma/sslwork/2013/client.crt \
	-cliKeyFile=/ad3/gma/sslwork/2013/client.key
set +x

