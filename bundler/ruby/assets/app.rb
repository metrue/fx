require 'sinatra'
require 'json'

require_relative 'fx.rb'

set :port, 3000

post '/' do
    ctx = {
      :request => request,
      :response => response,
      :status => status,
      :headers => headers,
    }
    fx ctx
end

get '/' do
    ctx = {
      :request => request,
      :response => response,
      :status => status,
      :headers => headers,
    }
    fx ctx
end
