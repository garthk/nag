source := $(shell find . -name \*.go)

.PHONY: deps all install test clean

nag: $(source)
	go build

all: deps nag test

deps:
	go get -v ./...

install: nag
	go install

lint:
	go vet ./...

test:
	go test -v -covermode=count -coverprofile=coverage.out ./naglib

clean:
	go clean
