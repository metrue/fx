OUTPUT_DIR=./build

test:
	${OUTPUT_DIR}/fx up index.js func.go
	@docker ps
build:
	go build -o ${OUTPUT_DIR}/fx .
start:
	${OUTPUT_DIR}/fx
list:
	@${OUTPUT_DIR}/fx list
clean:
	rm -rf ${OUTPUT_DIR}
.PHONY: test build start list clean
