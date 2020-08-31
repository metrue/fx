FROM metrue/fx-crystal-base

EXPOSE 3000

COPY . .

RUN crystal build --verbose  -o ./fxcr ./*.cr

CMD ["./fxcr"]
