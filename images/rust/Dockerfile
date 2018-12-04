FROM rust:1.23.0

WORKDIR /usr/src/myapp
COPY . .
RUN cargo install
EXPOSE 3000
CMD ["cargo", "run"]
