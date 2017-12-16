#!/usr/bin/env bash

HERE=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ROOT=$( cd ${HERE}/.. && pwd )

dest_dir="${ROOT}/tmp/protoc"

has() {
  type "$1" > /dev/null 2>&1
}

install_protoc() {
    local url=$1
    mkdir -p ${dest_dir}
    wget ${url} -O ${dest_dir}/protoc.zip
    cd ${dest_dir} && unzip ./protoc.zip && chmod +x ./bin/protoc
    ./bin/protoc --version
}

if has "${dest_dir}/bin/protoc"; then
    protoc --version
else
    if [ "$(uname)" == "Darwin" ]; then
        url="https://github.com/google/protobuf/releases/download/v3.5.0/protoc-3.5.0-osx-x86_64.zip"
        install_protoc ${url}
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        url="https://github.com/google/protobuf/releases/download/v3.5.0/protoc-3.5.0-linux-x86_64.zip"
        install_protoc ${url}
    else
        echo 'Sorry, this script works on Linux/Mac now. Please, Create a PR to provide Windows support!'
    fi
fi
