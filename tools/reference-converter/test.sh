#!/bin/bash

set -euo pipefail

go test -v -count 1 -race ./... 2>&1 | go-junit-report -iocopy -set-exit-code -out "$1"
