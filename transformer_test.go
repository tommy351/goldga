package goldga

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NopTransformer", func() {
	It("should return the same value", func() {
		t := &NopTransformer{}
		v := struct{}{}
		Expect(t.Transform(v)).To(BeIdenticalTo(v))
	})
})
