# Make a Ruby function a service with fx

Write a function like,

```python
def fx(requenst):
    return "hello world"
```

then deploy it with `fx up` command,

```shell
$ fx up -p 8080:3000 func.go
```

test it using `curl`

```shell
$ curl 127.0.0.1:8080

HTTP/1.0 200 OK
Content-Length: 11
Content-Type: text/html; charset=utf-8
Date: Thu, 08 Aug 2019 05:33:32 GMT
Server: Werkzeug/0.12.2 Python/3.6.3

hello world
```

### ctx

The `ctx` object is a Hash have contains [request](https://rubydoc.info/github/rack/rack/master/Rack/Request<Paste>), [response](https://rubydoc.info/github/rack/rack/master/Rack/Response), [status](https://www.rubydoc.info/gems/sinatra/Sinatra%2FHelpers:status), and [headers](https://www.rubydoc.info/gems/sinatra/Sinatra%2FHelpers:headers) of [Sinatra](https://github.com/sinatra/sinatra)
