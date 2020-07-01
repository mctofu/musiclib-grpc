#!/bin/bash -e
# Builds the cgo c client
go build -o cgo/build/client.so -buildmode=c-shared cgo/client/client.go
