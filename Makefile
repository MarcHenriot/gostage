SHELL=/bin/bash

.PHONY: server demo fmt ui

run:
	@go run cmd/main.go

demo:
	@rm -rf examples/destination/*
	@go run cmd/main.go

fmt:
	@go fmt ./...

go.sum: go.mod
	@go mod tidy