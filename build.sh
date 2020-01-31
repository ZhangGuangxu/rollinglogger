#!/bin/sh

gofmt -s -l -w *.go && go vet *.go && go build -i

exit 0
