#!/usr/bin/env bash

set -e

git config --global user.email "h.minghe@gmail.com"
git config --global user.name "Minghe Huang"

branch=$(git rev-parse --abbrev-ref HEAD)
commit=$(git rev-parse --short HEAD)
version=${branch}-${commit}
if [[ ${branch} == "production" ]];then
  version=$(cat fx.go| grep Version | awk -F'"' '{print $2}')
fi

git tag -a ${version} -m "auto release"
goreleaser --skip-validate
