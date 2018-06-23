#!/usr/bin/env bash

echo "mode: atomic\n" > coverage.txt

# TODO it's very weird, tests of api/service always fails on CirCleCI environment,
# but works on local
for d in `go list ./... | grep -v 'third_party\|examples\|assets'`; do
  echo $d
  go test -race -coverprofile=profile.out -covermode=atomic $d
  if [ -f profile.out ]; then
      cat profile.out | grep -v "mode: atomic" >> coverage.txt
      rm profile.out
  fi
done
