#!/usr/bin/env bash

dir=${1}
mkdir -p ${dir}

CURRENT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cp -rf ${CURRENT_DIR}/../images/node/* ${dir}
cd ${dir} && npm install
