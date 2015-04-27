#!/bin/bash
#
set -x
. ../shbasic
# Number of messages to send
npms=100
# Number of messages to consume
ncms=10
#
proto=1.1
#
port=61613 # ActiveMQ Here
#port=62613 # Apollo Here
#port=41613 # RabbitMQ Here
#
vh=localhost # OK for AMQ or Apollo
#vh=/ # OK for RMQ
#
STOMP_NMSGS=$npms STOMP_VHOST=$vh STOMP_PORT=$port STOMP_PROTOCOL=$proto STOMP_DEST=/queue/client01-hang.1 go run $basic/publish/publish.go
STOMP_NMSGS=$ncms STOMP_VHOST=$vh STOMP_PORT=$port STOMP_PROTOCOL=$proto STOMP_DEST=/queue/client01-hang.1 go run client02-hang.go
set +x

