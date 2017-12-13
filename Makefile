OUTPUT_DIR=./build
DIST_DIR=./dist

build-assets:
	go-bindata -pkg common -o common/asset.go ./assets/dockerfiles/fx/...
install-deps:
	# since we need go-bindata to build the asserts/dockerfiles/fx/* to binary
	go get -u github.com/jteeuwen/go-bindata/...
	dep ensure -v
build: build-assets
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
