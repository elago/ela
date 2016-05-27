#!/bin/bash

go get golang.org/x/tools/cmd/cover
go get github.com/mattn/goveralls
goversion=`go version | grep 1.6`
if [ ${#goversion} !== "0" ];then
	go test
else
	go test -v -covermode=count -coverprofile=coverage.out
	$HOME/gopath/bin/goveralls -coverprofile=coverage.out -service=travis-ci -repotoken $COVERALLS_TOKEN
fi