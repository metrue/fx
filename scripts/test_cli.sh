#!/usr/bin/env bash

set -e

service='fx-service-abc'

run() {
  local lang=$1
  ./build/fx up --name ${service}_${lang} examples/functions/func.${lang}
  ./build/fx list # | jq ''
  ./build/fx down ${service}_${lang} # | grep "Down Service ${service}"
}

./build/fx init
for lang in 'js' 'rb' 'py' 'go' 'php' 'jl' 'java' 'd' 'rs'; do
  run $lang &
done

wait
