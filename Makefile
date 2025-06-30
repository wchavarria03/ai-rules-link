# Makefile for ai-rules-link

.PHONY: help test build lint

default: help

help:
	@echo "Available targets:"
	@echo "  test   Run all Go tests with coverage"
	@echo "  build  Build the ai-rules-link binary"
	@echo "  lint   Run gofmt and go vet on all Go files"

test:
	go test -v -cover ./... 

build:
	go build -o ai-rules-link . 

lint:
	gofmt -l -s .
	go vet ./... 