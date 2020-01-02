# Make a Perl function a service with fx

[![asciicast](https://asciinema.org/a/aXpr0jquwhhwhghiDCdC7nY8r.svg)](https://asciinema.org/a/aXpr0jquwhhwhghiDCdC7nY8r)


### Hello World

```perl
sub fx {
  return 'hello fx'
}

1;
```

then deploy it with `fx up` command,

```shell
$ fx up -p 8080 --name helloworld func.pl
```

test it using `curl`

```shell
$ curl 127.0.0.1:8080

HTTP/1.1 200 OK
Connection: keep-alive
Content-Length: 11
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Aug 2019 15:58:41 GMT

hello fx
```

### Sum

```perl
sub fx {
  my $ctx = shift;
  my $a = $ctx->req->json->{"a"};
  my $b = $ctx->req->json->{"b"};
  return int($a) + int($b)
}

1;
```

```shell
fx up --name add --port 40002 --force add.pl
```

Then test it with httpie.
```shell
$ http post 0.0.0.0:40002 a=1 b=2

HTTP/1.1 200 OK
Content-Length: 1
Content-Type: application/json;charset=UTF-8
Date: Thu, 02 Jan 2020 15:39:49 GMT
Server: Mojolicious (Perl)

3
```

### ctx

The `ctx` object is exactly the [Controller](https://mojolicious.org/perldoc/Mojolicious/Controller) of [Mojolicious](https://mojolicious.org/perldoc/Mojolicious) framework.
