fx
------

Poor man's function as a service.
<p>
  <img src="https://circleci.com/gh/metrue/fx.svg?style=svg&circle-token=bd62abac47802f8504faa4cf8db43e4f117e7cd7"/>
</p>

### Usage

* clone and build

```
$ git clone https://github.com/metrue/fx
$ make install-deps && make build
```

* start server

```
sudo ./build/fx server start            # since fx server is running as system service, so sudo needed
```

now you can make a function to service in a second.

```
./build/fx up fx/example/functions/func.js
```

of course you can do more.

```
Usage:
$ fx up   func1 func2 ...       deploy a function or a group of functions
$ fx down func1 func2 ...       destroy a function or a group of functions
$ fx list                       list deployed services
$ fx server start               start fx server
$ fx server stop                stop fx server
$ fx server status              show status of fx server
$ fx --version                  show current version of f(x)
```

### Architecture

TODO

### Features

* no API Gateway
* no Function Watchdog
* no Docker Swarm
* no Kubernets
* no fancy web dashboard

but f(x)

* **makes a function to be a service in seconds**.
* **supports all major programming languages (Node, Golang, Ruby, Python) functions to services**.


### Contributors

<table>
  <tbody>
    <tr>
      <td align="center" valign="top">
        <img width="150" height="150" src="https://github.com/metrue.png?s=150">
        <br>
        <a href="https://github.com/metrue">Minghe</a>
        <p>Core</p>
        <br>
        <p>Founder of f(x)</p>
      </td>
      <td align="center" valign="top">
        <img width="150" height="150" src="https://github.com/pplam.png?s=150">
        <br>
        <a href="https://github.com/pplam">Tim</a>
        <p>Core</p>
        <br>
        <p>Founder of f(x)</p>
      </td>
     </tr>
  </tbody>
</table>
