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

check: vet lint errcheck interfacer aligncheck structcheck varcheck unconvert gosimple staticcheck unused test

vet:
	go vet $(PACKAGES)

lint:
	golint -set_exit_status $(PACKAGES)

errcheck:
	errcheck $(PACKAGES)

interfacer:
	interfacer $(PACKAGES)

aligncheck:
	aligncheck $(PACKAGES)

structcheck:
	structcheck $(PACKAGES)

varcheck:
	varcheck $(PACKAGES)

unconvert:
	unconvert -v $(PACKAGES)

gosimple:
	gosimple $(PACKAGES)

staticcheck:
	staticcheck $(PACKAGES)

unused:
	unused $(PACKAGES)

test:
	go test -cover $(PACKAGES)

coverage:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

clean:
	go clean && rm -rf dist/*

deps:
	go get -u github.com/FiloSottile/gvt
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/mdempsky/unconvert
	go get -u github.com/mvdan/interfacer/cmd/interfacer
	go get -u github.com/opennota/check/...
	go get -u github.com/yosssi/goat/...
	go get -u honnef.co/go/tools/...
