#!/bin/zsh

TAG=$(git describe --abbrev=0 --tags | head -n1)

set -e
go mod tidy
go build

GOOS=linux GOARCH=amd64 go build -o redis-flush.linux.amd64
GOOS=linux GOARCH=arm64 go build -o redis-flush.linux.arm64
GOOS=darwin GOARCH=amd64 go build -o redis-flush.darwin.amd64
GOOS=darwin GOARCH=arm64 go build -o redis-flush.darwin.arm64
GOOS=windows GOARCH=amd64 go build -o redis-flush.windows.amd64.exe

gh release create -p --generate-notes "${TAG}" redis-flush.*
