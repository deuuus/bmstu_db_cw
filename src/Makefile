.PHONY: build
build:
	go build -v ./cmd/apiserver

.PHONY: test
test:
	go test -v -race -timeout 30s ./...

db:
	psql -U polina -d med -h localhost -W

.DEFAULT_GOAL := build