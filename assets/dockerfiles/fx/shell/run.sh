#!/bin/bash

script=$1

filename=$(basename -- "$script")
extension="${filename##*.}"

BIN="node"
case "${extension}" in
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

${BIN} ${script}
