FROM julia:0.7

COPY . /app

RUN apt-get update && apt-get install -y gcc apt-utils unzip make libhttp-parser-dev
RUN julia /app/deps.jl

CMD julia /app/app.jl
EXPOSE 3000
