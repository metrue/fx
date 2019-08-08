# Make a Ruby function a service with fx

Write a function like,

```ruby
def fx(ctx)
  ctx[:response].body = "hello world"
end
```

then deploy it with `fx up` command,

```shell
$ fx up -p 8080:3000 func.go
```

test it using `curl`

```shell
$ curl 127.0.0.1:8080

HTTP/1.1 200 Created
Connection: Keep-Alive
Content-Length: 11
Content-Type: text/html;charset=utf-8
Date: Thu, 08 Aug 2019 02:39:55 GMT
Server: WEBrick/1.4.2 (Ruby/2.6.3/2019-04-16)
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-Xss-Protection: 1; mode=block

hello world
```

### ctx

The `ctx` object is a Hash have contains [request](https://rubydoc.info/github/rack/rack/master/Rack/Request<Paste>), [response](https://rubydoc.info/github/rack/rack/master/Rack/Response), [status](https://www.rubydoc.info/gems/sinatra/Sinatra%2FHelpers:status), and [headers](https://www.rubydoc.info/gems/sinatra/Sinatra%2FHelpers:headers) of [Sinatra](https://github.com/sinatra/sinatra)
