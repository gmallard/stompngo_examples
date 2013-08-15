
Cautionary note:

These TLS examples have been successfully tested with Apache Apollo, using a
snapshot build:

apache-apollo-99-trunk-20130803.032858-250-unix-distro.tar.gz

Use cases 3 and 4 fail when being tested against an ActiveMQ 5.7.0 or 5.8.0 
instance that is configured to require client certificates.

The exact cause of this failure is unknown at present.  SSL logging seems to 
indicate that AMQ does not correctly read the client certificate chain.

#-------------------------------------------------------------------------------

Update: Use cases 3 and 4 pass running against ActiveMQ.

Code changes: none

Modifications made:  correctly build ActiveMQ's trust store!

