FROM ubuntu:16.04

RUN apt-get update && apt-get install -y build-essential curl libcurl3 \
  && curl -OL http://downloads.dlang.org/releases/2.x/2.077.1/dmd_2.077.1-0_amd64.deb \
  && apt install ./dmd_2.077.1-0_amd64.deb