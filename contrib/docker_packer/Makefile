GOBIN ?= ./build
GIT_VERSION := $(shell git describe --tags)
VERSION ?= $(GIT_VERSION)

REPO ?= "metrue/fx-docker"
TAG ?= "latest"

build:
	CGO_ENABLED=0 go build -ldflags "-X main.Version=$(VERSION)" -v -o $(GOBIN)/docker_packer main.go
linux-build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=$(VERSION)" -v -o $(GOBIN)/docker_packer main.go
docker-build:
	docker build -t ${REPO}:${TAG} .
docker-publish:
	docker push ${REPO}:${TAG}
test:
	docker run -v /var/run/docker.sock:/var/run/docker.sock ${REPO}:${TAG} docker_packer  'eyJEb2NrZXJmaWxlIjoiRlJPTSBtZXRydWUvZngtbm9kZS1iYXNlXG5cbkNPUFkgLiAuXG5FWFBPU0UgMzAwMFxuQ01EIFtcIm5vZGVcIiwgXCJhcHAuanNcIl1cbiIsImFwcC5qcyI6ImNvbnN0IEtvYSA9IHJlcXVpcmUoJ2tvYScpO1xuY29uc3QgYm9keVBhcnNlciA9IHJlcXVpcmUoJ2tvYS1ib2R5cGFyc2VyJyk7XG5jb25zdCBmeCA9IHJlcXVpcmUoJy4vZngnKTtcblxuY29uc3QgYXBwID0gbmV3IEtvYSgpO1xuYXBwLnVzZShib2R5UGFyc2VyKCkpO1xuYXBwLnVzZShmeCk7XG5cbmFwcC5saXN0ZW4oMzAwMCk7XG4iLCJmeC5qcyI6IlxubW9kdWxlLmV4cG9ydHMgPSAoY3R4KSA9XHUwMDNlIHtcblx0Y3R4LmJvZHkgPSAnaGVsbG8gd29ybGQnXG59XG4ifQ==' app-hello
	docker run --rm -d -p 3000:3000 --name test-app-hello-container app-hello
	sleep 2
	curl 127.0.0.1:3000
	docker stop test-app-hello-container
