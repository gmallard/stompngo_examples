#
# Makefile for stompngo_examples
#

dirs = 	ack \
	conndisc \
	conndisc_tls \
	jinterop/gorecv \
	jinterop/gosend \
	publish \
	receipts/onack \
	receipts/onsend \
	recv_mds \
	srmgor_1conn \
	srmgor_1smrconn \
	srmgor_2conn \
	srmgor_manyconn \
	subscribe \
	tlsexamps/tlsuc1 \
	tlsexamps/tlsuc2 \
	tlsexamps/tlsuc3 \
	tlsexamps/tlsuc4

.PHONY: $(dirs) packages clean format

all: $(dirs)
	@for i in $(dirs); do \
	echo $$i; \
	curd=`pwd`; \
	cd $$i && go build; \
	cd $$curd; \
	done

clean:
	@for i in $(dirs); do \
	echo $$i; \
	curd=`pwd`; \
	cd $$i && go clean; \
	cd $$curd; \
	done

format:
	@for i in $(dirs); do \
	echo $$i; \
	curd=`pwd`; \
	cd $$i && gofmt -w *.go; \
	cd $$curd; \
	done

