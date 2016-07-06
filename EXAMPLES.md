# stompngo_examples - Details #

## Introduction  ##

This repository is a set of STOMP client examples using go.

These examples use the stomp client package here:

[stompngo](https://github.com/gmallard/stompngo)

## List of Examples  ##

A brief explanation of these examples follows.

This is a list of example go, Java, and properties files:

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
jinterop/Constants.java
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jndi.properties
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/gorecv/gorecv.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/gosend/gosend.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/Receiver.java
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/Sender.java
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
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
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
receipts/onsend/onsend.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
recv_mds/recv_mds.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
sngecomm/environment.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
sngecomm/utilities.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_1conn/srmgor_1conn.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_1smrconn/srmgor_1smrconn.go
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
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
TODO
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
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/compile.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/cp.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/gorecv.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/gosend.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/jrecv.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
jinterop/jsend.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
srmgor_1conn/run.sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
TODO
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
r
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A channel of type stompngo.MessageData.  Used in example code to
receive messages and metadata from the broker.
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
sh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Send headers.  A stompngo.Headers instance.
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
wh
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Work headers.  A stompngo.Headers instance.
</td>
</tr>

</table>

