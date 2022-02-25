package goldga

import (
	"bytes"
	"encoding/json"

	"github.com/BurntSushi/toml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

// nolint: gochecknoglobals
var serializerTestData = map[string]interface{}{
	"a": "str",
	"b": true,
	"c": 3.14,
}

type stringer struct {
	Value string
}

func (s stringer) String() string {
	return s.Value
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

var _ = Describe("StringSerializer", func() {
	DescribeTable("Using default fallback serializer", func(input interface{}, expected string) {
		var buf bytes.Buffer
		serializer := &StringSerializer{}
		Expect(serializer.Serialize(&buf, input)).To(Succeed())

		Expect(buf.String()).To(Equal(expected))
	},
		Entry("string", "abc", "abc"),
		Entry("[]byte", []byte("abc"), "abc"),
		Entry("fmt.Stringer", stringer{Value: "abc"}, "abc"),
		Entry("int", 42, "(int) 42\n"),
	)

	Describe("Using custom fallback serializer", func() {
		It("should use custom fallback serializer on other types", func() {
			var buf bytes.Buffer
			serializer := &StringSerializer{
				FallbackSerializer: &JSONSerializer{},
			}
			Expect(serializer.Serialize(&buf, serializerTestData)).To(Succeed())

			expected, err := json.Marshal(serializerTestData)
			Expect(err).NotTo(HaveOccurred())
			Expect(buf.Bytes()).To(MatchJSON(expected))
		})
	})
})
