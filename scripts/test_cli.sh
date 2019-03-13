#!/usr/bin/env bash

set -e

# start fx server
./build/fx serve > server_output 2>&1 &
sleep 20 # waiting fx server to pulling resource done

service='fx-service-abc'

for lang in 'js' 'rb' 'py' 'go' 'php' 'jl' 'java' 'd'; do
  ./build/fx up --name ${service} examples/functions/func.${lang} | grep 'info Run Service:'
  ./build/fx call examples/functions/func.js a=1 b=2 | grep '3'
  ./build/fx list | jq ''
  ./build/fx down ${service} | grep "Down Service ${service}"
done
