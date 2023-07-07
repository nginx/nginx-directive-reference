ARG GO_VERSION

FROM docker.io/library/golang:${GO_VERSION} as builder
ARG GO_JUNIT_REPORT_VERSION
ARG GOLANGCI_LINT_VERSION

WORKDIR /app
RUN go install github.com/jstemmer/go-junit-report/v2@${GO_JUNIT_REPORT_VERSION}
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCI_LINT_VERSION}
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /reference-converter .

FROM scratch

COPY --from=builder /reference-converter /reference-converter
ENTRYPOINT ["/reference-converter"]
