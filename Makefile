OUTPUT_DIR=./build
DIST_DIR=./dist

install-deps:
	@dep ensure
generate:
	@go generate api/fx.go
build: generate
	go build -o ${OUTPUT_DIR}/fx fx.go
cross:
	goreleaser --snapshot --skip-publish --skip-validate
release:
	goreleaser --skip-validate
clean:
	rm -rf ${OUTPUT_DIR}
	rm -rf ${DIST_DIR}
zip:
	zip -r images.zip images/
.PHONY: test build start list clean generate
