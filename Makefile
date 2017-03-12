.PHONY: build doc fmt lint run test vet

GOPATH := ${PWD}
export GOPATH

default: build

build: vet
	go build -v -o ./bin/qsim ./src/main

doc:
	godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./src/...

# https://github.com/golang/lint
lint:
	./bin/golint ./...

run: build
	./bin/qsim

test:
	go test -cover -v ./test/...

vet:
	go vet ./src/...

tar:
	tar czvf p1.tar.gz --exclude=".DS_Store" Makefile readme.md ./src ./test ./paper
