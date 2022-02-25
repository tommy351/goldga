# Goldga

[![GoDoc](https://godoc.org/github.com/tommy351/goldga?status.svg)](https://godoc.org/github.com/tommy351/goldga) [![GitHub release](https://img.shields.io/github/release/tommy351/goldga.svg)](https://github.com/tommy351/goldga/releases) [![CircleCI](https://circleci.com/gh/tommy351/goldga/tree/master.svg?style=svg)](https://circleci.com/gh/tommy351/goldga/tree/master) [![codecov](https://codecov.io/gh/tommy351/goldga/branch/master/graph/badge.svg)](https://codecov.io/gh/tommy351/goldga)

A golden file testing (snapshot testing) library for [gomega](http://onsi.github.io/gomega/).

## Installation

```sh
go get github.com/tommy351/goldga
```

## Usage

```go
import (
  . "github.com/onsi/ginkgo/v2"
  . "github.com/onsi/gomega"
  "github.com/tommy351/goldga"
)

Describe("Example", func() {
  It("works", func() {
    Expect("foobar").To(goldga.Match())
  })
})
```

See [examples](examples) folder for more examples.
