#!/usr/bin/env bash

set -e

fx="./build/fx"
service='fx-service-abc'

run() {
  local lang=$1
  $fx up --name ${service}_${lang} test/functions/func.${lang}
  $fx list # | jq ''
  $fx down ${service}_${lang} # | grep "Down Service ${service}"
}

# main
$fx provision

for lang in 'js' 'rb' 'py' 'go' 'php' 'jl' 'java' 'd'; do
  run $lang &
done

wait
