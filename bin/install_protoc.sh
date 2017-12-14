#!/usr/bin/env bash

targetDir="../vendor/protoc"

has() {
  type "$1" > /dev/null 2>&1
}

install_protoc() {
    local url=$1
    mkdir -p ${targetDir}
    wget ${url} -O ${targetDir}/protoc.zip
    cd ${targetDir} && unzip ./protoc.zip && chmod +x ./bin/protoc
    ./bin/protoc --version
}

if has "${targetDir}/bin/protoc"; then
    protoc --version
else
    if [ "$(uname)" == "Darwin" ]; then
        url="https://github.com/google/protobuf/releases/download/v3.5.0/protoc-3.5.0-osx-x86_64.zip"
        install_protoc ${url}
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        url="https://github.com/google/protobuf/releases/download/v3.5.0/protoc-3.5.0-linux-x86_64.zip"
        install_protoc ${url}
    else
        echo 'Sorry, this script on works on Linux/Mac now'
    fi
fi
