# stompngo_examples - Java JMS Interoperability #

## Assumptions ##

* ActiveMQ Message Broker
* Openwire port at 61616
* STOMP port at 61613

## Java and go Preparation ##

* Modify cp.sh to provide your AMQ lib location
* Run ./compile.sh

## Send With go ##

* Run ./gosend.sh

## Receive With go ##

* Run ./gorecv.sh

## Send With Java ##

* Run ./jsend.sh

## Receive With Java ##

* Run ./jrecv.sh

## Cleanup After Testing ##

* Run ./clean.sh

## Miscellaneous ##

One can send and receive using Java of course, or send and receive using go.

