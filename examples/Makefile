build:
	cd .. && make install-deps && make build
start-server: build
	sudo ../build/fx server stop
	sudo ../build/fx server start
deploy-functions: start-server
	../build/fx up functions/func.go
test: deploy-functions
	echo "Testing done"
