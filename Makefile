init:
	go install github.com/cosmtrek/air@latest

run:
	air

checks:
	sh hack/checks.sh

misspell:
	sh hack/misspell.sh

gocyclo:
	sh hack/gocyclo.sh

staticcheck:
	sh hack/staticcheck.sh
