# fx

### Up and Run

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

```Shell
$ fx up -p 8080:3000 func.go
$ curl 127.0.0.1

HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 11
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Aug 2019 15:58:41 GMT

hello world
```

### ctx

The `ctx` object is exactly the [context](https://github.com/gin-gonic/gin/blob/master/context.go#L43) of [Gin](https://github.com/gin-gonic/gin)
