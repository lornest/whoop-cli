.PHONY: build test lint clean run install

export GOROOT := /opt/homebrew/Cellar/go/1.25.6/libexec

build:
	go build -o whoop-cli ./cmd/whoop

install:
	go install ./cmd/whoop

test:
	go test ./... -v

lint:
	golangci-lint run

clean:
	rm -f whoop-cli
	go clean -testcache

run: build
	./whoop-cli

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
