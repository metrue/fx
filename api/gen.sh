#!/usr/bin/env bash

PROTOSRC="./fx.proto"

# generate the gRPC code
protoc  -I/usr/local/include -I. --go_out=plugins=grpc:. $PROTOSRC

# generate the JSON interface code
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  $PROTOSRC

# generate the reverse proxy
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:. \
  $PROTOSRC

# generate the swagger definitions
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --swagger_out=logtostderr=true:../swagger \
  $PROTOSRC
