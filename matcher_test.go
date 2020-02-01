package goldga

import (
	"encoding/json"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
)

var _ = Describe("Matcher", func() {
	var (
		fs      afero.Fs
		matcher *Matcher
		actual  interface{}
		success bool
		err     error
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
		success, err = matcher.Match(actual)
	})

	getFile := func() *goldenFile {
		file, err := fs.Open(matcher.Path)
		Expect(err).NotTo(HaveOccurred())
		defer file.Close()

		gf, err := readGoldenFile(file)
		Expect(err).NotTo(HaveOccurred())

		return gf
	}

	resetFileTime := func(path string) {
		zero := time.Unix(0, 0)
		Expect(fs.Chtimes(path, zero, zero)).To(Succeed())
	}

	writeFile := func(gf *goldenFile) {
		file, err := fs.Create(matcher.Path)
		Expect(err).NotTo(HaveOccurred())

		defer resetFileTime(matcher.Path)
		defer file.Close()

		Expect(writeGoldenFile(file, gf)).To(Succeed())
	}

	genCorrectGoldenFile := func() *goldenFile {
		return &goldenFile{
			Version: goldenFileVersion,
			Snapshots: snapshotMap{
				getGinkgoTestName(): getDefaultDumpConfig().Sdump(actual),
			},
		}
	}

	testSucceed := func() {
		It("should succeed", func() {
			Expect(success).To(BeTrue())
		})

		It("should not return error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	}

	testFailed := func() {
		It("should fail", func() {
			Expect(success).To(BeFalse())
		})

		It("should not return the error", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should generate failure message", func() {
			Expect(matcher.FailureMessage(actual)).To(HavePrefix(fmt.Sprintf("Expected to match the golden file %q", matcher.Path)))
		})

		It("should generate negated failure message", func() {
			Expect(matcher.NegatedFailureMessage(actual)).To(HavePrefix(fmt.Sprintf("Expected not to match the golden file %q", matcher.Path)))
		})
	}

	testFileUnchanged := func() {
		It("should not update the golden file", func() {
			stat, err := fs.Stat(matcher.Path)
			Expect(err).NotTo(HaveOccurred())
			Expect(stat.ModTime()).To(Equal(time.Unix(0, 0)))
		})
	}

	When("golden file exists", func() {
		When("UpdateFile = true", func() {
			BeforeEach(func() {
				matcher.UpdateFile = true
			})
		})

		When("golden file match", func() {
			BeforeEach(func() {
				writeFile(genCorrectGoldenFile())
			})

			testSucceed()
			testFileUnchanged()
		})

		When("golden file not match", func() {
			BeforeEach(func() {
				writeFile(&goldenFile{
					Version: goldenFileVersion,
					Snapshots: snapshotMap{
						getGinkgoTestName(): "foo",
					},
				})
			})

			testFailed()
			testFileUnchanged()
		})
	})

	When("golden file does not exist", func() {
		testSucceed()

		It("should write a new golden file", func() {
			Expect(getFile()).To(Equal(genCorrectGoldenFile()))
		})

		When("Serializer is set", func() {
			BeforeEach(func() {
				matcher.Serializer = &JSONSerializer{}
			})

			testSucceed()

			It("should serialize into JSON format", func() {
				file := getFile()
				expected, err := json.Marshal(actual)
				Expect(err).NotTo(HaveOccurred())
				Expect(file.Snapshots).To(HaveKeyWithValue(getGinkgoTestName(), MatchJSON(expected)))
			})
		})
	})
})
