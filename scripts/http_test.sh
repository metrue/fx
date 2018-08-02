#!/usr/bin/env bash

set -e

# start fx server
./build/fx serve > server_output 2>&1 &
sleep 20 # waiting fx server to pulling resource done

base_url='http://127.0.0.1:30080/v1'

# up
echo '{"Functions":[{"Content": "module.exports = (input) => {return parseInt(input.a, 10) + parseInt(input.b, 10)}", "Lang": "node"}]}' | http post ${base_url}/up

# call
http post ${base_url}/call Content='module.exports = (input) => {return parseInt(input.a, 10) + parseInt(input.b, 10)}' Params='a=1 b=2' Lang=node | jq '.Data'

# list
http post http://127.0.0.1:30080/v1/list | jq '.Instances'

# done
echo '{"ID":["*"]}' | http post ${base_url}/down | jq '.Instances'
