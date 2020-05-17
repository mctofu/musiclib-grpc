#!/bin/bash -e

docker run --rm -it -v `pwd`:/mlibgrpc protoc /protobuf/bin/protoc -I/protobuf --go_out=plugins=grpc:/mlibgrpc/go/mlibgrpc --proto_path=/mlibgrpc /mlibgrpc/musiclib.proto
