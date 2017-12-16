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

### Introduction

fx is a tool to help you do Function as a Service on your own server. fx can make your stateless function a service in seconds. The most exciting thing is that you can write your functions with most programming languages, you can refer to the [doc](https://github.com/metrue/fx/blob/master/NEW_LANGUAGE_SUPPORT.md) to make fx support the language not listed bellow.

| Language      | Status        | Contributor   |
| ------------- |:-------------:|:-------------:|
| Go            | Supported     | fx            |
| Node          | Supported     | fx            |
| Python        | Supported     | fx            |
| Ruby          | Supported     | fx            |
| Java          | Supported     | fx            |
| PHP           | Supported     | [@chlins](https://github.com/chlins)|
| Julia         | Supported     | [@mbesancon](https://github.com/mbesancon)|
| R             | Working on [need your help](https://github.com/metrue/fx/issues/31)   | |

Welcome to tweet [me](https://twitter.com/_metrue) or [Buy me a coffee](https://www.paypal.me/minghe).

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
    </tr>
  </tbody>
</table>


### Installation

* MacOS

```
brew tap metrue/homebrew-fx
brew install fx
```

* Linux/Unix

To install fx, you can use the [install script](https://github.com/metrue/fx/blob/master/bin/install.sh) using cURL:

```
curl -o- https://raw.githubusercontent.com/metrue/fx/master/bin/install.sh | bash
```

or Wget:

```
wget -qO- https://raw.githubusercontent.com/metrue/fx/master/bin/install.sh | bash
```

fx will be installed into /usr/local/bin, if fx not found after installation, you may need to checkout if `/usr/local/bin/fx` exists.
sometimes you may need `source ~/.zshrc` or `source ~/.bashrc` to make fx available on $PAHT.

* Window

You can go the release page to [download](https://github.com/metrue/fx/releases) fx manually;

### Usage

Make sure [Docker](https://docs.docker.com/engine/installation/) installed and running on your server first.

* start server

```
fx serve
```

now you can make a function to service in a second.

```
fx up ./examples/functions/func.js
```

the function defined in *examples/functions/func.js* is quite simple, it calculates the sum of two numbers then returns:
```
module.exports = (input) => {
    return parseInt(input.a, 10) + parseInt(input.b, 10)
}
```

then you can test your service:
```
curl -X POST <service url> -H "Content-Type: application/json" -d '{"a": 1, "b": 1}'
```

of course you can do more.

```
Usage:
$ fx serve                                      start f(x) server
$ fx up   func1.js func2.py func3.go ...        deploy a function or a group of functions
$ fx down [service ID] ...                      destroy a function or a group of functions
$ fx list                                       list deployed services
$ fx --version                                  show current version of f(x)
```

#### How to write your function

functions example with Go, Ruby, Python, Node, PHP, Java, Julia.

* Go
```
package main

type Input struct {
	A int32
	B int32
}

type Output struct {
	Sum int32
}

func Fx(input *Input) (output *Output) {
	output = &Output{
		Sum: input.A + input.B,
	}
	return
}
```

* Ruby
```
def fx(input)
    return input['a'] + input['b']
end
```

* Java
```
package fx;

import org.json.JSONObject;

public class Fx {
    public int handle(JSONObject input) {
        String a = input.get("a").toString();
        String b = input.get("b").toString();
        return Integer.parseInt(a) + Integer.parseInt(b);
    }
}
```

* Python
```
def fx(input):
    return input['a'] + input['b']
```

* Node
```
module.exports = (input) => {
    return parseInt(input.a, 10) + parseInt(input.b, 10)
}
```

* PHP
```
<?php
    function Fx($input) {
        return $input["a"]+$input["b"];
    }
```

* Julia
```
struct Input
    a::Number
    b::Number
end

fx = function(input::Input)
    return input.a + input.b
end
```

### Contributing

##### Requirements
* Docker: make sure [Docker](https://docs.docker.com/engine/installation/) installed and running on your server.
* dep: fx project uses [dep](https://github.com/golang/dep) to do the golang dependency management.
* protoc / grpc: Used for RPC and types definition (See a [setup script](https://gist.github.com/muka/4cc42c478b2699f0969450a1ec1ce44c) example)

##### Build and Run

```
$ git clone https://github.com/metrue/fx.git
$ cd fx
$ make install-deps && make build
$ ./build/fx serve                      # start your fx server
$ ./build/fx up func.js                 # deploy a function
```

### LICENSE

MIT
