package goldga

import (
	"path/filepath"
	"strings"

	"github.com/onsi/ginkgo"
)

func getCurrentGinkgoTestDescription() ginkgo.GinkgoTestDescription {
	return ginkgo.CurrentGinkgoTestDescription()
}

func getGinkgoPath() string {
	desc := getCurrentGinkgoTestDescription()
	path := desc.FileName

	if path == "" {
		panic("current file name is empty")
	}

	name := filepath.Base(desc.FileName)

	if ext := filepath.Ext(name); ext != "" {
		name = strings.TrimSuffix(name, ext)
		name = strings.TrimSuffix(name, "_test")
	}

	return filepath.Join("testdata", name+".golden")
}

func getGinkgoTestName() string {
	testName := getCurrentGinkgoTestDescription().FullTestText

	if testName == "" {
		panic("current test name is empty")
	}

	return testName
}
