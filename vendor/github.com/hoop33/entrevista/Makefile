packages = .

default: build

build: check
	go build

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

deps:
	go get -u github.com/FiloSottile/gvt
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	go get -u github.com/mvdan/interfacer/cmd/interfacer
