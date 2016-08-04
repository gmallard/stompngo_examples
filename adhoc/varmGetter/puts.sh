#!/usr/bin/env bash
# ------------------------------------------------------------------------------
# Entry condition: all queues to be used start off empty.
# Use this script to 'prime' queues for further experiments.
# ------------------------------------------------------------------------------
cmd_base=$(dirname $0)
source $cmd_base/funcs.sh
# ------------------------------------------------------------------------------
init
# ------------------------------------------------------------------------------
showparms
# ------------------------------------------------------------------------------
cqn=1
while [ "$cqn" -le "${MAX_QUEUE}" ]; do
	echo "---------------------------------------------------------------------"
	nqn="${QNBASE}${cqn}"
	ncmd="STOMP_DEST=${nqn} STOMP_NMSGS=${MSG_COUNT} go run $PUTTER"
	echo "Next Puts: ${ncmd}"
	eval ${ncmd}
	let cqn=cqn+1
done

