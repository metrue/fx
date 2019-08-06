# Make a JavaScript function a service with fx

Write a function like,

```JavaScript
module.exports = (ctx) => {
  ctx.body = 'hello world'
}
```

then deploy it with `fx up` command,

```shell
$ fx up -p 8080:3000 func.js
```

test it using `curl`

```shell
$ curl 127.0.0.1:8080

HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 11
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Aug 2019 15:58:41 GMT

hello world
```

### ctx

The `ctx` object is exactly the [ctx](https://github.com/koajs/koa/blob/master/docs/api/context.md) of [Koa](https://github.com/koajs/koa)
