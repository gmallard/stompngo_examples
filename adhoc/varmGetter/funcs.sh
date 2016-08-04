# ------------------------------------------------------------------------------
function init {
	if [ "${PUTTER}X" = "X" ]; then
		PUTTER=${cmd_base}/../../publish/publish.go
	fi
	# --------------------------------------------------------------------------
	# Default queue name prefix
	QNBASE=${QNBASE:-/queue/varmGetter.}
	# Default number of queues to prime
	MAX_QUEUE=${MAX_QUEUE:-9}
	# Default number of messages put to each queue
	MSG_COUNT=${MSG_COUNT:-10}
}
# ------------------------------------------------------------------------------
function showparms {
	echo "Queue Name Base: ${QNBASE}"
	echo "Number of Queues: ${MAX_QUEUE}"
	echo "Message Count: ${MSG_COUNT}"
	echo "Put program: ${PUTTER}"
}
# ------------------------------------------------------------------------------
