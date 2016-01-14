source := $(shell find . -name \*.go)

.PHONY: default deps all install test clean

default: nag test

nag: $(source)
	go build

all: deps nag test

deps:
	go get github.com/stretchr/testify
	go get -v ./...

install: nag
	go install

lint:
	go vet ./...

test:
	go test -v -covermode=count -coverprofile=coverage.out ./naglib

clean:
	go clean

format:
	go fmt ./...
