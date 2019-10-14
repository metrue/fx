#!/usr/bin/env bash

set -e

fx="./build/fx"
service='fx-service-abc'

run() {
  local lang=$1
  local port=$2
  $fx up --name ${service}_${lang} --port ${port} --healthcheck test/functions/func.${lang}
  $fx list # | jq ''
  $fx down ${service}_${lang} # | grep "Down Service ${service}"
}

deploy() {
  local lang=$1
  local port=$2
  if [[ -z "$DOCKER_USERNAME" || -z "$DOCKER_PASSWORD" ]];then
    echo "skip deploy test since no DOCKER_USERNAME and DOCKER_PASSWORD set"
  else
    $fx deploy --name ${service}-${lang} --port ${port} test/functions/func.${lang}
    docker ps
    $fx destroy ${service}-${lang}
  fi
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
docker stop fx-agent || true && docker rm fx-agent || true

$fx infra activate localhost
port=20000
for lang in 'js' 'rb' 'py' 'go' 'php' 'java' 'd'; do
  run $lang $port
  ((port++))

  deploy $lang $port

  build_image $lang "test-fx-image-build-${lang}"
  mkdir -p /tmp/${lang}/images
  export_image ${lang} /tmp/${lang}/images
  rm -rf /tmp/${lang}/images
done

wait
