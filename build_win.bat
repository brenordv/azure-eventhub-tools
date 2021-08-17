echo off
echo Setting OS and Architecture for windows / amd64...
set GOOS=windows
set GOARCH=amd64

echo Building: HUBSEND
go build -o hubsend.exe ./cmd/hubsend

echo Building: HUBREAD
go build -o hubread.exe ./cmd/hubread

echo Building: HUBEXPORT
go build -o hubexport.exe ./cmd/hubexport


7z a -tzip .\.dist\az-eventhub-tools--windows-amd64--%*.zip *.exe readme.md