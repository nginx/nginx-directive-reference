ROOT_DIR:=$(shell git rev-parse --show-toplevel)
include $(ROOT_DIR)/tools/devtools/Makefile

.PHONY: build test lint

build:
	@$(MAKE) .run args="go build -v -o dist/ ."

test:
	@$(MAKE) .run args="./test.sh results/test.xml"

lint:
	@$(MAKE) .run args="golangci-lint run -v"
