FROM golang:1.15

RUN apt-get update && apt-get install -y unzip wget

ARG PROTOBUF_VERSION=3.14.0

RUN mkdir -p /protobuf && \
  mkdir -p /tools && \
  mkdir -p /mlibgrpc && \
  wget "https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOBUF_VERSION/protoc-$PROTOBUF_VERSION-linux-x86_64.zip" && \
  unzip "protoc-$PROTOBUF_VERSION-linux-x86_64.zip" -d "/protobuf"

WORKDIR /tools

RUN go mod init tools && \
  go get google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0 && \
  go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0.1

WORKDIR /mlibgrpc
