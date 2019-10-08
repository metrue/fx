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
