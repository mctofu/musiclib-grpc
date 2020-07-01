#!/bin/bash -e

docker-compose run --rm protoc /protobuf/bin/protoc -I/protobuf --go_out=plugins=grpc:/mlibgrpc/go/mlibgrpc --proto_path=/mlibgrpc /mlibgrpc/musiclib.proto
