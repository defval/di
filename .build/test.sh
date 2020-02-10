#!/usr/bin/env bash

set -e

if [[ -f coverage.txt ]]; then
  rm coverage.txt
fi

for d in $(go list ./... | grep -v ditest); do
  go test -coverprofile=profile.out -coverpkg=./... -covermode=atomic "$d"
  if [[ -f profile.out ]]; then
    cat profile.out >>coverage.txt
    rm profile.out
  fi
done
