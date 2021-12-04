VERSION	:= $(shell git describe  --always --long)

.PHONY: clean build all

default: all

help:   ## show this help
	@echo 'usage: make [target] ...'
	@echo ''
	@echo 'targets:'
	@egrep '^(.+)\:\ .*##\ (.+)' ${MAKEFILE_LIST} | sed 's/:.*##/#/' | column -t -c 2 -s '#'

clean:  ## clean go environment for this package
	go clean

build:  ## build this package
	go build -ldflags="-X 'main.version=$(VERSION)'"

all:  ## clean and build this package (default)
	$(info Building $(VERSION))
	make clean
	make build
