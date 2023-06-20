FROM docker.io/library/golang:1.19 as builder
WORKDIR /app
RUN go install github.com/jstemmer/go-junit-report/v2@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /reference-converter .

FROM scratch

COPY --from=builder /reference-converter /reference-converter
ENTRYPOINT ["/reference-converter"]
