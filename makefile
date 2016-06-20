default: check

all: check build

build: osx linux windows

osx:
	GOOS=darwin go build -o ./dist/osx/limo

linux:
	GOOS=linux go build -o ./dist/linux/limo

windows:
	GOOS=windows go build -o ./dist/windows/limo.exe

check: lint errcheck test

lint:
	golint --set_exit_status .
	golint --set_exit_status ./cmd
	golint --set_exit_status ./config
	golint --set_exit_status ./output
	golint --set_exit_status ./service

test:
	go test . ./cmd ./config ./output ./service

errcheck:
	errcheck ./cmd
	errcheck ./config
	errcheck ./output
	errcheck ./service

clean:
	rm -rf dist/*
