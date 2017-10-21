OUTPUT_DIR=./build
GOPATH=$(shell pwd)/vendor

start:
	${OUTPUT_DIR}/fx
build:
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o ${OUTPUT_DIR}/fx fx.go
clean:
	rm -rf ${OUTPUT_DIR}
.PHONY: start build clean
