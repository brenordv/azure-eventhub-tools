echo off

echo Checking required build folders...
if not exist ".tmp" mkdir .tmp
if not exist ".tmp\win64" mkdir .tmp\win64

echo Setting OS and Architecture for windows / amd64...
set GOOS=windows
set GOARCH=amd64


echo Building: HUBSEND...
go build -o .tmp\win64\hubsend.exe .\cmd\hubsend

echo Building: HUBREAD...
go build -o .tmp\win64\hubread.exe .\cmd\hubread

echo Building: HUBEXPORT...
go build -o .tmp\win64\hubexport.exe .\cmd\hubexport


echo Creating archive for version %*...
7z a -tzip .\.dist\az-eventhub-tools--windows-amd64--%*.zip .\.tmp\win64\*.exe readme.md

echo Copying application to current directory...
copy /b/v/y .\.tmp\win64\*.exe .\

echo All done!
