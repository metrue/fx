# Make a Perl function a service with fx

Write a function like,

```perl
sub fx {
  my $ctx = shift;
  return 'hello fx'
}

1;
```

then deploy it with `fx up` command,

```shell
$ fx up -p 8080:3000 func.pl
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

### ctx

The `ctx` object is exactly the [Controller](https://mojolicious.org/perldoc/Mojolicious/Controller) of [Mojolicious](https://mojolicious.org/perldoc/Mojolicious) framework.
