SHELL := /bin/bash

.PHONY: module-sync build run tidy

module-sync:
	./scripts/update-module.sh

build:
	go build -o pokemon-fetch-cli ./...

run: build
	./pokemon-fetch-cli

tidy:
	go mod tidy
