@echo off
cd frontend
cmd /C "pnpm build"
cd ..

go build -o .\build\go-sync.exe .\main.go
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o .\build\go-sync .\main.go