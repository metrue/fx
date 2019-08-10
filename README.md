fx
------

Poor man's function as a service.
<br/>
![build](https://circleci.com/gh/metrue/fx.svg?style=svg&circle-token=bd62abac47802f8504faa4cf8db43e4f117e7cd7)
[![codecov](https://codecov.io/gh/metrue/fx/branch/master/graph/badge.svg)](https://codecov.io/gh/metrue/fx)
[![Go Report Card](https://goreportcard.com/badge/github.com/metrue/fx?style=flat-square)](https://goreportcard.com/report/github.com/metrue/fx)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/metrue/fx)
![](https://img.shields.io/github/license/metrue/fx.svg)
[![Release](https://img.shields.io/github/release/metrue/fx.svg?style=flat-square)](https://github.com/metrue/fx/releases/latest)

fx is a tool to help you do Function as a Service on your own server. fx can make your stateless function a service in seconds. The most exciting thing is that you can write your functions with most programming languages.

Feel free hacking fx to support the languages not listed. Welcome to tweet me [@_metrue](https://twitter.com/_metrue) on Twitter, [@metrue](https://www.weibo.com/u/2165714507) on Weibo.


| Language      | Status        | Contributor   | Example        |
| ------------- |:-------------:|:-------------:| :-------------:|
| Go            | Supported     | fx            | [/examples/Golang](https://github.com/metrue/fx/tree/master/examples/functions/Golang) |
| Rust          | Supported     | [@FrontMage](https://github.com/FrontMage)| [/examples/Rust](https://github.com/metrue/fx/tree/master/examples/functions/Rust) |
| Node          | Supported     | fx            | [/examples/Rust](https://github.com/metrue/fx/tree/master/examples/functions/JavaScript) |
| Python        | Supported     | fx            | [/examples/Python](https://github.com/metrue/fx/tree/master/examples/functions/Python) |
| Ruby          | Supported     | fx            | [/examples/Ruby](https://github.com/metrue/fx/tree/master/examples/functions/Ruby) |
| Java          | Supported     | fx            | [/examples/Java](https://github.com/metrue/fx/tree/master/examples/functions/Java) |
| PHP           | Supported     | [@chlins](https://github.com/chlins)| [/examples/PHP](https://github.com/metrue/fx/tree/master/examples/functions/PHP) |
| Julia         | Supported     | [@mbesancon](https://github.com/mbesancon)| [/examples/Julia](https://github.com/metrue/fx/tree/master/examples/functions/Julia) |
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
curl -o- https://raw.githubusercontent.com/metrue/fx/master/scripts/install.sh | bash
```

or Wget

```shell
wget -qO- https://raw.githubusercontent.com/metrue/fx/master/scripts/install.sh | bash
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
   0.5.1

COMMANDS:
   infra     manage infrastructure of fx
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

Give your service a port with `--port`, and name with `--name` if you want.

```shell
$ fx up -name fx_service_name -p 10001 func.js

2019/08/10 13:26:37  info Pack Service: ✓
2019/08/10 13:26:39  info Build Service: ✓
2019/08/10 13:26:39  info Run Service: ✓
2019/08/10 13:26:39  info Service (fx_service_name) is running on: 0.0.0.0:10001
2019/08/10 13:26:39  info up function fx_service_name(func.js) to machine localhost: ✓
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

## Contribute

fx uses [Project](https://github.com/metrue/fx/projects) to manage the development.

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
        <a href="https://github.com/pplam" target="_blank">
            <img alt="pplam" src="https://avatars2.githubusercontent.com/u/12783579?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/muka" target="_blank">
            <img alt="muka" src="https://avatars2.githubusercontent.com/u/1021269?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/xwjdsh" target="_blank">
            <img alt="xwjdsh" src="https://avatars2.githubusercontent.com/u/11025519?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/mbesancon" target="_blank">
            <img alt="mbesancon" src="https://avatars2.githubusercontent.com/u/7623090?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/avelino" target="_blank">
            <img alt="avelino" src="https://avatars2.githubusercontent.com/u/31996?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/DaidoujiChen" target="_blank">
            <img alt="DaidoujiChen" src="https://avatars0.githubusercontent.com/u/670441?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/chlins" target="_blank">
            <img alt="chlins" src="https://avatars2.githubusercontent.com/u/31262637?v=4&s=50" width="50">
        </a>
        <a href="https://github.com/andre2007" target="_blank">
            <img alt="andre2007" src="https://avatars1.githubusercontent.com/u/1451047?s=50&v=4" width="50">
        </a>
        <a href="https://github.com/steventhanna" target="_blank">
            <img alt="andre2007" src="https://avatars1.githubusercontent.com/u/2541678?s=50&v=4" width="50">
        </a>
    </tr>
  </tbody>
</table>
