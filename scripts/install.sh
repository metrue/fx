#!/usr/bin/env bash

fx_has() {
  type "$1" > /dev/null 2>&1
}

get_package_url() {
    label=""
    if [ "$(uname)" == "Darwin" ]; then
        label="macOS"
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        label="Tux"
    elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ]; then
        label="windows"
    elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW64_NT" ]; then
        label="windows"
    fi

    curl -s https://api.github.com/repos/metrue/fx/releases/latest | grep browser_download_url | awk -F'"' '{print $4}' | grep ${label}
}

download_and_install() {
    local url=$1
    # TODO we can do it on one line
    rm -rf fx.tar.gz
    curl -o fx.tar.gz -L -O ${url} && tar -xvzf ./fx.tar.gz --exclude=*.md -C /usr/local/bin
    rm -rf ./fx.tar.gz
}

main() {
    if fx_has "docker"; then
        url=$(get_package_url)
        if [ ${url}"X" != "X" ];then
            download_and_install ${url}
        fi
    else
        echo "No Docker found on this host"
        echo "  - Docker installation: https://docs.docker.com/engine/installation"
    fi
}

# main
main
