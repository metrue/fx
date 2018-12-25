#!/usr/bin/env bash

set -e

# start fx server
./build/fx serve > server_output 2>&1 &
sleep 20 # waiting fx server to pulling resource done

for lang in 'js' 'rb' 'py' 'go' 'php' 'jl' 'java' 'd'; do
  ./build/fx up examples/functions/func.${lang} | jq '.Instances'
  ./build/fx call examples/functions/func.js a=1 b=2 | jq '.Data'
  ./build/fx list | jq '.Instances'
  ./build/fx down '*' | jq '.Instances'
done
