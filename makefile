packages = . ./cmd ./config ./model ./output ./service

default: check

all: check build

build: osx linux windows

osx:
	GOOS=darwin go build -o ./dist/osx/limo

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
