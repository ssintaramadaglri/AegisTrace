BINARY ?= aegis

.PHONY: all build test race vet lint fmt check clean smoke

all: check build
