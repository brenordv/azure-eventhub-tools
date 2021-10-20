echo off

echo Checking required build folders...
if not exist ".tmp" mkdir .tmp
if not exist ".tmp\linux" mkdir .tmp\linux

echo Setting OS and Architecture for linux / amd64...
set GOOS=linux
set GOARCH=amd64


echo Building: HUBSEND...
go build -o .tmp\linux\hubsend.exe .\cmd\hubsend

echo Building: HUBREAD...
go build -o .tmp\linux\hubread.exe .\cmd\hubread

echo Building: HUBEXPORT...
go build -o .tmp\linux\hubexport.exe .\cmd\hubexport


echo Creating archive for version %*...
7z a -tzip .\.dist\az-eventhub-tools--linux-amd64--%*.zip .\.tmp\linux\*.exe readme.md

echo All done!
