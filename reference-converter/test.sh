#!/bin/bash

set -euo pipefail

mkdir -p $(dirname "$1")
CGO_ENABLED=1 GOOS=linux go test -v -count 1 -race ./... 2>&1 | go-junit-report -iocopy -set-exit-code -out "$1"
