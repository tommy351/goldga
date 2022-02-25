#!/bin/bash

set -euo pipefail

go run github.com/onsi/ginkgo/v2/ginkgo -r --randomize-all --randomize-suites --fail-on-pending --cover --trace --race --progress --junit-report=junit.xml $@
$(dirname "${BASH_SOURCE[0]}")/collect-coverage.sh
