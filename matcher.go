package goldga

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/onsi/gomega/types"
	"github.com/spf13/afero"
)

func getUpdateFile() bool {
	update, _ := strconv.ParseBool(os.Getenv("UPDATE_GOLDEN"))
	return update
}

func Match() *Matcher {
	return &Matcher{
		Serializer:  DefaultSerializer,
		Transformer: DefaultTransformer,
		Storage: &SuiteStorage{
			Path: getGinkgoPath(),
			Name: getGinkgoTestName(),
			Fs:   defaultFs,
		},
		Differ:     DefaultDiffer,
		UpdateFile: getUpdateFile(),
	}
}

var _ types.GomegaMatcher = (*Matcher)(nil)

type Matcher struct {
	Serializer  Serializer
	Transformer Transformer
	Storage     Storage
	Differ      Differ
	UpdateFile  bool
}

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

		if err := m.Storage.Write(actualContent); err != nil {
			return false, fmt.Errorf("faield to write file: %w", err)
		}

		return true, nil
	}

	return bytes.Equal(expected, actualContent), nil
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

	return fmt.Sprintf("Expected %s match the golden file\n%s",
		message,
		m.Differ.Diff(expectedContent, actualContent))
}

func (m *Matcher) getExpectedContent() ([]byte, error) {
	if m.UpdateFile {
		return nil, afero.ErrFileNotFound
	}

	return m.Storage.Read()
}

func (m *Matcher) getActualContent(actual interface{}) ([]byte, error) {
	var buf bytes.Buffer
	transformed, err := m.Transformer.Transform(actual)

	if err != nil {
		return nil, err
	}

	if err := m.Serializer.Serialize(&buf, transformed); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *Matcher) FailureMessage(actual interface{}) string {
	return m.getMessage(actual, "to")
}

func (m *Matcher) NegatedFailureMessage(actual interface{}) string {
	return m.getMessage(actual, "not to")
}
