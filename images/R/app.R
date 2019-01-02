library(jug)
library(jsonlite)
source('./fx.R')

jug() %>%
    post("/", function(req, res, err) {
             input <- fromJSON(req$body)
             fx(input)
    }) %>%
    simple_error_handler_json() %>%
    serve_it(port = 3000)
