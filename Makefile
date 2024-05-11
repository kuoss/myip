GOLANGCI_LINT_VER := v1.58.1
GO_LICENSES_VER := v1.6.0

.PHONY: run
run:
	go install github.com/cosmtrek/air@latest || true
	air

.PHONY: test
test:
	go test -v ./... -race -failfast

.PHONY: cover
cover:
	go test -coverprofile=cover.out ./...
	go tool cover -func=cover.out
	go tool cover -func=cover.out | grep ^total: | grep 100.0%

.PHONY: lint
lint:
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VER) || true
	golangci-lint run

.PHONY: licenses
licenses:
	go install -v github.com/google/go-licenses@$(GO_LICENSES_VER) || true
	go-licenses check .

.PHONY: checks
checks: test cover lint licenses
