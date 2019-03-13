require 'sinatra'
require 'json'

require_relative 'fx.rb'

set :port, 3000

post '/' do
    request.body.rewind
    request_payload = JSON.parse request.body.read
    ret = fx request_payload
    body ret.to_json
end
