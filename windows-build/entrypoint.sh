#!/bin/bash
export CGO_ENABLED=1
export GOOS=windows
GOARCH=amd64 go build -o dist/ktanemod-remote-math-interface-amd64.dll -buildmode=c-shared .
GOARCH=386 go build -o dist/ktanemod-remote-math-interface-386.dll -buildmode=c-shared .
