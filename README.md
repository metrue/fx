fx
------
Poor man's function as a service.
<br/>
![ci](https://github.com/metrue/fx/workflows/ci/badge.svg)
![build](https://circleci.com/gh/metrue/fx.svg?style=svg&circle-token=bd62abac47802f8504faa4cf8db43e4f117e7cd7)
[![codecov](https://codecov.io/gh/metrue/fx/branch/master/graph/badge.svg)](https://codecov.io/gh/metrue/fx)
[![Go Report Card](https://goreportcard.com/badge/github.com/metrue/fx?style=flat-square)](https://goreportcard.com/report/github.com/metrue/fx)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/metrue/fx)
![](https://img.shields.io/github/license/metrue/fx.svg)
[![Release](https://img.shields.io/github/release/metrue/fx.svg?style=flat-square)](https://github.com/metrue/fx/releases/latest)

## Table of Contents
- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Contribute](#contribute)


## Introduction

fx is a tool to help you do Function as a Service on your own server. fx can make your stateless function a service in seconds. The most exciting thing is that you can write your functions with most programming languages.

Feel free hacking fx to support the languages not listed. Welcome to tweet me [@_metrue](https://twitter.com/_metrue) on Twitter, [@metrue](https://www.weibo.com/u/2165714507) on Weibo.


| Language      | Status        | Contributor   | Example        |
| ------------- |:-------------:|:-------------:| :-------------:|
| Go            | Supported     | fx            | [/examples/Golang](https://github.com/metrue/fx/tree/master/examples/functions/Golang) |
| Rust          | Supported     | [@FrontMage](https://github.com/FrontMage)| [/examples/Rust](https://github.com/metrue/fx/tree/master/examples/functions/Rust) |
| Node          | Supported     | fx            | [/examples/JavaScript](https://github.com/metrue/fx/tree/master/examples/functions/JavaScript) |
| Python        | Supported     | fx            | [/examples/Python](https://github.com/metrue/fx/tree/master/examples/functions/Python) |
| Ruby          | Supported     | fx            | [/examples/Ruby](https://github.com/metrue/fx/tree/master/examples/functions/Ruby) |
| Java          | Supported     | fx            | [/examples/Java](https://github.com/metrue/fx/tree/master/examples/functions/Java) |
| PHP           | Supported     | [@chlins](https://github.com/chlins)| [/examples/PHP](https://github.com/metrue/fx/tree/master/examples/functions/PHP) |
| Julia         | Supported     | [@matbesancon](https://github.com/matbesancon)| [/examples/Julia](https://github.com/metrue/fx/tree/master/examples/functions/Julia) |
| D             | Supported     | [@andre2007](https://github.com/andre2007)| [/examples/D](https://github.com/metrue/fx/tree/master/examples/functions/D) |
| R             | Working on [need your help](https://github.com/metrue/fx/issues/31)   | ||

# Installation

* MacOS

```
brew tap metrue/homebrew-fx
brew install metrue/fx/fx
```

* Linux/Unix

via cURL

```shell
# Install to local directory
curl -o- https://raw.githubusercontent.com/metrue/fx/master/scripts/install.sh | bash

# Install to /usr/local/bin/
curl -o- https://raw.githubusercontent.com/metrue/fx/master/scripts/install.sh | sudo bash
```

fx will be installed into /usr/local/bin, sometimes you may need `source ~/.zshrc` or `source ~/.bashrc` to make fx available in `$PAHT`.

* Window

You can go the release page to [download](https://github.com/metrue/fx/releases) fx manually;

## Usage

Make sure [Docker](https://docs.docker.com/engine/installation/) installed and running on your server first. then type `fx -h` on your terminal to check out basic help.

```
NAME:
   fx - makes function as a service

USAGE:
   fx [global options] command [command options] [arguments...]

VERSION:
   0.6.0

COMMANDS:
   infra     manage infrastructure of fx
   image     manage image of service
   doctor    health check for fx
   up        deploy a function or a group of functions
   down      destroy a service
   list, ls  list deployed services
   call      run a function instantly
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

1. List your current machines and activate you machine

```shell
$ fx infra ls     # list machines

{
	"localhost": {
		"Host": "localhost",
		"User": "",
		"Password": "",
		"Enabled": true,
		"Provisioned": false
	}
}

$ fx infra activate localhost    # activate 'localhost'

2019/08/10 13:21:20  info Provision:pull python Docker base iamge: ✓
2019/08/10 13:21:21  info Provision:pull d Docker base image: ✓
2019/08/10 13:21:23  info Provision:pull java Docker base image: ✓
2019/08/10 13:21:28  info Provision:pull julia Docker base image: ✓
2019/08/10 13:21:31  info Provision:pull node Docker base image: ✓
2019/08/10 13:22:09  info Provision:pull go Docker base image: ✓
2019/08/10 13:22:09  info provision machine localhost: ✓
2019/08/10 13:22:09  info enble machine localhost: ✓
```
It may take seconds since `fx` needs to download some basic resources

*Note* you can add a remote host as fx machine also,
```
$ fx infra add --name my_aws_vm --host 13.121.202.227 --user root --password yourpassword

$ fx infra list
{
	"my_aws_vm": {
		"Host": "13.121.202.227",
		"User": "root",
		"Password": "yourpassword",
		"Enabled": false,
		"Provisioned": false
	},
	"localhost": {
		"Host": "localhost",
		"User": "",
		"Password": "",
		"Enabled": true,
		"Provisioned": true
	}
}

$ fx infra activate my_aws_vm

```
then your function will be deployed onto remote host also.

2. Write a function

You can check out [examples](https://github.com/metrue/fx/tree/master/examples/functions) for reference. Let's write a function as an example,  it calculates the sum of two numbers then returns:

```js
module.exports = (ctx) => {
    ctx.body = 'hello world'
}
```
Then save it to a file `func.js`.

3. Deploy your function as a service

Give your service a port with `--port`, and name with `--name`, heath checking with `--healthcheck` if you want.

```shell
$ fx up -name fx_service_name -p 10001 --healthcheck func.js

2019/08/10 13:26:37  info Pack Service: ✓
2019/08/10 13:26:39  info Build Service: ✓
2019/08/10 13:26:39  info Run Service: ✓
2019/08/10 13:26:39  info Service (fx_service_name) is running on: 0.0.0.0:10001
2019/08/10 13:26:39  info up function fx_service_name(func.js) to machine localhost: ✓
```

if you want see what the source code of your service looks like, you can export it into a dirctory,

```shell
$ fx image export -o <path of dir> func.js
2019/09/25 19:31:19  info exported to <path of dir>: ✓
```

4. Test your service

then you can test your service:

```shell
$ curl -v 0.0.0.0:10001


GET / HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: 0.0.0.0:10001
User-Agent: HTTPie/1.0.2



HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 11
Content-Type: text/plain; charset=utf-8
Date: Sat, 10 Aug 2019 05:28:03 GMT

hello world

```

## Docker

  TODO

## Kubernetes

**fx** supports deploy function to be a service onto Kubernetes cluster infrasture, and we encourage you to do that other than on bare Docker environment, there are lots of advantage to run your function on Kubernetes like self-healing, load balancing, easy horizontal scaling, etc. It's pretty simple to deploy your function onto Kubernetes with **fx**, you just set KUBECONFIG in your enviroment.

```shell
KUBECONFIG=<Your KUBECONFIG> fx deploy -n fx-service-abc_js -p 12349 examples/functions/JavaScript/func.js   # function will be deploy to your Kubernetes cluster and expose a IP address of your loadbalencer
```

or

```shell
$ export KUBECONFIG=<Your KUBECONFIG>
$ fx deploy -n fx-service-abc_js -p 12349 examples/functions/JavaScript/func.js   # function will be deploy to your Kubernetes cluster and expose a IP address of your loadbalencer
```

* Local Kubernetes Cluster

Docker for Mac and Docker for Windows already support Kubernetes with single node cluster, we can use it directly, and the default `KUBECONFIG` is `~/.kube/config`.

```shell
$ export KUBECONFIG=~/.kube/config  # then fx will take the config to deloy function
```

if you have multiple Kubernetes clusters configured, you have to set context correctly. FYI [configure-access-multiple-clusters](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/)

* Azure Kubernetes Service (AKS)

You should create a Kubernetes cluster if you don't have one on AKS, detail document can found [here](https://docs.microsoft.com/en-us/azure/aks/kubernetes-walkthrough).

```shell
$ az group create --name <myResourceGroup> --location eastus
$ az aks create --resource-group <myResourceGroup> --name myAKSCluster --node-count <count>
$ az aks get-credentials --resource-group <myResourceGroup> --name <myAKSCluster>
```

Then you can verify it with `kubectl`,

```shell
$ kubectl get nodes

NAME                       STATUS   ROLES   AGE     VERSION
aks-nodepool1-31718369-0   Ready    agent   6m44s   v1.12.8
```

Since AKS's config will be merged into `~/.kube/config` and set to be current context after you run `az aks get-credentials` command, so you can just set KUBECONFIG to default config also,

```shell
$ export KUBECONFIG=~/.kube/config  # then fx will take the config to deloy function
```

But we would suggest you run `kubectl config current-context` to check if the current context is what you want.

* Amazon Elastic Kubernetes Service (EKS)
  TODO

* Google Kubernetes Engine (GKET)
  TODO

## Contribute

fx uses [Project](https://github.com/metrue/fx/projects/4) to manage the development.

#### Prerequisites

Docker: make sure [Docker](https://docs.docker.com/engine/installation/) installed and running on your server.


#### Build & Test

```
$ git clone https://github.com/metrue/fx
$ cd fx
$ make build
```

Then you can build and test:

```
$ make build
$ ./build/fx -h
```


## Contributors

Thank you to all the people who already contributed to fx!

<table>
  <tbody>
    <tr>
        <a href="https://github.com/metrue" target="_blank">
            <img alt="metrue" src="https://avatars2.githubusercontent.com/u/1001246?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/muka" target="_blank">
            <img alt="muka" src="https://avatars2.githubusercontent.com/u/1021269?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/pplam" target="_blank">
            <img alt="pplam" src="https://avatars2.githubusercontent.com/u/12783579?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/matbesancon" target="_blank">
            <img alt="mbesancon" src="https://avatars2.githubusercontent.com/u/7623090?s=60&v=4" width="50">
        </a>
        <a href="https://github.com/chlins" target="_blank">
            <img alt="chlins" src="https://avatars2.githubusercontent.com/u/31262637?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/xwjdsh" target="_blank">
            <img alt="xwjdsh" src="https://avatars2.githubusercontent.com/u/11025519?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/DaidoujiChen" target="_blank">
            <img alt="DaidoujiChen" src="https://avatars0.githubusercontent.com/u/670441?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/avelino" target="_blank">
            <img alt="avelino" src="https://avatars2.githubusercontent.com/u/31996?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/andre2007" target="_blank">
            <img alt="andre2007" src="https://avatars1.githubusercontent.com/u/1451047?s=50&v=4" width="50">
        </a>
        <a href="https://github.com/polyrabbit" target="_blank">
            <img alt="polyrabbit" src="https://avatars0.githubusercontent.com/u/2657334?s=60&v=4" width="50">
        </a>
        <a href="https://github.com/johnlunney" target="_blank">
            <img alt="johnlunney" src="https://avatars3.githubusercontent.com/u/536947?s=60&v=4" width="50">
        </a>
        <a href="https://github.com/tbrand" target="_blank">
            <img alt="tbrand" src="https://avatars0.githubusercontent.com/u/3483230?s=60&v=4" width="50">
        </a>
        <a href="https://github.com/steventhanna" target="_blank">
            <img alt="andre2007" src="https://avatars1.githubusercontent.com/u/2541678?s=50&v=4" width="50">
        </a>
        <a href="https://github.com/border-radius" target="_blank">
            <img alt="border-radius" src="https://avatars0.githubusercontent.com/u/3204785?s=60&v=4" width="50">
        </a>
        <a href="https://github.com/Russtopia" target="_blank">
            <img alt="Russtopia" src="https://avatars1.githubusercontent.com/u/2966177?s=60&v=4<Paste>" width="50">
        </a>
        <a href="https://github.com/FrontMage" target="_blank">
            <img alt="FrontMage" src="https://avatars2.githubusercontent.com/u/17007026?s=60&v=4" width="50">
        </a>
        <a href="https://github.com/DropNib" target="_blank">
            <img alt="DropNib" src="https://avatars0.githubusercontent.com/u/32019589?s=60&v=4" width="50">
        </a>
    </tr>
  </tbody>
</table>
