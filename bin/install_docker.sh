#!/bin/bash

install_docker() {
  local docker_installer="https://get.docker.com"
  local docker_binary_src="https://download.docker.com/linux/static/stable/x86_64/docker-18.06.3-ce.tgz"
  which docker >/dev/null 2>&1
  if [ $? -eq 0 ]; then
    return
  fi
  curl -fsSL  | sh >/dev/null 2>&1
  if [ ! $? -eq 0 ]; then
    curl -fsSL "$docker_binarys_src" -o docker.tgz
    tar zxvf docker.tgz >/dev/null 2>&1
    sudo mv docker/* /usr/bin && rm -rf docker docker.tgz
  fi
}

install_docker_compose() {
  local docker_compose_src="https://github.com/docker/compose/releases/download/1.24.0/docker-compose-$(uname -s)-$(uname -m)"
  which docker-compose >/dev/null 2>&1
  if [ $? -eq 0 ]; then
    return
  fi
  sudo curl -L "$docker_compose_src" -o /usr/local/bin/docker-compose >/dev/null 2>&1
  sudo chmod +x /usr/local/bin/docker-compose
}

launch_docker_daemon() {
  docker ps >/dev/null 2>&1
  if [ $? -eq 0 ]; then
    return
  fi
  if [ ! -S /var/run/docker.sock ]; then
    (sudo dockerd >/dev/null 2>&1 &)
    local next_wait_time=0
    until [ -S /var/run/docker.sock ] || [ $next_wait_time -eq $1 ]; do
      sleep $(( next_wait_time++ ))
    done
  fi
  if [ ! -S /var/run/docker.sock ]; then
    echo "timeout to wait docker daemon ready"
    exit 1
  fi
  sudo setfacl -m user:$(whoami):rw /var/run/docker.sock
  docker ps >/dev/null 2>&1
}

main() {
  install_docker
  install_docker_compose
  launch_docker_daemon 10
}

main
