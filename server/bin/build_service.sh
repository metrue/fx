#!/usr/bin/env bash

TARGET_DIR=${1}
SERVICE_NAME=${2}

cd ${TARGET_DIR} && docker build -t ${SERVICE_NAME} .
