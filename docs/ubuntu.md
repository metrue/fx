# fx on Amazon Lightsai

> The guide is verified on Amazon Lightsail ubuntu 18.08 instance

## Install Docker

```shell
apt-get remove -y docker docker-engine docker.io containerd runc
apt-get update -y
apt-get install -y apt-transport-https ca-certificates curl software-properties-common lsb-core
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
apt-get update -y
apt-get install -y docker-ce
docker run hello-world
```

## Install fx

```shell
curl -o- https://raw.githubusercontent.com/metrue/fx/master/scripts/install.sh | sudo bash
```

## Deploy a function onto localhost

```shell
$ cat func.js

module.exports = (ctx) => {
  ctx.body = 'hello world'
}

# fx up -n test -p 2000 func.js
# curl 127.0.0.1:2000
```

##  Deploy a function onto remote host

* make sure your instance can be ssh login
* make sure your instance accept port 8866


If you're first time to deploy a function to remote host, you need init it first
```shell
DOCKER_REMOTE_HOST_ADDR=<your host> DOCKER_REMOTE_HOST_USER=<your user> DOCKER_REMOTE_HOST_PASSWORD=<your password> fx init
```

then you can deploy function to remote host

```shell
DOCKER_REMOTE_HOST_ADDR=<your host> DOCKER_REMOTE_HOST_USER=<your user> DOCKER_REMOTE_HOST_PASSWORD=<your password> fx up -p 2000 test/functions/func.js
```
