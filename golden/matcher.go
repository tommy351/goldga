package golden

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega/types"
	"github.com/spf13/afero"
)

// Default configs
const (
	DefaultFixtureDir     = "testdata"
	DefaultFileNamePrefix = ""
	DefaultFileNameSuffix = ".golden"
)

// nolint: gochecknoglobals
var defaultFs = afero.NewCacheOnReadFs(
	afero.NewOsFs(),
	afero.NewMemMapFs(),
	time.Minute,
)

func getGinkgoFileName() string {
	desc := ginkgo.CurrentGinkgoTestDescription()
	path := desc.FileName

	if path == "" {
		panic("current file name is empty")
	}

	name := filepath.Base(desc.FileName)

	if ext := filepath.Ext(name); ext != "" {
		name = strings.TrimSuffix(name, ext)
		name = strings.TrimSuffix(name, "_test")
	}

	return name
}

func getGinkgoTestName() string {
	testName := ginkgo.CurrentGinkgoTestDescription().FullTestText

	if testName == "" {
		panic("current test name is empty")
	}

	return testName
}

func getUpdateFile() bool {
	update, _ := strconv.ParseBool(os.Getenv("UPDATE_GOLDEN"))
	return update
}

func getColor() bool {
	// TODO: Detect if tty is colorable.
	return true
}

// Match succeeds if actual matches the golden file.
func Match() *Matcher {
	return &Matcher{
		FixtureDir:     DefaultFixtureDir,
		FileNamePrefix: DefaultFileNamePrefix,
		FileNameSuffix: DefaultFileNameSuffix,
		FileName:       getGinkgoFileName(),
		TestName:       getGinkgoTestName(),
		Serializer:     DefaultSerializer,
		Transformer:    DefaultTransformer,
		UpdateFile:     getUpdateFile(),
		Color:          getColor(),
		fs:             defaultFs,
	}
}

var _ types.GomegaMatcher = (*Matcher)(nil)

// Matcher implements GomegaMatcher.
type Matcher struct {
	// Path of the fixture directory.
	FixtureDir string

	// Name of the golden file.
	FileName string

	// Prefix of the file name.
	FileNamePrefix string

	// Suffix of the file name.
	FileNameSuffix string

	// Name of the test.
	TestName string

	Serializer Serializer

	Transformer Transformer

	// Display colored output.
	Color bool

	// Force update the golden file.
	UpdateFile bool

	fs afero.Fs
}

func (m *Matcher) getPath() string {
	return filepath.Join(m.FixtureDir, m.FileNamePrefix+m.FileName+m.FileNameSuffix)
}

// Match implements GomegaMatcher.
func (m *Matcher) Match(actual interface{}) (bool, error) {
	actualContent, err := m.getActualContent(actual)

	if err != nil {
		return false, err
	}

	expected, err := m.getExpectedContent()

	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return false, err
		}

		if err := m.writeFile(actualContent); err != nil {
			return false, err
		}

		return true, nil
	}

	return expected == actualContent, nil
}

func (m *Matcher) writeFile(actualContent string) error {
	path := m.getPath()

	if err := m.fs.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	file, err := m.fs.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		return err
	}

	defer file.Close()

	gf, err := readGoldenFile(file)

	if err != nil {
		return err
	}

	gf.Snapshots[m.TestName] = actualContent

	if err := writeGoldenFile(file, gf); err != nil {
		return err
	}

	return nil
}

func (m *Matcher) getMessage(actual interface{}, message string) string {
	expectedContent, err := m.getExpectedContent()

	if err != nil {
		panic(err)
	}

	actualContent, err := m.getActualContent(actual)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Expected %s match the golden file %q\n%s",
		message,
		m.getPath(),
		diffString(m.Color, expectedContent, actualContent))
}

func (m *Matcher) getExpectedContent() (string, error) {
	if m.UpdateFile {
		return "", os.ErrNotExist
	}

	file, err := m.fs.Open(m.getPath())

	if err != nil {
		return "", err
	}

	defer file.Close()

	gf, err := readGoldenFile(file)

	if err != nil {
		return "", err
	}

	content, ok := gf.Snapshots[m.TestName]

	if ok {
		return content, nil
	}

	return "", os.ErrNotExist
}

func (m *Matcher) getActualContent(actual interface{}) (string, error) {
	var buf bytes.Buffer
	transformed, err := m.Transformer.Transform(actual)

	if err != nil {
		return "", err
	}

	if err := m.Serializer.Serialize(&buf, transformed); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// FailureMessage implements GomegaMatcher.
func (m *Matcher) FailureMessage(actual interface{}) string {
	return m.getMessage(actual, "to")
}

// NegatedFailureMessage implements GomegaMatcher.
func (m *Matcher) NegatedFailureMessage(actual interface{}) string {
	return m.getMessage(actual, "not to")
}
