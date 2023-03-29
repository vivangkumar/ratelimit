.PHONY: test build lint fmt

goimports := go run golang.org/x/tools/cmd/goimports
lint := go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

build:
	go build ./...

fmt:
	gofmt -s -w .
	$(lint) run --fast --fix

lint:
	$(lint) run

test:
	go test -v -race ./...
