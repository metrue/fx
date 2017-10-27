OUTPUT_DIR=./build
DIST_DIR=./dist
GOPATH=$(shell pwd)/vendor

install-deps:
	@mkdir -p ./vendor/src
	@glide install --strip-vendor --strip-vcs
	@mkdir -p ./vendor/src
	@mv ./vendor/* ./vendor/src > /dev/null 2>&1;true
test:
	${OUTPUT_DIR}/fx up functions/func.go
	@docker ps
build:
	GOPATH=${GOPATH} go build -o ${OUTPUT_DIR}/fx fx.go
start:
	${OUTPUT_DIR}/fx
list:
	@${OUTPUT_DIR}/fx list
cross:
	GOPATH={GOPATH} goreleaser --snapshot --skip-publish --skip-validate
clean:
	rm -rf ${OUTPUT_DIR}
	rm -rf ${DIST_DIR}
.PHONY: test build start list clean
