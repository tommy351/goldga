package goldga

import (
	"bytes"
	"encoding/json"

	"github.com/BurntSushi/toml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

// nolint: gochecknoglobals
var serializerTestData = map[string]interface{}{
	"a": "str",
	"b": true,
	"c": 3.14,
}

func testSerializer(s Serializer) []byte {
	var buf bytes.Buffer
	Expect(s.Serialize(&buf, serializerTestData)).To(Succeed())

	return buf.Bytes()
}

var _ = Describe("DumpSerializer", func() {
	It("works", func() {
		serializer := &DumpSerializer{
			Config: newDefaultDumpConfig(),
		}
		expected := newDefaultDumpConfig().Sdump(serializerTestData)
		Expect(testSerializer(serializer)).To(Equal([]byte(expected)))
	})
})

var _ = Describe("YAMLSerializer", func() {
	It("works", func() {
		expected, err := yaml.Marshal(serializerTestData)
		Expect(err).NotTo(HaveOccurred())
		Expect(testSerializer(&YAMLSerializer{})).To(MatchYAML(expected))
	})
})

var _ = Describe("JSONSerializer", func() {
	It("works", func() {
		expected, err := json.Marshal(serializerTestData)
		Expect(err).NotTo(HaveOccurred())
		Expect(testSerializer(&JSONSerializer{})).To(MatchJSON(expected))
	})
})

var _ = Describe("TOMLSerializer", func() {
	marshalTOML := func(input interface{}) []byte {
		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		Expect(enc.Encode(input)).To(Succeed())

		return buf.Bytes()
	}

	It("works", func() {
		expected := marshalTOML(serializerTestData)
		Expect(testSerializer(&TOMLSerializer{})).To(Equal(expected))
	})
})
