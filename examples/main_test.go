package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tommy351/goldga"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "examples")
}

var _ = Describe("Examples", func() {
	It("string", func() {
		Expect("abc").To(goldga.Match())
	})

	It("bool", func() {
		Expect(true).To(goldga.Match())
	})

	It("map", func() {
		Expect(map[string]interface{}{
			"a": "str",
			"b": true,
			"c": 123,
			"d": 3.14,
			"e": []string{"a", "b", "c"},
		}).To(goldga.Match())
	})

	It("multiple gold files in the same test", func() {
		Expect("foo").To(goldga.Match(goldga.WithDescription("first gold file")))
		Expect("bar").To(goldga.Match(goldga.WithDescription("second gold file")))
		Expect("foobar").To(goldga.Match(goldga.WithDescription("third gold file")))
	})
})
