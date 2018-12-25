#!/usr/bin/env bash

HERE=$( cd "$( dirname "${0}" )" && pwd )
ROOT=$( cd ${HERE}/.. && pwd )
PROTOSRC="./fx.proto"

protoc_bin="${ROOT}/third_party/protoc/bin/protoc"
protoc_include="${ROOT}/third_party/protoc/include"

VENDOR="${ROOT}/vendor"

# generate the gRPC code
${protoc_bin} -I/usr/local/include \
    -I${protoc_include} \
    -I. \
    --go_out=plugins=grpc:. \
    $PROTOSRC

# generate the JSON interface code
${protoc_bin} -I/usr/local/include \
    -I${protoc_include} \
    -I. \
    -I$GOPATH/src \
    -I${VENDOR}/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --go_out=plugins=grpc:. \
    $PROTOSRC

# generate the reverse proxy
${protoc_bin} -I/usr/local/include \
    -I${protoc_include} \
    -I. \
    -I$GOPATH/src \
    -I${VENDOR}/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc-gateway_out=logtostderr=true:. \
    $PROTOSRC

# generate the swagger definitions
${protoc_bin} -I/usr/local/include \
    -I. \
    -I${protoc_include} \
    -I$GOPATH/src \
    -I${VENDOR}/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --swagger_out=logtostderr=true:../swagger \
    $PROTOSRC
