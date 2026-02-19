SHELL := /bin/sh

.PHONY: test lint generate check check-generated

test:
	go test ./...

lint:
	go vet ./...

generate:
	go run ./cmd/generator

check: lint test

# Requires network access to Rapid7 docs.
check-generated:
	go run ./cmd/generator
	git diff --exit-code -- pkg/msfrpc/generated
