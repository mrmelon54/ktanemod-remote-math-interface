#!/bin/bash
export CGO_ENABLED=1
export GOOS=windows
GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -buildvcs=false -o dist/ktanemod-remote-math-interface-windows-amd64.dll -buildmode=c-shared .
GOARCH=386 CC=i686-w64-mingw32-gcc go build -buildvcs=false -o dist/ktanemod-remote-math-interface-windows-386.dll -buildmode=c-shared .
