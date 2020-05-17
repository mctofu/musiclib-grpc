#!/bin/bash -e

go build -o cgo/build/client.o -buildmode=c-shared cgo/client/client.go