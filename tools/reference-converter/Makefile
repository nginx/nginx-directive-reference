# TODO: use the devtools image

ROOT_DIR:=$(shell git rev-parse --show-toplevel)
GO_VERSION=$(shell grep --max-count 1 'go ' go.mod |  cut -d' ' -f2)
GO_JUNIT_REPORT_VERSION?=latest
GOLANGCI_LINT_VERSION?=1.53.3

DOCKER_BUILD_ARGS := $(DOCKER_BUILD_ARGS) \
		--build-arg "GO_VERSION=$(GO_VERSION)" \
		--build-arg "GO_JUNIT_REPORT_VERSION=$(GO_JUNIT_REPORT_VERSION)" \
		--build-arg "GOLANGCI_LINT_VERSION=$(GOLANGCI_LINT_VERSION)"

.PHONY: build test lint

build:
	docker buildx build -t reference-converter:builder $(DOCKER_BUILD_ARGS) --target builder .
	docker buildx build -t reference-converter $(DOCKER_BUILD_ARGS) .

/results:
	@mkdir -p $(PWD)/results

test: /results
	docker run --rm -v "$(PWD)/results:/results" \
		reference-converter:builder ./test.sh /results/test.xml

lint:
	docker run --rm -v "$(PWD):/app" \
		-v "$(ROOT_DIR)/.golangci.yml:/app/.golangci.yml" \
		-w /app \
		reference-converter:builder golangci-lint run -v
