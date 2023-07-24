#!/usr/bin/make -f
DEVTOOLS_IMG?=ampex/devtools
ROOT_DIR:=$(shell git rev-parse --show-toplevel)
DEVTOOLS_DIR:=$(realpath $(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

# - if we're running within the devtools container (eg. when running within a
#   CI pipeline), then we don't need to issue "docker" commands
ifeq ("$(wildcard /opt/devtools)", "")
.run: .run-outside-container
else
.run: .run-inside-container
endif

# - for more info about ssh agent forwarding in docker-for-mac, see:
#   https://docs.docker.com/desktop/mac/networking/#ssh-agent-forwarding
OS := $(shell uname -s)
ifeq ("$(OS)", "Darwin")
	LOCAL_AUTH_SOCK := /run/host-services/ssh-auth.sock
else
	LOCAL_AUTH_SOCK := $(SSH_AUTH_SOCK)
endif

devtools-image: ## build a image containing all our tooling
	docker build -t $(DEVTOOLS_IMG) $(ROOT_DIR)/tools/devtools

devtools-versions: ## print the versions of software in the devtools image
	@grep ARG $(DEVTOOLS_DIR)/Dockerfile | grep _VERSION= | sed -e 's/ARG //'

devtools-shell: ## explore / debug the devtools image
	@$(MAKE) .run args="bash"

#######################################
## private targets
#######################################

.run-outside-container:
	@docker run --rm \
		$(shell tty > /dev/null && echo "-it") \
		--net=host \
		-v $(ROOT_DIR):/mnt/src \
		-v $(LOCAL_AUTH_SOCK):/ssh-auth.sock \
		-e SSH_AUTH_SOCK=/ssh-auth.sock \
		-w $(subst $(ROOT_DIR), /mnt/src, $(PWD)) \
		$(shell for item in $$(env | grep -E '^(TEST|CI|GITHUB)_' | cut -d '=' -f 1); do echo "-e $$item"; done) \
		-e TZ=$(shell echo "$$(readlink /etc/localtime)" | sed 's#.*zoneinfo/##') \
		-e GOOS="$(shell uname -s | tr '[:upper:]' '[:lower:]')" \
		$(DEVTOOLS_IMG) $(args)

.run-inside-container:
	$(args)
