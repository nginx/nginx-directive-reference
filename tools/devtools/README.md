# devtools

This image wraps up all the development tools used locally or in the CI pipeline.

Inspired by:

- https://gitlab.com/f5/nginx/nginxazurelb/tools/nlb-devtools
- https://github.com/nginxinc/amptest/tree/main/docker/devtools

## Usage

Include the makefile in your app's folder, and call `make .run` with your desired command in `args`. Example:

```make
ROOT_DIR:=$(shell git rev-parse --show-toplevel)
include $(ROOT_DIR)/tools/devtools/Makefile

build:
	@$(MAKE) .run args="go build -v -o dist/ ."
```

If we're running outside devtools container, then we'll `docker run` the `args` with a lot of funky options to Just Work.

## Versions

Versions of common dev software are specified as docker build ARGs, so the dockerfile can be easily re-used as a vscode devcontainer without repeating version numbers in multiple files.

These versions will be shared by all apps using the devtools container, so be conservative with what is included.
