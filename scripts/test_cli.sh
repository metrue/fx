#!/usr/bin/env bash

set -e

fx="./build/fx"
service='fx-service-abc'

run() {
  local lang=$1
  $fx up --name ${service}_${lang} examples/functions/func.${lang}
  $fx list # | jq ''
  $fx down ${service}_${lang} # | grep "Down Service ${service}"
}

# main
$fx init

for lang in 'js' 'rb' 'py' 'go' 'php' 'jl' 'java' 'd' 'rs'; do
  run $lang &
done

wait
