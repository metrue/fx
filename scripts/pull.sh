#!/usr/bin/env bash

set -e

# enable docker remote api on host
# docker run -d -v /var/run/docker.sock:/var/run/docker.sock -p 127.0.0.1:1234:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock

# pull basic docker image and try build
ROOT=`pwd`
echo $ROOT
images='api/images'
for lang in `ls ${images}`; do
  cd ${ROOT}/${images}/${lang} && docker build .
done
