SHELL=/bin/bash

.PHONY: run fmt

run:
	@rm -rf examples/destination/*
	@go run cmd/main.go

fmt:
	@go fmt ./...

go.sum: go.mod
	@go mod tidy