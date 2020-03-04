#![feature(proc_macro_hygiene, decl_macro, plugin)]

#[macro_use]
extern crate rocket;
extern crate rocket_contrib;
extern crate serde;
#[macro_use]
extern crate serde_derive;
extern crate serde_json;

mod fns;

use rocket_contrib::json::Json;

#[post("/", format = "application/json", data = "<req>")]
fn index(req: Json<fns::fns::Request>) -> Json<fns::fns::Response> {
    Json(fns::fns::func(req.0))
}

fn main() {
    rocket::ignite().mount("/", routes![index]).launch();
}
