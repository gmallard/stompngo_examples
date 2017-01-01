@echo off
::
IF [%CERTBASE%]==[] SET CERTBASE=%USERPROFILE%\Documents\tls-certs
IF [%CLICERT%]==[] SET CLICERT=client.crt
IF [%CLIKEY%]==[] SET CLIKEY=client.key
@echo on
go build
tlsuc3.exe -cliCertFile=%CERTBASE%\%CLICERT% -cliKeyFile=%CERTBASE%\%CLIKEY%
@echo off

