#!/usr/bin/env bash

echo "mode: atomic\n" > coverage.txt

for d in `go list ./... | grep -v 'mocks\|images\|examples\|assets'`; do
  echo $d
  go test -race -coverprofile=profile.out -covermode=atomic $d
  if [ -f profile.out ]; then
      cat profile.out | grep -v "mode: atomic" >> coverage.txt
      rm profile.out
  fi
done
