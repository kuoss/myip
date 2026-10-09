.PHONY: test
test:
	go test -v ./... -race -failfast

.PHONY: cover
cover:
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out
	go tool cover -func=cover.out | grep ^total: | grep 100.0%

## golangci-lint with its default linters. Built with this module's Go version so it can load the code.
GOLANGCI_LINT_VERSION ?= latest

.PHONY: lint
lint:
	GOTOOLCHAIN=$$(go env GOVERSION) go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run

.PHONY: checks
checks: test cover lint
