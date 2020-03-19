# fx on Ubuntu

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
$ curl -o- https://raw.githubusercontent.com/metrue/fx/master/scripts/install.sh | sudo bash
```

## Deploy a function onto localhost

```shell
$ cat func.js

module.exports = (ctx) => {
  ctx.body = 'hello world'
}

$ fx up -n test -p 2000 func.js
$ curl 127.0.0.1:2000
```

##  Deploy a function onto remote host

* make sure you can ssh login to target host with root

Update `/etc/ssh/sshd_config` to allow login with root.

```
PermitRootLogin yes
```

Then restart sshd with,

```shell
$ sudo service sshd restart
```

* make sure your instance accept port 8866

[FYI](https://lightsail.aws.amazon.com/ls/docs/en_us/articles/understanding-firewall-and-port-mappings-in-amazon-lightsail)

then you can deploy function to remote host

```shell
fx up --host root@<your host> test/functions/func.js 
```
