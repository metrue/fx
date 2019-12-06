# Make a Rust function a service with fx

Write a function like,

```rust
pub mod fns {
    #[derive(Serialize)]
    pub struct Response {
        pub result: i32,
    }

    #[derive(Deserialize)]
    pub struct Request {
        pub a: i32,
        pub b: i32,
    }

    pub fn func(req: Request) -> Response {
        Response {
            result: req.a + req.b,
        }
    }
}

```

then deploy it with `fx up` command,

```shell
$ fx up -p 8080 func.rs
```

test it using `curl`

```shell
$ curl -X 'POST' --header 'Content-Type: application/json' --data '{"a":1,"b":1}' '0.0.0.0:3000'

HTTP/1.1 200 OK
Content-Length: 12
Content-Type: application/json
Date: Fri, 06 Dec 2019 06:45:14 GMT
Server: Rocket

{
    "result": 2
}
```
