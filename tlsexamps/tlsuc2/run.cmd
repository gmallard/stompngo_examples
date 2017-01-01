@echo off
::
IF [%CERTBASE%]==[] SET CERTBASE=%USERPROFILE%\Documents\tls-certs
IF [%CACERT%]==[] SET CACERT=ca.crt
@echo on
go build
tlsuc2.exe -srvCAFile=%CERTBASE%\%CACERT%
@echo off
