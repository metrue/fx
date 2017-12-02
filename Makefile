OUTPUT_DIR=./build
DIST_DIR=./dist

install-deps:
	@dep ensure -update
	@dep ensure
build:
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
.PHONY: test build start list clean
