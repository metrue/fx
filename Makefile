OUTPUT_DIR=./build
DIST_DIR=./dist

install-deps:
	# install protobuf and grpc
	go get -u github.com/olekukonko/tablewriter
	go get -u github.com/jteeuwen/go-bindata/...

	# mkdir -p ./tmp
	# git clone --depth 1 https://github.com/googleapis/googleapis.git tmp/googleapis

	# install protoc
	./scripts/install_protoc.sh third_party
	ls -al third_party
	cd third_party/bin && pwd && ls -al .

	# install the other dependencies
	# @dep ensure
generate:
	# generate gRPC related code
	cd api && ./gen.sh
	# bundle assert into binary
	go-bindata -pkg common -o common/asset.go ./assets/dockerfiles/fx/...
build: generate
	go build -o ${OUTPUT_DIR}/fx fx.go
cross:
	goreleaser --snapshot --skip-publish --skip-validate
release:
	goreleaser --skip-validate
clean:
	rm -rf ${OUTPUT_DIR}
	rm -rf ${DIST_DIR}
test-unit: generate
	./scripts/coverage.sh
integration-test: generate
	./scripts/test_cli.sh
zip:
	zip -r images.zip images/
.PHONY: test build start list clean generate
