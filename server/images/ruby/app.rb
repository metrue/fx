require 'sinatra'

require_relative 'fx.rb'

set :port, 3000

post '/' do
    request.body.rewind
    request_payload = JSON.parse request.body.read
    body fx request_payload
end

