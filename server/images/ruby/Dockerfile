FROM ruby:latest

RUN gem install sinatra

COPY . .
CMD ruby app.rb -p 3000 -o 0.0.0.0
