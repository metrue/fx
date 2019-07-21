OUTPUT_DIR=./build
DIST_DIR=./dist

lint:
	golangci-lint run --no-config \
		--issues-exit-code=0 \
		--deadline=30m \
		--disable-all \
		--enable=deadcode \
		--enable=gocyclo \
		--enable=golint \
		--enable=varcheck \
		--enable=structcheck \
		--enable=maligned \
		--enable=errcheck \
		--enable=dupl \
		--enable=ineffassign \
		--enable=interfacer \
		--enable=unconvert \
		--enable=goconst \
		--enable=gosec \
		--enable=megacheck

generate:
	packr

build: generate
	go build -o ${OUTPUT_DIR}/fx fx.go

pull:
	./scripts/pull.sh

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

zip:
	zip -r images.zip images/
.PHONY: test build start list clean generate
