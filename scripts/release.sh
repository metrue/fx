#!/usr/bin/env bash

set -e

git config --global user.email "h.minghe@gmail.com"
git config --global user.name "Minghe Huang"

branch=$(git rev-parse --abbrev-ref HEAD)
commit=$(git rev-parse --short HEAD)
version=$(cat fx.go| grep Version | awk -F'"' '{print $2}')
if [[ ${branch} == "master" ]];then
  version=${version}-alpha.${commit}
  echo "alpha release $version"
elif [[ "${branch}" == *--autodeploy ]];then
  version=${version}-alpha.${commit}
  echo "alpha release $version"
elif [[ ${branch} == "production" ]];then
  echo "official release $version"
else
  exit 0
fi

git tag -a ${version} -m "auto release"
curl -sL https://git.io/goreleaser | bash -s  -- --skip-validate
