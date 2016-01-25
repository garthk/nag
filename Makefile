source := $(shell find . -name \*.go)

.PHONY: default deps all install test clean

default: nag test

nag: $(source)
	go build

all: deps nag test

deps:
	go get github.com/stretchr/testify
	go get github.com/cpuguy83/go-md2man/md2man
	go get -v ./...

install: nag
	go install

lint:
	go vet ./...

test:
	go test -covermode=count -coverprofile=naglib.out ./naglib
	go test -covermode=count -coverprofile=rfs.out ./pkg/readable-fs
	go test -covermode=count -coverprofile=userinfo.out ./pkg/user-info-shim
	go tool cover -func=naglib.out
	go tool cover -func=rfs.out
	go tool cover -func=userinfo.out

clean:
	go clean

format:
	go fmt ./...
