OUTPUT_DIR=./build
GOPATH=$(shell pwd)/vendor

install-deps:
	@mkdir -p ./vendor/src
	@glide install #--strip-vendor
	@mkdir -p ./vendor/src
	@mv ./vendor/* ./vendor/src > /dev/null 2>&1
start:
	${OUTPUT_DIR}/fx
build:
	GOPATH=$(GOPATH) go build -o ${OUTPUT_DIR}/fx fx.go
clean:
	rm -rf ${OUTPUT_DIR}
.PHONY: start build clean
