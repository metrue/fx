FROM liuchong/rustup

WORKDIR /usr/src/myapp
COPY . .
RUN cp ./config ~/.cargo/ && rustup default nightly && cargo build
EXPOSE 3000
CMD ["cargo", "run"]
