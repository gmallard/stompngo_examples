@echo on
::
:: Demo _some_ of the stompngo examples.
::
IF [%SNGD%]==[] SET SNGD=C:\gosrc\src\github.com\gmallard\stompngo_examples
echo %SNGD%
IF [%STOMP_HOST%]==[] SET STOMP_HOST=192.168.1.200
echo %STOMP_HOST%
::
set STOMP_NMSGS=1
set STOMP_NQS=1
::
cd %SNGD%
::
set STOMP_DEST=/queue/snge.ack
go run publish\publish.go
set STOMP_DEST=/queue/snge.ack.1
go run ack\ack.go
::
set STOMP_USEEOF=y
set STOMP_DEST=/queue/snge.rdr
go run publish\publish.go
set STOMP_USEEOF=y
set STOMP_DEST=/queue/snge.rdr.1
go run adhoc\reader\reader.go
::
go run conndisc\conndisc.go
::
:: Port 61611: AMQ, TLS, No cert required
set STOMP_PORT=61611
go run conndisc_tls\conndisc_tls.go
::
set STOMP_PORT=61613
go run srmgor_1conn\srmgor_1conn.go
go run srmgor_1smrconn\srmgor_1smrconn.go
go run srmgor_2conn\srmgor_2conn.go
go run srmgor_manyconn\srmgor_manyconn.go
::
set STOMP_USEEOF=y
set STOMP_DEST=/queue/snge.sub
go run publish\publish.go
set STOMP_USEEOF=y
set STOMP_DEST=/queue/snge.sub.1
set STOMP_NMSGS=2
go run subscribe\subscribe.go
::
