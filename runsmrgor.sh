#!/bin/sh
#
export STOMP_NQS=10
export STOMP_NMSGS=100
#
export STOMP_PORT=61613
#export STOMP_PORT=62613
#
for d in srmgor_1conn srmgor_1smrconn srmgor_2conn srmgor_manyconn;do
	echo "---------${d}---------"
	cd $d
	go run $d.go
	cd ..
done

