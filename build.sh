#!/bin/sh

DIR=$PWD
export GOOS=linux
export GOARCH=amd64

echo "Download dependencies"
go get gopkg.in/yaml.v2

./build-plugins.sh

cd $DIR || exit 1
echo "Building Dasher"
go build -ldflags="-s -w" src/dasher.go
upx --brute dasher > /dev/null