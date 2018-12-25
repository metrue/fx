FROM ruby:latest

RUN gem install sinatra

COPY . .
EXPOSE 3000
CMD ruby app.rb -p 3000 -o 0.0.0.0
