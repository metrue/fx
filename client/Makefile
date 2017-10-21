OUTPUT_DIR=./build

test:
	${OUTPUT_DIR}/fx up index.js
	# docker ps
build:
	go build -o ${OUTPUT_DIR}/fx .
start:
	${OUTPUT_DIR}/fx
clean:
	rm -rf ${OUTPUT_DIR}
.PHONY: test build start clean
