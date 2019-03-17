# Artemis Java / JMS and _stompngo_ Interoperability

## The Included JAR File

I finally found this in the Artemis documentation:

<blockquote>
Apache ActiveMQ Artemis requires just a single jar on the client classpath.

Warning

The client jar mentioned here can be found in the lib/client directory of the Apache ActiveMQ Artemis distribution. Be sure you only use the jar from the correct version of the release, you must not mix and match versions of jars from different Apache ActiveMQ Artemis versions. Mixing and matching different jar versions may cause subtle errors and failures to occur.
</blockquote>

The currently commited JAR is from Artemis 2.4.0.  YMMV.

## Port Numbers

Port numbers (local) are coded in:

* jndi.properties

This file likely needs to change to fit an indivadual's environment.  It should
be an Artemis port that supports JMS.

The STOMP port should be specified in the environment, i.e.:

export STOMP_PORT=50613

## Destinations / Queue Names

Take care in spcifying queue names in jndi.properties.  This must be carefully
done in order to match what is required for STOMP destination names.
