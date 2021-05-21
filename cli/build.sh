GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o ./bin/toast64.exe ./*.go
GOOS=windows GOARCH=386 go build -ldflags -H=windowsgui -o ./bin/toast32.exe ./*.go
