update:
	protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative proto/fxd.proto
install:
	export GO111MODULE=on
	go get github.com/golang/protobuf/protoc-gen-go@v1.3
