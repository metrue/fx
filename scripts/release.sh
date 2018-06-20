#!/usr/bin/env bash

set -e

branch=$(git rev-parse --abbrev-ref HEAD)
commit=$(git rev-parse --verify HEAD)
git tag -a v${branch}-${commit} -m "auto release"
goreleaser --skip-validate
