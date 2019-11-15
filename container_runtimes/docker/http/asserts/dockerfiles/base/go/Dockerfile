FROM golang:1.12.0-alpine3.9 AS builder

LABEL maintainer="h.minghe@gmail.com"

ENV GO111MODULE on
ENV GOPROXY https://goproxy.io
ENV PATH $GOPATH/bin:$PATH

RUN go get "github.com/gin-gonic/gin"

RUN apk upgrade --no-cache \
    && apk add ca-certificates
