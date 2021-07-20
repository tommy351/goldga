#!/bin/bash

set -euo pipefail

export GO111MODULE="on"

go get github.com/golang/mock/mockgen@v1.6.0
go generate ./...
go mod tidy
