#!/bin/bash
set -e
# ++
# verified on Ubuntu 16.04 x64
# ++
user_host=$1

ssh ${user_host} 'bash -s' <<EOF
  apt-get remove -y docker docker-engine docker.io containerd runc
  apt-get update -y
  apt-get install -y apt-transport-https ca-certificates curl software-properties-common lsb-core
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
  add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu \$(lsb_release -cs) stable"
  apt-get update -y
  apt-get install -y docker-ce
  docker run hello-world
EOF
