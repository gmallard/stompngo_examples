#!/bin/bash
#
set -x
. ../shbasic
export STOMP_NMSGS=10
export STOMP_NQS=1 # Looks funny, but is not for this example
#
#
proto=1.2
#
port=61613 # ActiveMQ Here
port=62613 # Apollo Here
port=41613 # RabbitMQ Here
#
vh=localhost # OK for AMQ or Apollo
vh=/ # OK for RMQ
#

for i in 1 2 3 4 5 6 7 8 9 10;do
#for i in 1 2;do
	STOMP_VHOST=$vh STOMP_PORT=$port STOMP_PROTOCOL=$proto STOMP_DEST=/queue/sngissue25.$i go run $basic/publish/publish.go
done
