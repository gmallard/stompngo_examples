#!/usr/bin/env bash
#
# Demo _some_ of the stompngo examples.
#
SNGD=${SNGD:-/home/opt/gosrc/src/github.com/gmallard/stompngo_examples}
echo $SNGD
pushd $SNGD
#
export STOMP_NMSGS=1
export STOMP_NQS=1
#
STOMP_DEST=/queue/snge.ack go run publish/publish.go
STOMP_DEST=/queue/snge.ack.1 go run ack/ack.go
#
STOMP_USEEOF=y STOMP_DEST=/queue/snge.rdr go run publish/publish.go
STOMP_USEEOF=y STOMP_DEST=/queue/snge.rdr.1 go run adhoc/reader/reader.go
#
go run conndisc/conndisc.go
#
# Port 61611: AMQ, TLS, No cert required
STOMP_PORT=61611 go run conndisc_tls/conndisc_tls.go
#
go run srmgor_1conn/srmgor_1conn.go
go run srmgor_1smrconn/srmgor_1smrconn.go
go run srmgor_2conn/srmgor_2conn.go
go run srmgor_manyconn/srmgor_manyconn.go
#
STOMP_USEEOF=y STOMP_DEST=/queue/snge.sub go run publish/publish.go
STOMP_USEEOF=y STOMP_DEST=/queue/snge.sub.1 STOMP_NMSGS=2 go run subscribe/subscribe.go
#
popd
