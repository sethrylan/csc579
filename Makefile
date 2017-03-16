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
	go fmt mm1k main

# https://github.com/golang/lint
lint:
	./bin/golint mm1k main

run: build
	./bin/qsim

test:
	go test -cover -v ./test/...

vet:
	go vet mm1k main

p1:
	tar czvf p1.tar.gz --exclude=".DS_Store" Makefile readme.md ./src ./test ./p1/paper/
