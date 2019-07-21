FROM metrue/fx-golang-base:latest

COPY . /go/src/github.com/metrue/fx
WORKDIR /go/src/github.com/metrue/fx

RUN go build -ldflags "-w -s" -o fx fx.go app.go

EXPOSE 3000

CMD ["./fx"]
