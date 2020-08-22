.PHONY: build
build:

		go build -v ./cmd/apiserver


.PHONY: build-linux
build-linux:

		env GOOS=linux go build -v ./cmd/apiserver

.PHONY: test
test:

		go test -v -race -timeout 30s ./...


.DEFAULT_GOAL := build