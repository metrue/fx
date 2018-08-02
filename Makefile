OUTPUT_DIR=./build
DIST_DIR=./dist

install-deps:
	git submodule update --init --recursive
	go get -u github.com/olekukonko/tablewriter
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u golang.org/x/sys/...
	go get -u golang.org/x/text/...
	go get -u google.golang.org/grpc
	go get -u github.com/urfave/cli
	go get github.com/goreleaser/goreleaser
	# install protoc and plugins
	./scripts/install_protoc.sh third_party/protoc

generate:
	# generate gRPC related code
	cd api && ./gen.sh
	# bundle assert into binary
	go-bindata -pkg common -o common/asset.go ./assets/dockerfiles/fx/...

build: generate
	go build -o ${OUTPUT_DIR}/fx fx.go

cross: generate
	goreleaser --snapshot --skip-publish --skip-validate

release: generate
	./scripts/release.sh

clean:
	rm -rf ${OUTPUT_DIR}
	rm -rf ${DIST_DIR}

unit-test: generate
	./scripts/coverage.sh

cli-test: generate
	./scripts/test_cli.sh

http-test: generate
	./scripts/http_test.sh

grpc-test: generate
	echo "TODO"

zip:
	zip -r images.zip images/
.PHONY: test build start list clean generate
