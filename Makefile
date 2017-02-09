packages = . ./cmd ./config ./model ./output ./service

default: check

install: build
	go install

all: check build

build:
	go build

dist: macos linux windows

macos:
	GOOS=darwin go build -o ./dist/macos/limo

linux:
	GOOS=linux go build -o ./dist/linux/limo

windows:
	GOOS=windows go build -o ./dist/windows/limo.exe

check: vet lint errcheck test

vet:
	go vet $(packages)

lint:
	$(foreach package,$(packages),golint --set_exit_status $(package);)

test:
	go test -cover $(packages)

errcheck:
	errcheck $(packages)

clean:
	rm -rf dist/*

get-deps:
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
