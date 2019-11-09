#!/usr/bin/env bash

set -e

fx="./build/fx"
service='fx-service'

run() {
  local lang=$1
  local port=$2
  # localhost
  $fx up --name ${service}_${lang} --port ${port} --healthcheck test/functions/func.${lang}
  # remote host
  DOCKER_REMOTE_HOST_ADDR=${REMOTE_HOST_ADDR} \
  DOCKER_REMOTE_HOST_USER=${REMOTE_HOST_USER} \
  DOCKER_REMOTE_HOST_PASSWORD=${REMOTE_HOST_PASSWORD} \
  # force destroy it first
  echo 'stopping ...'
  $fx down ${service}_${lang} || true

  ./build/fx up -p 1234 test/functions/func.js
  $fx list # | jq ''
  # $fx down ${service}_${lang}
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
# clean up
# docker stop fx-agent || true && docker rm fx-agent || true

port=20000
for lang in 'js' 'rb' 'py' 'go' 'php' 'java' 'd'; do
  run $lang $port
  ((port++))

  build_image $lang "test-fx-image-build-${lang}"
  mkdir -p /tmp/${lang}/images
  export_image ${lang} /tmp/${lang}/images
  rm -rf /tmp/${lang}/images
done

wait
