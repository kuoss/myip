#!/bin/bash
cd $(dirname $0)/..

set -xeuo pipefail
go mod tidy
go fmt ./...
sh hack/misspell.sh
go vet ./...
which goimports || go install golang.org/x/tools/cmd/goimports@latest
goimports -local -v -w .
sh hack/staticcheck.sh
which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
golangci-lint run --timeout 5m
sh hack/test-cover.sh
sh hack/go-licenses.sh
