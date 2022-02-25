package goldga

import (
	"bytes"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/spf13/afero"
)

var _ = Describe("Matcher", func() {
	var (
		matcher  *Matcher
		mockCtrl *gomock.Controller
		storage  *MockStorage
		actual   interface{}
		success  bool
		err      error
	)

	BeforeEach(func() {
		matcher = Match()
		mockCtrl = gomock.NewController(GinkgoT())
		storage = NewMockStorage(mockCtrl)
		matcher.Storage = storage
		actual = map[string]interface{}{
			"a": "str",
			"b": true,
			"c": 3.14,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	JustBeforeEach(func() {
		success, err = matcher.Match(actual)
	})

	testSucceed := func() {
		It("should succeed", func() {
			Expect(success).To(BeTrue())
		})

		It("should not return error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	}

	testFail := func() {
		It("should fail", func() {
			Expect(success).To(BeFalse())
		})

		It("should not return error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	}

	testError := func(expected types.GomegaMatcher) {
		It("should fail", func() {
			Expect(success).To(BeFalse())
		})

		It("should return error", func() {
			Expect(err).To(expected)
		})
	}

	getFileContent := func() []byte {
		var buf bytes.Buffer
		Expect(matcher.Serializer.Serialize(&buf, actual)).To(Succeed())

		return buf.Bytes()
	}

	testUpdateFile := func() {
		When("write success", func() {
			BeforeEach(func() {
				storage.EXPECT().Write(getFileContent()).Return(nil)
			})

			testSucceed()
		})

		When("write failed", func() {
			BeforeEach(func() {
				storage.EXPECT().Write(getFileContent()).Return(errors.New("error"))
			})

			testError(HaveOccurred())
		})
	}

	When("golden file exists", func() {
		When("UpdateFile = true", func() {
			BeforeEach(func() {
				matcher.UpdateFile = true
			})

			testUpdateFile()
		})

		When("match", func() {
			BeforeEach(func() {
				storage.EXPECT().Read().Return(getFileContent(), nil)
			})

			testSucceed()
		})

		When("not match", func() {
			BeforeEach(func() {
				storage.EXPECT().Read().Return([]byte{}, nil)
			})

			testFail()

			Context("failure message", func() {
				BeforeEach(func() {
					storage.EXPECT().Read().Return([]byte{}, nil)
				})

				It("positive", func() {
					Expect(matcher.FailureMessage(actual)).To(HavePrefix("Expected to match the golden file"))
				})

				It("negative", func() {
					Expect(matcher.NegatedFailureMessage(actual)).To(HavePrefix("Expected not to match the golden file"))
				})
			})
		})
	})

	When("golden file does not exist", func() {
		BeforeEach(func() {
			storage.EXPECT().Read().Return(nil, afero.ErrFileNotFound)
		})

		testUpdateFile()
	})

	When("failed to read golden file", func() {
		BeforeEach(func() {
			storage.EXPECT().Read().Return(nil, errors.New("error"))
		})

		testError(HaveOccurred())
	})
})

var _ = Describe("Options", func() {
	Describe("WithDescription", func() {
		It("should append a description to the test name, allowing multiple gold files per test", func() {
			Expect("foo").To(Match(WithDescription("First Gold File")))
			Expect("bar").To(Match(WithDescription("Second Gold File")))
			Expect("foobar").To(Match(WithDescription("Third Gold File")))
		})
	})
})
