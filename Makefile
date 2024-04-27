init:
	go install github.com/cosmtrek/air@latest

run:
	air

test:
	go test -v -race -failfast ./...

checks:
	sh hack/checks.sh

cover:
	sh hack/cover.sh

lint:
	sh hack/lint.sh

licenses:
	sh hack/licenses.sh

