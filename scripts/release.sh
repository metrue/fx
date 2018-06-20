#!/usr/bin/env bash

set -e

branch=$(git rev-parse --abbrev-ref HEAD)
commit=$(git rev-parse --verify HEAD)
git config --global user.email "h.minghe@gmail.com"
git config --global user.name "Minghe Huang"
git tag -a v${branch}-${commit} -m "auto release"
goreleaser --skip-validate
