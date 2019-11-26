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
	# TODO enable remote test
	echo 'run testing on remote host'
	DOCKER_REMOTE_HOST_ADDR=${REMOTE_HOST_ADDR} DOCKER_REMOTE_HOST_USER=${REMOTE_HOST_USER} DOCKER_REMOTE_HOST_PASSWORD=${REMOTE_HOST_PASSWORD} ./scripts/test_cli.sh

http-test:
	./scripts/http_test.sh

zip:
	zip -r images.zip images/
.PHONY: test build start list clean generate

start_docker_infra:
	docker build -t fx-docker-infra -f test/Dockerfile ./test
	docker run --rm --name fx-docker-infra -p 2222:22 -v /var/run/docker.sock:/var/run/docker.sock -d fx-docker-infra

test_docker_infra:
	SSH_PORT=2222 SSH_KEY_FILE=./test/id_rsa ./build/fx infra create --name docker-local -t docker --host root@127.0.0.1

stop_docker_infra:
	docker stop fx-docker-infra

start_k3s_infra:
	multipass launch --name k3s-master --cpus 1 --mem 512M --disk 3G  --cloud-init ./test/k3s/ssh-cloud-init.yaml
	multipass launch --name k3s-worker1 --cpus 1 --mem 512M --disk 3G  --cloud-init ./test/k3s/ssh-cloud-init.yaml
	multipass launch --name k3s-worker2 --cpus 1 --mem 512M --disk 3G  --cloud-init ./test/k3s/ssh-cloud-init.yaml

test_k3s_infra:
	./scripts/test_k3s_infra.sh

stop_k3s_infra:
	multipass delete k3s-master
	multipass delete k3s-worker1
	multipass delete k3s-worker2
	multipass purge
