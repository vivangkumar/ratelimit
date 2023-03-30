.PHONY: test bench build lint fmt

lint := go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2

install-tools:
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2

build:
	@echo "Building..."
	go build ./...

fmt:
	@echo "Formatting..."
	gofmt -s -w .
	$(lint) run --fast --fix

lint:
	@echo "Linting..."
	$(lint) run

test:
	@echo "Running tests..."
	go test -v -race ./...

bench:
	@echo "Running benchmarks..."
	go test -bench=. -count 5 -run=^# ./...
