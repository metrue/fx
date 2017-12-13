#!/usr/bin/env bash

has() {
  type "$1" > /dev/null 2>&1
}

install_protoc_mac() {
    brew tap grpc/grpc
    brew install --with-plugins grpc
    protoc --version
}

install_protoc_linux() {
    tmp_dir=/tmp/protoc$(date '+%s')
    mkdir -p ${tmp_dir}
    wget https://github.com/google/protobuf/releases/download/v3.5.0/protoc-3.5.0-linux-x86_64.zip -O ${tmp_dir}/protoc.zip
    cd ${tmp_dir} && unzip ./protoc.zip && mv ./bin/protoc /usr/local/bin/
    rm -rf ${tmp_dir}
    chmod +x /usr/local/bin/protoc
    protoc --version
}

if has "protoc"; then
    protoc --version
else
    if [ "$(uname)" == "Darwin" ]; then
        install_protoc_mac
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        install_protoc_linux
    else
        echo 'Sorry, this script on works on Linux/Mac now'
    fi
fi
