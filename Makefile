OUTPUT_DIR=./build
DIST_DIR=./dist

lint:
	golangci-lint run

generate:
	packr

b:
	go build -o ${OUTPUT_DIR}/fx fx.go

build:
	go build -o ${OUTPUT_DIR}/fx fx.go

pull:
	./scripts/pull.sh

cross: generate
	goreleaser --snapshot --skip-publish --skip-validate

clean:
	rm -rf ${OUTPUT_DIR}
	rm -rf ${DIST_DIR}

unit-test:
	./scripts/coverage.sh

cli-test:
	echo 'run testing on localhost'
	./scripts/test_cli.sh
	echo 'run testing on remote host'
	DOCKER_REMOTE_HOST_ADDR=${REMOTE_HOST_ADDR} DOCKER_REMOTE_HOST_USER=${REMOTE_HOST_USER} DOCKER_REMOTE_HOST_PASSWORD=${REMOTE_HOST_PASSWORD} ./scripts/test_cli.sh

http-test:
	./scripts/http_test.sh

zip:
	zip -r images.zip images/
.PHONY: test build start list clean generate
