#!/bin/bash

DISABLE_AUTH="True" go test -cover ./...
golint $(go list ./... | grep -v /vendor/)
errcheck -ignoretests ./...
golangci-lint run --tests=False --exclude=auth ./...