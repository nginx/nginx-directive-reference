ARG GO_VERSION=1.19-bullseye
ARG BASE_IMG=mcr.microsoft.com/devcontainers/typescript-node:0-18-bullseye

FROM docker.io/library/golang:$GO_VERSION as golang
ARG GO_JUNIT_REPORT_VERSION=latest
ARG GOPLS_VERSION=latest
ARG GOLANGCI_LINT_VERSION=1.53.3

RUN go install github.com/jstemmer/go-junit-report/v2@${GO_JUNIT_REPORT_VERSION} \
    && go install -v golang.org/x/tools/gopls@latest \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCI_LINT_VERSION}

FROM $BASE_IMG
ARG DOCKER_VERSION=23.0.6

# copy golang tools in
COPY --from=golang /go/bin/* /go/bin/
COPY --from=golang /usr/local/go /usr/local/go
ENV PATH=/go/bin:/usr/local/go/bin:${PATH}

RUN apt-get update -y \
    && apt-get install -y xz-utils curl \
    # install docker
    && export ARCH=$(dpkg --print-architecture) \
    && export DOCKER_ARCH=$(case $ARCH in "arm64") echo "aarch64" ;; "amd64") echo "x86_64" ;; esac;) \
    && export DOCKER_URL=https://download.docker.com/linux/static/stable/${DOCKER_ARCH}/docker-${DOCKER_VERSION}.tgz \
    && curl -sSfL "$DOCKER_URL" -o /tmp/docker.tgz \
    && tar xzvf /tmp/docker.tgz --strip 1 -C /usr/local/bin docker/docker \
    # create default location for devtools scripts
    && mkdir /opt/devtools && chmod -R 777 /opt/devtools \
    # clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/* /tmp/*

CMD ["/bin/bash"]

