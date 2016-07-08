# stompngo_examples - Details #

## Introduction  ##

This repository is a set of STOMP client examples using go.

These examples use the stomp client package here:

[stompngo STOMP client library](https://github.com/gmallard/stompngo)

The reader is urged to become familiar with the go documentaion
for the *stompngo* package:

[stompngo documentation](http://godoc.org/github.com/gmallard/stompngo)<br />
[stompngo wiki](https://github.com/gmallard/stompngo/wiki)

## List of Individual Examples  ##

A brief explanation of the individual examples follows. The list of consistes
of example go, Java, and properties files:

<table border="1" style="width:80%;border: 1px solid black;">
<tr>
<th style="width:20%;border: 1px solid black;padding-left: 10px;" >
Example Name
</th>
<th style="width:60%border: 1px solid black;padding-left: 10px;" >
Explanation
</th>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
ack/ack.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Example of using ACK to acknowledge received messages.  The STOMP
destination should have been previously loaded with message(s).
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
conndisc/conndisc.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A very basic demonstration of a 'CONNECT' / 'DISCONNECT' sequence.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
conndisc_tls/conndisc_tls.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A very basic demonstration of a 'CONNECT' / 'DISCONNECT' sequence with
ssl (tls).  You must connect to a broker port that is 'ssl/tls' enabled.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A subdirectory which contains go and Java examples of a producer and consumer.<br />
This demonstrates interoperability between go and Java, and between STOMP and JMS.<br />
See individual files for details.<br />
This example is ActiveMQ specific, but can be easily adapted for other brokers.<br />
It is assumed that the reader is familiar with Java, JMS, and JNDI.<br />
A number of helper shell scripts are provided.  See the script descriptions below.<br />
These interoperability examples have some hard coded port numbers and path names.<br />
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/Constants.java
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Static constants used in the other Java programs in this directory.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/jndi.properties
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The JNDI properties file definition.
</td>
</tr>


<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/log4j.properties
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The log4j properties file used by the Java code.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/gorecv/gorecv.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message consumer written in go.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/gosend/gosend.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message producer written in go.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/Receiver.java
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message consumer written in Java using JMS.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/Sender.java
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message producer written in Java using JMS.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
publish/publish.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A classic 'putter'.  Used to send an arbitrary number of messages to
a given destination.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
receipts/onack/onack.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Demonstrate obtaining a RECEIPT for an ACK request.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
receipts/onsend/onsend.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Demonstrate obtaining a RECEIPT for a SEND request.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
recv_mds/recv_mds.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An example intended to demonstrate how different brokers distribute output
messages when a client subscribes multiple times to the same destination.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
sngecomm
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A subdirectory which defines helper code for these examples.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
sngecomm/environment.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Handle overrides from the environment for these examples.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
sngecomm/utilities.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Utility routines used by these examples.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_1conn/srmgor_1conn.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send and receive, many go routines, one *stompngo.Connection.<br />
One sender go routine per destination.<br />
One receiver go routine per destination.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_1smrconn/srmgor_1smrconn.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send and receive, many go routines.<br />
One sender connection, with one go routine per destination.<br />
Many receiver connections: one per destination.<br />
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_2conn/srmgor_2conn.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_manyconn/srmgor_manyconn.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
subscribe/subscribe.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A basic demonstration of subscribing and receiving messages.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc1/tlsuc1.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc2/tlsuc2.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc3/tlsuc3.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc4/tlsuc4.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

</table>

## Shell Scripts  ##

This is a list of script files used in the examples.  There are very short
scripts and it should be trivial to convert them to Windows .bat / .cmd
files if necessary.

<table border="1" style="width:80%;border: 1px solid black;">
<tr>
<th style="width:20%;border: 1px solid black;padding-left: 10px;" >
Script Name
</th>
<th style="width:60%border: 1px solid black;padding-left: 10px;" >
Description
</th>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/clean.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This is the clean up script for the Java / JMS interoperability examples.<br />
It removes all three .class files and the two go executables.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/compile.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This is the compile script for the Java / JMS interoperability examples.<br />
It compiles the Java and the go interoperability code.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/cp.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An internal script, sourced by most of the other scripts.<br />
It builds a list of jar files that will be included in the Java CLASSPATH.<br />
This script should be modified to support your environment.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/gorecv.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This script runs the go receiver/consumer/getter.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/gosend.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This script runs the go sender/producer/putter.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/jrecv.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This script runs the Java receiver/consumer/getter.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/jsend.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This script runs the Java sender/producer/putter.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_1conn/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send and receive, many go routines, one *stompngo.Connection.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_1smrconn/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_2conn/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_manyconn/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc1/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc2/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc3/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc4/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

</table>

## Variable Names in the Examples  ##

Note the author is accustomed to idiomatic go variable names (short, 1
character if possible).

<table border="1" style="width:80%;border: 1px solid black;">
<tr>
<th style="width:20%;border: 1px solid black;padding-left: 10px;" >
Variable Name
</th>
<th style="width:60%border: 1px solid black;padding-left: 10px;" >
Common Use
</th>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
ah
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Ack headers.  A stompngo.Headers instance.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
c
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A loop couner.  Because in some situations it makes sense. And I enjoy
writing 'c++' for the end of loop condition.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
ch
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Connect headers.  A stompngo.Headers instance, used for the initial
CONNECT frame sent to the broker.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
conn
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A instance of a *stompngo.Connection.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
d
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A stomp destination.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
dh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Disconnect headers.  A stompngo.Headers instance.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
dt
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A time.Duration instance.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
e
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A go 'error' instance.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
h
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A broker's (DNS) host name.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
i
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A loop counter.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
id
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A (type 4) UUID.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jhp
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A joined host and port pair, returned from net.JoinHostPort(h, p).
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
ll
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A go logger instance.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
mc
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message count.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
mcs
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message count, with type string.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
md
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A stompngo.MessageData instance.  The instance is retrieved from r (i.e.
md := &lt;-r)
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
ms
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message body, with type string.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
n
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An instance of net.Conn.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
nmsgs
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The number of messages to process (produce/consume).
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
nqs
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The number of destinations to use.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
nr
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The number of receiver go routines.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
p
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A broker's listener port.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
pbc
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The print byte count.  When message bodies are printed, this is the
maximum number of bytes to print.  Useful when message body sizes are
large.  This is arbitrarily set to 64 unless overridden.  If this is
set to 0, message bodies are (usually) not printed at all.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
qn
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A queue number identifier.  Used in a looping control structure for a
variable number of queues.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
qns
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A queue number identifier, type string (from qn).
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
r
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A channel of type stompngo.MessageData.  Used in example code to
receive messages and metadata from the broker.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
rd
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A receipt message from the broker. An instance of type stompngo.MessageData.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
rf
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A receive wait time multiplier.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
rid
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A receipt id.
</td>
</tr>


<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
rw
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A wait flag used by receivers. When set to true, receivers wait for a
random amount of time after each message red.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
sbh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Subscribe headers.  A stompngo.Headers instance.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
sc
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A channel of type stompngo.MessageData.  Used in example code to
receive messages and metadata from the broker.
</td>
</tr>


<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
sf
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A send wait time multiplier.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send headers.  A stompngo.Headers instance.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
sw
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A wait flag used by senders. When set to true, senders wait for a
random amount of time after each send.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tc
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An instance of *tls.Config.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
td
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A time.Duration.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tmr
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A time.Timer.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
wga
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Wait group for all sender and receiver go routines.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
wgr
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Wait group for all receiver go routines.
</td>
</tr>


<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
wgs
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Wait group for all sender go routines.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
wh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Work headers.  A stompngo.Headers instance.
</td>
</tr>

</table>

