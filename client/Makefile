OUTPUT_DIR=./build
GOPATH=$(shell pwd)/vendor

install-deps:
	@mkdir -p ./vendor/src
	@glide install --strip-vendor
	@mkdir -p ./vendor/src
	@mv ./vendor/* ./vendor/src > /dev/null 2>&1;true
test:
	${OUTPUT_DIR}/fx up functions/func.go
	@docker ps
build:
	GOPATH=${GOPATH} go build -o ${OUTPUT_DIR}/fx .
start:
	${OUTPUT_DIR}/fx
list:
	@${OUTPUT_DIR}/fx list
clean:
	rm -rf ${OUTPUT_DIR}
.PHONY: test build start list clean
