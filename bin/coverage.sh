#!/usr/bin/env bash

set -e
echo "" > coverage.txt

for d in ./image; do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

#run service test at once
go test -race -coverprofile=profile.out -covermode=atomic $d ./api/service/*_test.go
if [ -f profile.out ]; then
    cat profile.out >> coverage.txt
    rm profile.out
fi
