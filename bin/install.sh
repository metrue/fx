#!/usr/bin/env bash

# v0.0.1
# https://github.com/metrue/fx/releases/download/v0.0.1/fx_0.0.1_checksums.txt
# https://github.com/metrue/fx/releases/download/v0.0.1/fx_0.0.1_macOS_64-bit.tar.gz
# https://github.com/metrue/fx/releases/download/v0.0.1/fx_0.0.1_Tux_64-bit.tar.gz
# https://github.com/metrue/fx/releases/download/v0.0.1/fx_0.0.1_windows_64-bit.tar.gz

get_package_url() {
    local label=$1
    curl -s https://api.github.com/repos/metrue/fx/releases/latest | grep browser_download_url | awk -F'"' '{print $4}' | grep ${label}
}

download_and_unzip() {
    local url=$1
    echo ${url}
    # TODO we can do it on one line
    rm -rf fx.tar.gz
    curl -o fx.tar.gz -L -O ${url} && tar -xvzf ./fx.tar.gz -C /usr/local/bin
}

url=""
if [ "$(uname)" == "Darwin" ]; then
    url=$(get_package_url 'macOS')
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    url=$(get_package_url 'Tux')
# elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ]; then
#     # TODO no support
# elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW64_NT" ]; then
    # TODO no support
fi

if [ ${url}"X" != "X" ];then
    download_and_unzip ${url}
fi
