FROM clux/muslrust:nightly AS builder
WORKDIR /build
COPY . .
RUN ln -s /usr/bin/g++ /usr/bin/musl-g++ && cargo build --release

FROM scratch
WORKDIR /usr/src/myapp
COPY --from=builder /build/target/x86_64-unknown-linux-musl/release/rust /usr/src/myapp/
COPY ./Rocket.toml /usr/src/myapp/
EXPOSE 3000
ENV ROCKET_ENV=prod
CMD ["/usr/src/myapp/rust"]
