# Make a Golang function a service with fx
(v0.5.4 or newer)

Write a function like,

```Go
package main

import "github.com/gin-gonic/gin"

func fx(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "hello world",
	})
}
```

into a file named ```fx_hello.go```

then deploy it with `fx up` command,

```shell
$ fx up -name fx_hello -p 10001 --healthcheck fx_hello.go
2019/08/30 10:12:11  info Build Service fx_hello: ✓
2019/08/30 10:12:11  info Run Service: ✓           
2019/08/30 10:12:11  info Service (fx_hello) is running on: 0.0.0.0:10001
2019/08/30 10:12:11  info service is running       
2019/08/30 10:12:12  info service is running       
2019/08/30 10:12:13  info service is running       
2019/08/30 10:12:14  info up function fx_hello(fx_hello.go) to machine localhost: ✓
```

test it using `curl`

```shell
$ curl -i 127.0.0.1:10001
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 30 Aug 2019 17:12:33 GMT
Content-Length: 25

{"message":"hello world"}
```

### ctx

The `ctx` object is exactly the [context](https://github.com/gin-gonic/gin/blob/master/context.go#L43) of [Gin](https://github.com/gin-gonic/gin)
