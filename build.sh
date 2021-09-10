#!/bin/sh

DIR=$PWD

export GO111MODULE=auto

echo "Download dependencies"
go get gopkg.in/yaml.v2
go get github.com/yalp/jsonpath

./build-plugins.sh

cd $DIR || exit 1
echo "Building Dasher"
go build -ldflags="-s -w" src/dasher.go
