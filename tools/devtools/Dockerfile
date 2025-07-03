ARG GO_VERSION=1.24
ARG NODE_VERSION=18
ARG BASE_IMG=docker.io/library/node:${NODE_VERSION}-bullseye

FROM docker.io/library/golang:${GO_VERSION}-bullseye AS golang
ARG GO_JUNIT_REPORT_VERSION=latest
ARG GOPLS_VERSION=latest
ARG DELVE_VERSION=latest
ARG GOLANGCI_LINT_VERSION=1.64.8

RUN go install github.com/jstemmer/go-junit-report/v2@${GO_JUNIT_REPORT_VERSION} \
    && go install -v golang.org/x/tools/gopls@${GOPLS_VERSION} \
    && go install github.com/go-delve/delve/cmd/dlv@${DELVE_VERSION} \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCI_LINT_VERSION}

FROM $BASE_IMG

# copy golang tools in
COPY --from=golang /go/bin/* /go/bin/
COPY --from=golang /usr/local/go /usr/local/go
ENV PATH=/go/bin:/usr/local/go/bin:${PATH} \
    CGO_ENABLED=0

RUN apt-get update -y \
    && apt-get install -y xz-utils curl \
    # create default location for devtools scripts
    && mkdir /opt/devtools && chmod -R 777 /opt/devtools \
    # tell git it's ok if the permissions are screwy
    && git config --system --add safe.directory /mnt/src \
    # clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/* /tmp/*

USER node

CMD ["/bin/bash"]

