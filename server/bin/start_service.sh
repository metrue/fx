#!/usr/bin/env bash
TARGET_DIR=${1}
SERVICE_NAME=${2}
PORT=${3}

cd ${TARGET_DIR} && docker run -d -p ${PORT}:3000 ${SERVICE_NAME}
