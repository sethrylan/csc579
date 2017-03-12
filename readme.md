
Queue Simulation (P1)
===================

#### Prerequisites
Install Go runtime: https://golang.org/dl/

#### Install Build Packages
GOPATH=`pwd` go get -u github.com/golang/lint/golint

#### Build/Test/Execute
make
make test
./bin/qsim 0.5 10 1000 0

#### Documentation
make doc
