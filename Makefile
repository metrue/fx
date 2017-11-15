OUTPUT_DIR=./build
DIST_DIR=./dist

install-deps:
	@glide install -v
build:
	go build -o ${OUTPUT_DIR}/fx fx.go
cross:
	goreleaser --snapshot --skip-publish --skip-validate
release:
	goreleaser --skip-validate
clean:
	rm -rf ${OUTPUT_DIR}
	rm -rf ${DIST_DIR}
.PHONY: test build start list clean
