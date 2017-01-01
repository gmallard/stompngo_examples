@echo off
::
IF [%CERTBASE%]==[] SET CERTBASE=%USERPROFILE%\Documents\tls-certs
IF [%CLICERT%]==[] SET CLICERT=client.crt
IF [%CLIKEY%]==[] SET CLIKEY=client.key
IF [%CACERT%]==[] SET CACERT=ca.crt
@echo on
go build
tlsuc4.exe -cliCertFile=%CERTBASE%\%CLICERT% -cliKeyFile=%CERTBASE%\%CLIKEY% -srvCAFile=%CERTBASE%\%CACERT%
@echo off
