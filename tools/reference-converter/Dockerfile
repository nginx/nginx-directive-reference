FROM docker.io/library/golang:1.19 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /reference-converter .
ENTRYPOINT ["/reference-converter"]

