#!/bin/sh
#
set -x
go build
# AMQ default, not TLS enabled
./tlsuc4 -srvCAFile=/ad3/gma/sslwork/2013/TestCA.crt \
	-cliCertFile=/ad3/gma/sslwork/2013/client.crt \
	-cliKeyFile=/ad3/gma/sslwork/2013/client.key
# AMQ TLS
STOMP_PORT=61612 ./tlsuc4 -srvCAFile=/ad3/gma/sslwork/2013/TestCA.crt \
	-cliCertFile=/ad3/gma/sslwork/2013/client.crt \
	-cliKeyFile=/ad3/gma/sslwork/2013/client.key
# Apollo TLS
STOMP_PORT=62614 ./tlsuc4 -srvCAFile=/ad3/gma/sslwork/2013/TestCA.crt \
	-cliCertFile=/ad3/gma/sslwork/2013/client.crt \
	-cliKeyFile=/ad3/gma/sslwork/2013/client.key
set +x

