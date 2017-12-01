OUTPUT_DIR=./build
DIST_DIR=./dist

install-deps:
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
sync:
	rsync  -avzhe ssh root@45.79.111.212:/root/.go-workspace/src/fx/vendor .
.PHONY: test build start list clean
