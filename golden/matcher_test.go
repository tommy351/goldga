package golden

import (
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Matcher", func() {
	var (
		fs      afero.Fs
		matcher *Matcher
		actual  interface{}
	)

	BeforeEach(func() {
		fs = afero.NewMemMapFs()
		matcher = Match()
		matcher.fs = fs
		actual = []interface{}{
			true,
			"str",
			123,
		}
	})

	JustBeforeEach(func() {
		Expect(actual).To(matcher)
	})

	getGoldenFile := func() *goldenFile {
		file, err := fs.Open(matcher.getPath())
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		gf, err := readGoldenFile(file)
		Expect(err).NotTo(HaveOccurred())

		return gf
	}

	When("golden file exists", func() {
	})

	When("golden file does not exist", func() {
		It("should write a new golden file", func() {
			Expect(getGoldenFile()).To(Equal(&goldenFile{
				Version: 1,
				Snapshots: snapshotMap{
					getGinkgoTestName(): spew.Sdump(actual),
				},
			}))
		})

		When("Serializer is set", func() {
			BeforeEach(func() {
				matcher.Serializer = &JSONSerializer{}
			})

			It("should serialize into JSON format", func() {
				file := getGoldenFile()
				expected, err := json.Marshal(actual)
				Expect(err).NotTo(HaveOccurred())
				Expect(file.Snapshots).To(HaveKeyWithValue(getGinkgoTestName(), MatchJSON(expected)))
			})
		})
	})
})
