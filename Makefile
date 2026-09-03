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

smoke: build
	./$(BINARY) compile examples/echo-agent.yaml -o smoke.aegis
	./$(BINARY) verify smoke.aegis

clean:
	rm -f $(BINARY) $(BINARY).exe smoke.aegis run.atlog proof.json
	rm -rf .aegis dist
