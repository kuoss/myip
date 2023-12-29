init:
	go install github.com/cosmtrek/air@latest

run:
	air

checks: misspell gocyclo

misspell:
	sh hack/misspell.sh

gocyclo:
	sh hack/gocyclo.sh
