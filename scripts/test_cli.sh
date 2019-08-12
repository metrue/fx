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

# main
$fx infra activate localhost
port=20000
for lang in 'js' 'rb' 'py' 'go' 'php' 'jl' 'java' 'd'; do
  run $lang $port
  ((port++))
done

wait
