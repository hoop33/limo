DIST_DIR ?= ./dist
PROG_NAME ?= limo
PACKAGES = $(shell go list ./... | grep -v /vendor/)

default: build

install: build
	go install

build: check
	go build

dist: macos linux windows

macos:
	GOOS=darwin go build -o $(DIST_DIR)/macos/$(PROG_NAME)

linux:
	GOOS=linux go build -o $(DIST_DIR)/linux/$(PROG_NAME)

windows:
	GOOS=windows go build -o $(DIST_DIR)/windows/$(PROG_NAME).exe

check: vet lint errcheck interfacer test

vet:
	go vet $(PACKAGES)

lint:
	golint -set_exit_status $(PACKAGES)

errcheck:
	errcheck $(PACKAGES)

interfacer:
	interfacer $(PACKAGES)

test:
	go test -cover $(PACKAGES)

clean:
	rm -rf dist/*

deps:
	go get -u github.com/FiloSottile/gvt
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/mvdan/interfacer/cmd/interfacer
