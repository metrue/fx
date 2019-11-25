#!/usr/bin/env bash

## install docker
sudo apt-get remove -y docker docker-engine docker.io containerd runc
apt-get update -y
sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common lsb-core curl
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update -y
sudo apt-get install -y docker-ce

docker run hello-world
