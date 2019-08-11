#!/usr/bin/env bash

## install docker
sudo apt-get remove -y docker docker-engine docker.io containerd runc
apt-get update -y
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common lsb-core curl
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu \$(lsb_release -cs) stable"
suod apt-get update -y
sudo apt-get install -y docker-ce

docker run hello-world

## start fx proxy agent
docker run -d --name=fx-agent --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:8866:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock
