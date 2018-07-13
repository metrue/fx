FROM ubuntu:14.04

# Install Node.js
RUN apt-get update --yes
RUN apt-get install --yes curl
RUN curl --silent --location https://deb.nodesource.com/setup | bash -
RUN apt-get install --yes nodejs
RUN apt-get install --yes build-essential
RUN apt-get install -y nodejs ruby python perl

ADD run.sh .

ENTRYPOINT ["./run.sh"]
