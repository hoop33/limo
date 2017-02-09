packages = . ./cmd ./config ./model ./output ./service

default: build

install: build
	go install

build: check
	go build

dist: macos linux windows

macos:
	GOOS=darwin go build -o ./dist/macos/limo

linux:
	GOOS=linux go build -o ./dist/linux/limo

windows:
	GOOS=windows go build -o ./dist/windows/limo.exe

check: vet lint errcheck interfacer test

vet:
	go vet $(packages)

lint:
	$(foreach package,$(packages),golint -set_exit_status $(package);)

errcheck:
	errcheck $(packages)

interfacer:
	interfacer $(packages)

test:
	go test -cover $(packages)

clean:
	rm -rf dist/*

get-deps:
	go get -u github.com/FiloSottile/gvt
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/mvdan/interfacer/cmd/interfacer
