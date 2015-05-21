#!/bin/bash
#
set -x
. ../shbasic
export STOMP_NQS=10
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
STOMP_DRAIN=yes STOMP_VHOST=$vh STOMP_PORT=$port STOMP_PROTOCOL=$proto go run getters.go

