#!/usr/bin/env bash

set -e

service='fx-service-abc'

for lang in 'js' 'rb' 'py' 'go' 'php' 'jl' 'java' 'd'; do
  ./build/fx up --name ${service}_${lang} examples/functions/func.${lang} # | grep 'info Run Service:'
  # when call the service, we have to make sure input params is correct (include correct type, since some statical language like Golang, it Unmashal payload into specific type)
  # ./build/fx call examples/functions/func.${lang} a=1 b=2 # | grep '3'
  ./build/fx list # | jq ''
  ./build/fx down ${service}_${lang} # | grep "Down Service ${service}"
done
