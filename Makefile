MAKEFLAGS += --silent
XK6_VERSION ?= v1.4.13
K6_VERSION ?= $(shell go list -m -f '{{.Version}}' go.k6.io/k6)

all: clean format test build

## help: Prints a list of available build targets.
help:
	echo "Usage: make <OPTIONS> ... <TARGETS>"
	echo ""
	echo "Available targets are:"
	echo ''
	sed -n 's/^##//p' ${PWD}/Makefile | column -t -s ':' | sed -e 's/^/ /'
	echo
	echo "Targets run by default are: `sed -n 's/^all: //p' ./Makefile | sed -e 's/ /, /g' | sed -e 's/\(.*\), /\1, and /'`"

## clean: Removes any previously created build artifacts.
clean:
	rm -f ./k6

## build: Builds a custom 'k6' with the local extension. 
build:
	go run go.k6.io/xk6/cmd/xk6@$(XK6_VERSION) build $(K6_VERSION) --with $(shell go list -m)=.

## format: Applies Go formatting to code.
format:
	go fmt ./...

## test: Executes any unit tests.
test:
	go test -cover -race .
	cd examples/backend && go test -cover -race ./...

.PHONY: build clean format help test
