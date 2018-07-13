#!/bin/bash

code=$1
lang=$2

echo "${code}" | base64 -d > fxScript

BIN="node"
case "${lang}" in
  js)
    BIN="nodejs"
    ;;
  rb)
    BIN="ruby"
    ;;
  py)
    BIN="python"
    ;;
  pl)
    BIN="perl"
    ;;
  *)
    echo "Not support run ${extension} yet"
    exit 1
    ;;
esac

${BIN} fxScript
