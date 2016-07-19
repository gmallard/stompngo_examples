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
Examples are provided for ActiveMQ and Artemis.  Only the ActiveMQ artifacts are described here.<br />
It is assumed that the reader is familiar with Java, JMS, and JNDI.<br />
A number of helper shell scripts are provided.  See the script descriptions below.<br />
These interoperability examples have some hard coded port numbers and path names.<br />
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/Constants.java
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Static constants used in the other Java programs in this directory.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/jndi.properties
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The JNDI properties file definition.
</td>
</tr>


<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/log4j.properties
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The log4j properties file used by the Java code.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/gorecv/gorecv.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message consumer written in go.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/gosend/gosend.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message producer written in go.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/Receiver.java
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message consumer written in Java using JMS.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/Sender.java
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
Many receiver connections:<br />
one receiver connection per destination (one go routine per connection).<br />
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_2conn/srmgor_2conn.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send and receive, two connections.<br />
Many go routines in the sender connection.<br />
Many go routines in the receiver connection.<br />
The receiver go routines illustrate a channel buffering technique that
can be used with in bound message data.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_manyconn/srmgor_manyconn.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send and receive.<br />
Many sender connections.<br />
Many receiver connections.
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
tlsexamps
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A subdirectory demonstrating the use of TLS connections.<br />
All four primary TLS use cases are demonstrated.<br />
UseCase1) Client does not authenticate the broker and broker does not
authenticate the client.<br />
UseCase2) Client does authenticate the broker but broker does not
authenticate the client.<br />
UseCase3) Client does not authenticate the broker but broker does
authenticate the client.<br />
UseCase4) Client does authenticate the broker and broker does
authenticate the client.<br />
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc1/tlsuc1.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
UseCase1) Client does not authenticate the broker and broker does not
authenticate the client.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc2/tlsuc2.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
UseCase2) Client does authenticate the broker but broker does not
authenticate the client.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc3/tlsuc3.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
UseCase3) Client does not authenticate the broker but broker does
authenticate the client.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc4/tlsuc4.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
UseCase4) Client does authenticate the broker and broker does
authenticate the client.
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
jinterop/activemq/clean.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This is the clean up script for the Java / JMS interoperability examples.<br />
It removes all three .class files and the two go executables.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/compile.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This is the compile script for the Java / JMS interoperability examples.<br />
It compiles the Java and the go interoperability code.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/cp.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An internal script, sourced by some of the other scripts.<br />
It builds a list of jar files that will be included in the Java CLASSPATH.<br />
This script should be modified to support your environment.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/gorecv.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This script runs the go receiver/consumer/getter.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/gosend.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This script runs the go sender/producer/putter.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/jrecv.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This script runs the Java receiver/consumer/getter.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/jsend.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This script runs the Java sender/producer/putter.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/activemq/runall.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This script runs all possible combinations of producers / consumers.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_1conn/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send and receive, many go routines, one *stompngo.Connection.<br />
A logging helper script.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_1smrconn/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send and receive, one sender connection with many go routines,
many receiver connections.<br />
A logging helper script.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_2conn/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send and receive, one receiver connection with many go routines, one
sender connection with many go routines.<br />
A logging helper script.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_manyconn/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send and receive, many receiver connections, many sender connections.<br />
A logging helper script.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc1/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TLS Use Case 1.<br />
Example run script.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc2/runall.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TLS Use Case 2.<br />
Example run script.<br />
Three separate runs, using different ports.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc2/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TLS Use Case 2.<br />
Example run script.
</td>
</tr>


<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc3/runall.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TLS Use Case 3.<br />
Example run script.<br />
Three separate runs, using different ports.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc3/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TLS Use Case 3.<br />
Example run script.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc4/runall.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TLS Use Case 4.<br />
Example run script.<br />
Three separate runs, using different ports.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tlsexamps/tlsuc4/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TLS Use Case 4.<br />
Example run script.
</td>
</tr>

</table>

## Variable Names in the Examples  ##

Note the author is accustomed to idiomatic go variable names (short if possible).

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
A loop couner.  Because in some situations it makes sense. And because
writing 'c++' for the end of loop action is amusing.
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
A instance of a *stompngo.Connection, usually obtained from a call to
stompngo.Connect(...).<br />
Rename candidate: sc.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
d
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A stomp destination, type string.  Format and semantics are broker specific.
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
An instance of the go error type.
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
hap
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A host and port name returned from net.JoinHostPort(h, p), type string.<br />
Example: localhost:61613.
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
A (type 4) UUID obtained from calling stompngo.Uuid(), type string.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
irid
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A receipt id, received from the broker, type string.
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
mse
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A message body, with addional data, type string.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
n
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An implementation of the net.Conn interface, obtained by calling
net.Dial.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
nmsgs
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The number of messages to process (produce/consume).<br />
Rename candidate: nm.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
nqs
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The number of destinations to use.<br />
Rename candidate: nd.
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
set to 0, message bodies are (usually) not printed at all.<br />
Rename candidate: pc.<br />
TODO: consistently print no message body if this is &lt;= 0.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
qc
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A channel of type bool, used for signaling event completion.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
qn
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A queue number identifier.  Used in a looping control structure for a
variable number of queues, or as a method parameter.
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
Used in the sngecomm subpackage only, for several reasons.<br />
TODO:  this needs to change to a more meaningful variable name and
documentation for same added.
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
A receipt id, type string.
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
si
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A loop counter formatted as a string.
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

## Variable Names Specific to the TLS Examples  ##

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
b
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A buffer, type []byte.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
c
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An instance of *x509.Certificate.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
cc
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An instance of tls.Certificate.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
cliCertFile
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The full file name of the client's public cert in PEM format.
</td>
</tr>


<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
cliKeyFile
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The full file name of the client's private key in PEM format.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
e
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An instance of error.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
k
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An instance of *pem.Block.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
ll
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A logger instance.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
nc
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A pointer to a tls configuration, type *tls.Config, obtained from
'tls.Client'.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
svrCAFile
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The full file name of the server CA's public cert in PEM format.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
t
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
An implementation of the net.Conn interface, obtained by calling
net.Dial.
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
tc
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A pointer to a tls configuration, type *tls.Config, obtained from 'new'.
</td>
</tr>

</table>

