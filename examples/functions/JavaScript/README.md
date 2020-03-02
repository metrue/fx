# Make a JavaScript function a service with fx

You have a function define like this in `func.js`

```JavaScript
module.exports = (ctx) => {
  ctx.body = 'hello world'
}
```

## Deploy function on localhost

```shell
$ fx up -p 8080 func.js
```

test it using `curl`

```shell
$ curl -v 127.0.0.1:8080

HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 11
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Aug 2019 15:58:41 GMT

hello world
```

## Deploy function to remote host

You have to make sure you can `ssh root@<your host>` without password, you can refer to [SSH login without password](http://www.linuxproblem.org/art_9.html) on how to enable ssh to host without
password.

```
$ export SSH_PORT=11398                                                       # if you changed the SSH port of you remote host
$ fx infra create --name <infra_name> --type docker --host root@<your_host>   # make this host to be one of fx's target hosts
$ fx infra <infra_name>                                                       # make this host to be the active host
$ fx infra list                                                               # show fx infra status
```

Then you can deploy your function to host with,
```
$ fx up --name hello --port 8080 fun.js
```

### What's the ctx of fx function

The `ctx` object is exactly the [ctx](https://github.com/koajs/koa/blob/master/docs/api/context.md) of [Koa](https://github.com/koajs/koa)
