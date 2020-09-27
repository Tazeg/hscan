#!/bin/bash

# This script is just to build and zip releases

go test

# install on OS
go build
go install

# build Linux
go build
7z a -t7z hscan-linux64.7z hscan
rm hscan

# build Windows
env GOOS=windows GOARCH=amd64 go build -o hscan.exe hscan.go
7z a -t7z hscan-win64.7z hscan.exe
rm hscan.exe

# build Raspberry Pi
env GOARM=7 GOARCH=arm go build hscan.go
7z a -t7z hscan-arm64.7z hscan
rm hscan
