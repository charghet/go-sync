@echo off

go build -o .\build\go-sync.exe .\cmd\go-sync\go-sync.go
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o .\build\go-sync .\cmd\go-sync\go-sync.go