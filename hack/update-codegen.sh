#!/bin/bash

set -euo pipefail

export GO111MODULE="on"
export PATH="$(go env GOPATH)/bin;$PATH"

go get github.com/golang/mock/mockgen@latest
go generate ./...
