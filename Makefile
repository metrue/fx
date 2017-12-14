OUTPUT_DIR=./build
DIST_DIR=./dist

install-deps:
	# install protoc
	./bin/install_protoc.sh

	# install protobuf and grpc
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u google.golang.org/grpc
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

	git clone --depth 1 https://github.com/googleapis/googleapis.git vendor/github.com/googleapis
	cp -rf vendor/github.com/googleapis/google/ api/google/

	# install the other dependencies
	@dep ensure
generate:
	go generate ./api/fx.go
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
