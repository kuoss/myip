GOLANGCI_LINT_VER := v1.58.1

run:
	go install github.com/cosmtrek/air@latest || true
	air

lint:
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VER) || true
	golangci-lint run

test:
	go test -v -race -failfast ./...

cover:
	sh hack/cover.sh

licenses:
	sh hack/licenses.sh

checks: lint cover licenses
