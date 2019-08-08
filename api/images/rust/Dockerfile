FROM liuchong/rustup

WORKDIR /usr/src/myapp
COPY . .

RUN cp ./config ~/.cargo/

RUN rustup default nightly
RUN cargo update
RUN cargo build

EXPOSE 3000

CMD ["cargo", "run"]
