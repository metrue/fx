#!/usr/bin/env bash

set -e

has() {
  type "$1" > /dev/null 2>&1
}

get_package_url() {
    platform=""
    if [ "$(uname)" == "Darwin" ]; then
        platform="macOS"
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        platform="Tux"
    elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ]; then
        platform="windows"
    elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW64_NT" ]; then
        platform="windows"
    fi

    curl https://api.github.com/repos/metrue/fx/releases/latest | grep browser_download_url | awk -F'"' '{print $4}' | grep ${platform}
}

download_and_install() {
    url=$(get_package_url)
    tarFile="fx.tar.gz"
    targetFile=$(pwd)

    userid=$(id -u)
    if [ "$userid" != "0" ]; then
      tarFile="$(pwd)/${tarFile}"
    else
      tarFile="/tmp/${tarFile}"
      targetFile="/usr/local/bin"
    fi

    if [ -e $tarFile ]; then
      rm -rf $tarFile
    fi

    echo "Downloading fx from $url"
    curl -sSLf $url --output $tarFile
    if [ "$?" == "0" ]; then
      echo "Download complete, saved to $tarFile"
    fi

    echo "Installing fx to ${targetFile}"
    tar -xvzf ${tarFile} --exclude=*.md -C ${targetFile}
    echo "fx installed successfully at ${targetFile}"
    ${targetFile}/fx -v

    echo "Cleaning up ${tarFile}"
    rm -rf ${tarFile}
}

main() {
    if has "curl";then
      download_and_install
    else
      echo "You need cURL to use this script"
      exit 1
    fi
}

# main
main
