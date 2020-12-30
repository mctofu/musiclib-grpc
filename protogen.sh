#!/bin/bash -e
# Regenerates clients & server stubs from the grpc proto def
docker-compose run --rm protoc /protobuf/bin/protoc -I/protobuf --go_out=/mlibgrpc --go_opt=module=github.com/mctofu/musiclib-grpc --go-grpc_out=/mlibgrpc --go-grpc_opt=module=github.com/mctofu/musiclib-grpc --proto_path=/mlibgrpc /mlibgrpc/musiclib.proto
