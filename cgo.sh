#!/bin/bash -e

go build -o cgo/build/client.so -buildmode=c-shared cgo/client/client.go
