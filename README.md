fx
------

Poor man's function as a service.
<br/>
![build](https://circleci.com/gh/metrue/fx.svg?style=svg&circle-token=bd62abac47802f8504faa4cf8db43e4f117e7cd7)
[![Go Report Card](https://goreportcard.com/badge/github.com/metrue/fx?style=flat-square)](https://goreportcard.com/report/github.com/metrue/fx)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/metrue/fx)
[![Release](https://img.shields.io/github/release/metrue/fx.svg?style=flat-square)](https://github.com/metrue/fx/releases/latest)

### Introduction

fx is a tool to help you do Function as a Service on your own server. fx can make your stateless function a service in seconds. The most exciting thing is that you can write your functions with most programming languages.

| Language      | Status        |
| ------------- |:-------------:|
| Go            | Supported     |
| Node          | Supported     |
| Python        | Supported     |
| Ruby          | Supported     |
| PHP           | Supported     |
| Perl          | Working on    |
| R             | Working on    |
| Rust          | Working on    |

tweet [@_metrue](https://twitter.com/_metrue) or issue is welcome.

### Usage

##### Requirements
* Docker: make sure [Docker](https://docs.docker.com/engine/installation/) installed and running on your server.
* dep: fx project uses [dep](https://github.com/golang/dep) to do the golang dependency management.

##### Build and Run

```
$ git clone https://github.com/metrue/fx.git
$ cd fx
$ dep ensure
$ go install ./
```

* start server

```
fx serve
```

now you can make a function to service in a second.

```
fx up ./example/functions/func.js
```

the function defined in *exmaple/functions/func.js* is quite simple, it calculates the sum of two numbers then returns:
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

functions example with Go, Ruby, Python, Node, PHP.

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

### Contributors

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
    </tr>
  </tbody>
</table>
