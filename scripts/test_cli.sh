#!/usr/bin/env bash

set -e

fx="./build/fx"
service='fx-service'

run() {
  local lang=$1
  local port=$2
  # localhost
  $fx up --name ${service}_${lang} --port ${port} --healthcheck test/functions/func.${lang}
  $fx list
  $fx down ${service}_${lang} || true
}

build_image() {
  local lang=$1
  local tag=$2
  $fx image build -t ${tag} test/functions/func.${lang}
}

export_image() {
  local lang=$1
  local dir=$2
  $fx image export -o ${dir} test/functions/func.${lang}
}

# main
port=20000
for lang in ${1}; do
  run $lang $port
  ((port++))

  build_image $lang "test-fx-image-build-${lang}"
  mkdir -p /tmp/${lang}/images
  export_image ${lang} /tmp/${lang}/images
  rm -rf /tmp/${lang}/images
done

wait
