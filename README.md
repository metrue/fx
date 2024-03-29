fx
------

Poor man's function as a service.
<br/>
![CI](https://github.com/metrue/fx/workflows/ci/badge.svg)
![GitHub contributors](https://img.shields.io/github/contributors/metrue/fx)
[![CodeCov](https://codecov.io/gh/metrue/fx/branch/master/graph/badge.svg)](https://codecov.io/gh/metrue/fx)
[![Go Report Card](https://goreportcard.com/badge/github.com/metrue/fx?style=flat-square)](https://goreportcard.com/report/github.com/metrue/fx)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/metrue/fx)
![visitors](https://visitor-badge.glitch.me/badge?page_id=https://github.com/metrue/fx)
![license](https://img.shields.io/github/license/metrue/fx.svg)
[![Release](https://img.shields.io/github/release/metrue/fx.svg?style=flat-square)](https://github.com/metrue/fx/releases/latest)

## Table of Contents
- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)

## Introduction

fx is a tool to help you do Function as a Service on your own server, fx can make your stateless function a service in seconds, both Docker host and Kubernetes cluster supported. The most exciting thing is that you can write your functions with most programming languages.

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
| Perl          | Supported     | fx            | [/examples/Perl](https://github.com/metrue/fx/tree/master/examples/functions/Perl) |
| Crystal       | Supported     | [@mvrilo](https://github.com/mvrilo)       | [/examples/Crystal](https://github.com/metrue/fx/tree/master/examples/functions/Crystal) |
| R             | Working on [need your help](https://github.com/metrue/fx/issues/31)   | ||

# Installation

Binaries are available for Windows, MacOS and Linux/Unix on x86. For other architectures and platforms, follow instructions to [build fx from source](#buildtest).

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

fx will be installed into /usr/local/bin, sometimes you may need `source ~/.zshrc` or `source ~/.bashrc` to make fx available in `$PATH`.

* Windows

You can go the release page to [download](https://github.com/metrue/fx/releases) fx manually;

## Usage

```
NAME:
   fx - makes function as a service

USAGE:
   fx [global options] command [command options] [arguments...]

COMMANDS:
   up        deploy a function
   down      destroy a service
   list, ls  list deployed services
   image     manage image of service
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Deploy function

#### Local Docker environment

By default, function will be deployed on localhost， make sure [Docker](https://docs.docker.com/engine/installation/) installed and running on your server first. then type `fx -h` on your terminal to check out basic help.

```
$ fx up --name hello ./examples/functions/JavaScript/func.js

+------------------------------------------------------------------+-----------+---------------+
|                                ID                                |   NAME    |   ENDPOINT    |
+------------------------------------------------------------------+-----------+---------------+
| 5b24d36608ee392c937a61a530805f74551ddec304aea3aca2ffa0fabcf98cf3 | /hello    | 0.0.0.0:58328 |
+------------------------------------------------------------------+-----------+---------------+
```

#### Remote host

Use `--host` to specify the target host for your function, or you can just set it to `FX_HOST` environment variable.

```shell
$ fx up --host roo@<your host> --name hello ./examples/functions/JavaScript/func.js

+------------------------------------------------------------------+-----------+---------------+
|                                ID                                |   NAME    |   ENDPOINT    |
+------------------------------------------------------------------+-----------+---------------+
| 5b24d36608ee392c937a61a530805f74551ddec304aea3aca2ffa0fabcf98cf3 | /hello    | 0.0.0.0:58345 |
+------------------------------------------------------------------+-----------+---------------+
```

#### Kubernetes

```
$ FX_KUBECONF=~/.kube/config fx up examples/functions/JavaScript/func.js --name hello

+-------------------------------+------+----------------+
| ID                     | NAME        |    ENDPOINT    |
+----+--------------------------+-----------------------+
|  5b24d36608ee392c937a  | hello-fx    | 10.0.242.75:80 |
+------------------------+-------------+----------------+
```

### Test service

then you can test your service:

```shell
$ curl -v 0.0.0.0:58328


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

## Use Public Cloud Kubernetes Service as infrastructure to run your functions

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
$ export FX_KUBECONF=~/.kube/config  # then fx will take the config to deloy function
```

But we would suggest you run `kubectl config current-context` to check if the current context is what you want.

* Amazon Elastic Kubernetes Service (EKS)
  TODO

* Google Kubernetes Engine (GKE)

First you should create a Kubernetes cluster in your GKE, then make sure your KUBECONFIG is ready in `~/.kube/config`, if not, you can run following commands,

``` shell
$ gcloud auth login
$ gcloud container clusters get-credentials <your cluster> --zone <zone> --project <project>
```

Then make sure you current context is GKE cluster, you can check it with command,

``` shell
$ kubectl config current-context
```

Then you can deploy your function onto GKE cluster with,

```shell
$ FX_KUBECONF=~/.kube/config fx up examples/functions/JavaScript/func.js --name hellojs
```

* Setup your own Kubernetes cluster

```shell
fx infra create --type k3s --name fx-cluster-1 --master root@123.11.2.3 --agents 'root@1.1.1.1,root@2.2.2.2'
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
