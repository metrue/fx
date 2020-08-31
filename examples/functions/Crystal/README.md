# Make a Crystal function a service with fx

Write a function like,

```crystal
def fx(ctx)
  "hello world, crystal"
end
```

then deploy it with `fx up` command,

```shell
$ fx up -p 8000 func.cr
```

test it using `curl`

```shell
$ curl -i localhost:8000

HTTP/1.1 200 OK
Connection: keep-alive
X-Powered-By: Kemal
Content-Type: text/html
Content-Length: 20

hello world, crystal
```

### ctx

The `ctx` argument is a Kemal HTTP request / response context [context](https://kemalcr.com/guide/#context)
