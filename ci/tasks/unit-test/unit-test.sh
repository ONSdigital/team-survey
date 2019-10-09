#!/usr/bin/env bash
set -xeuo pipefail

export GOPATH="$PWD/gopath"
export PATH="$PWD/gopath/bin:$PATH"

# Move the project to the GOPATH
mkdir -p "$GOPATH/src/armakuni"
cp -rf team-survey "$GOPATH/src/armakuni/"

# Change directory to the project in the GOPATH and run the tests
cd "$GOPATH/src/armakuni/team-survey"

go test -v ./...