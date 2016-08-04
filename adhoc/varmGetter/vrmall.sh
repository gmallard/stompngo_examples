#!/usr/bin/env bash
# ------------------------------------------------------------------------------
eval $DeBug
# ------------------------------------------------------------------------------
cmd_base=$(dirname $0)
# ------------------------------------------------------------------------------
pushd $cmd_base
# ------------------------------------------------------------------------------
go build varmGetter.go
# ------------------------------------------------------------------------------
export STOMP_ACKMODE=client
# ------------------- Phase 1 --------------------------------------------------
# ------------------------------------------------------------------------------
# One message from queue 1
echo "1111111111111111111111111111111111111111111111111111111111111111111111111"
export STOMP_DEST=/queue/varmGetter.1
export STOMP_NMSGS=1
./varmGetter
# ------------------------------------------------------------------------------
# Two messages from queue 2, skip UNSUBSCRIBE
echo "2222222222222222222222222222222222222222222222222222222222222222222222222"
export STOMP_DEST=/queue/varmGetter.2
export STOMP_NMSGS=2
VMG_NOUNSUB=y ./varmGetter
# ------------------------------------------------------------------------------
# Three messages from queue 3, skip DISCONNECT
echo "3333333333333333333333333333333333333333333333333333333333333333333333333"
export STOMP_DEST=/queue/varmGetter.3
export STOMP_NMSGS=3
VMG_NODISC=y ./varmGetter
# ------------------------------------------------------------------------------
# Four messages from queue 4, skip UNSUBSCRIBE and DISCONNECT
echo "4444444444444444444444444444444444444444444444444444444444444444444444444"
export STOMP_DEST=/queue/varmGetter.4
export STOMP_NMSGS=4
VMG_NODISC=y VMG_NOUNSUB=y ./varmGetter
# ------------------- Phase 2 --------------------------------------------------
export VMG_GETAR=y
# ------------------------------------------------------------------------------
# One message from queue 5
echo "5555555555555555555555555555555555555555555555555555555555555555555555555"
export STOMP_DEST=/queue/varmGetter.5
export STOMP_NMSGS=1
./varmGetter
# ------------------------------------------------------------------------------
# Two messages from queue 6, skip UNSUBSCRIBE
echo "6666666666666666666666666666666666666666666666666666666666666666666666666"
export STOMP_DEST=/queue/varmGetter.6
export STOMP_NMSGS=2
VMG_NOUNSUB=y ./varmGetter
# ------------------------------------------------------------------------------
# Three messages from queue 7, skip DISCONNECT
echo "7777777777777777777777777777777777777777777777777777777777777777777777777"
export STOMP_DEST=/queue/varmGetter.7
export STOMP_NMSGS=3
VMG_NODISC=y ./varmGetter
# ------------------------------------------------------------------------------
# Four messages from queue 8, skip UNSUBSCRIBE and DISCONNECT
echo "8888888888888888888888888888888888888888888888888888888888888888888888888"
export STOMP_DEST=/queue/varmGetter.8
export STOMP_NMSGS=4
VMG_NODISC=y VMG_NOUNSUB=y ./varmGetter
# ------------------------------------------------------------------------------
popd
set +x

