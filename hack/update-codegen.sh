#!/bin/bash

set -euo pipefail

export GO111MODULE="on"

go get github.com/golang/mock/mockgen@latest
go generate ./...
