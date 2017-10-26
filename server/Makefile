OUTPUT_DIR=./build
GOPATH=$(shell pwd)/vendor
DIST_DIR=./dist

install-deps:
	@mkdir -p ./vendor/src
	@glide install --strip-vendor --strip-vcs
	@mkdir -p ./vendor/src
	@mv ./vendor/* ./vendor/src > /dev/null 2>&1;true
start:
	${OUTPUT_DIR}/fx
build:
	GOPATH=$(GOPATH) go build -o ${OUTPUT_DIR}/fx fx.go
cross:
	GOPATH=${GOPATH} goreleaser --snapshot --skip-publish --skip-validate
release:
	GOPATH=${GOPATH} goreleaser --skip-validate
clean:
	rm -rf ${OUTPUT_DIR}
	rm -rf ${DIST_DIR}
.PHONY: start build clean
