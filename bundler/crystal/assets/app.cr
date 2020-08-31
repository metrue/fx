require "kemal"
require "gzip"
require "json"
require "./fx"

post "/" do |ctx|
    fx ctx
end

get "/" do |ctx|
    fx ctx
end

Kemal.run
