# TLS Examples - Usage Notes

Hopefully the four "use cases" present will suffice for any client scenario
required.  Please read the comments in the code for the details of each use case.

## Using Custom Ciphers

Recently added has been an experiment showing the optional use of a custom
cipher list when connecting to a TLS broker.

## Using ciphers not supported by the go stdlib

It might happen that a client wishes to use a cipher not supported by the
go stdlib runtime.

Experiments with that scenario are not included in these examples, but have been
attempted locally.

The result of those experiments is: this attempt will most likely fail.

A somewhat more lengthy description of these experiments follows.

Conclusions are based on content of Java based broker logs, with ssl
debugging on (-Djavax.net.debug=all).

### Some ciphers are go supported, some are not

Experimental details:

* Using 3 go supported ciphers
* Using 3 ciphers not supported by go
* Broker logs show that the ciphers not supported by go have
  been eliminated from the ClientHello message
* Connection succeds using one of the go supported ciphers

### No ciphers are go supported

Experimental details:

* Using 3 ciphers not supported by go
* Broker logs show that no ciphers are sent to the broker in the ClientHello
  message
* Connection fails with "no common cipher present"

