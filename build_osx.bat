echo off

echo Checking required build folders...
if not exist ".tmp" mkdir .tmp
if not exist ".tmp\osx" mkdir .tmp\osx

echo Setting OS and Architecture for OSX / amd64...
set GOOS=darwin
set GOARCH=amd64


echo Building: HUBSEND...
go build -o .tmp\osx\hubsend.exe .\cmd\hubsend

echo Building: HUBREAD...
go build -o .tmp\osx\hubread.exe .\cmd\hubread

echo Building: HUBEXPORT...
go build -o .tmp\osx\hubexport.exe .\cmd\hubexport


echo Creating archive for version %*...
7z a -tzip .\.dist\az-eventhub-tools--darwin-osx-amd64--%*.zip .\.tmp\osx\*.exe readme.md

echo All done!
