FROM julia:latest

COPY . /app

RUN apt-get update && apt-get install -y apt-utils unzip make libhttp-parser-dev
RUN julia /app/deps.jl

CMD julia /app/app.jl
