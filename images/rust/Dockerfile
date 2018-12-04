FROM shepmaster/rust-nightly

WORKDIR /usr/src/myapp
COPY . .
RUN cargo build
EXPOSE 3000
CMD ["cargo", "run"]
