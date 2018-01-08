BINARY ?= aegis

.PHONY: all build test race vet lint fmt check clean smoke

all: check build

build:
	go build -trimpath -o $(BINARY) ./cmd/aegis

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

fmt:
	gofmt -w $$(find . -name '*.go')

lint:
	golangci-lint run

check: vet race

