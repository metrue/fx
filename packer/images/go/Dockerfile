FROM golang:latest

COPY . /go/src/github.com/metrue/fx
WORKDIR /go/src/github.com/metrue/fx

# dependency management
RUN go get github.com/gin-gonic/gin

RUN go build -ldflags "-w -s" -o fx fx.go app.go

EXPOSE 3000

CMD ["./fx"]
