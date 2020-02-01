package goldga

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

// nolint: gochecknoglobals
var defaultFs = afero.NewCacheOnReadFs(
	afero.NewOsFs(),
	afero.NewMemMapFs(),
	time.Minute,
)

func getGinkgoPath() string {
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

	return filepath.Join("testdata", name+".golden")
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
		Path:        getGinkgoPath(),
		Name:        getGinkgoTestName(),
		Serializer:  DefaultSerializer,
		Transformer: DefaultTransformer,
		UpdateFile:  getUpdateFile(),
		Color:       getColor(),
		fs:          defaultFs,
	}
}

var _ types.GomegaMatcher = (*Matcher)(nil)

// Matcher implements GomegaMatcher.
type Matcher struct {
	Path        string
	Name        string
	Serializer  Serializer
	Transformer Transformer
	Color       bool
	UpdateFile  bool

	fs afero.Fs
}

// Match implements GomegaMatcher.
func (m *Matcher) Match(actual interface{}) (bool, error) {
	actualContent, err := m.getActualContent(actual)

	if err != nil {
		return false, fmt.Errorf("failed to get actual content: %w", err)
	}

	expected, err := m.getExpectedContent()

	if err != nil {
		if !errors.Is(err, afero.ErrFileNotFound) {
			return false, fmt.Errorf("failed to get expected content: %w", err)
		}

		if err := m.writeFile(actualContent); err != nil {
			return false, fmt.Errorf("failed to write file: %w", err)
		}

		return true, nil
	}

	return expected == actualContent, nil
}

func (m *Matcher) readGoldenFile() (*goldenFile, error) {
	file, err := m.fs.Open(m.Path)

	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()

	return readGoldenFile(file)
}

func (m *Matcher) writeFile(actualContent string) error {
	if err := m.fs.MkdirAll(filepath.Dir(m.Path), os.ModePerm); err != nil {
		return err
	}

	gf, err := m.readGoldenFile()

	if err != nil {
		if !errors.Is(err, afero.ErrFileNotFound) {
			return fmt.Errorf("failed to read golden file: %w", err)
		}

		gf = newGoldenFile()
	}

	gf.Snapshots[m.Name] = actualContent

	file, err := m.fs.Create(m.Path)

	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()

	if err := writeGoldenFile(file, gf); err != nil {
		return fmt.Errorf("failed to write golden file: %w", err)
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
		m.Path,
		diffString(m.Color, expectedContent, actualContent))
}

func (m *Matcher) getExpectedContent() (string, error) {
	if m.UpdateFile {
		return "", afero.ErrFileNotFound
	}

	gf, err := m.readGoldenFile()

	if err != nil {
		return "", fmt.Errorf("failed to read golden file: %w", err)
	}

	content, ok := gf.Snapshots[m.Name]

	if ok {
		return content, nil
	}

	return "", afero.ErrFileNotFound
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
