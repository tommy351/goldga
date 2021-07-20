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

type Option func(*Matcher)

// WithDescription adds an optional description to the gold file, allowing multiple gold files per test.
func WithDescription(description string) Option {
	return func(matcher *Matcher) {
		if s, ok := matcher.Storage.(*SuiteStorage); ok {
			s.Name = fmt.Sprintf("%s (%s)", s.Name, description)
		}
	}
}

func getUpdateFile() bool {
	update, _ := strconv.ParseBool(os.Getenv("UPDATE_GOLDEN"))

	return update
}

func Match(options ...Option) *Matcher {
	m := &Matcher{
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
	for _, option := range options {
		option(m)
	}

	return m
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

	data, err := m.Storage.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

func (m *Matcher) getActualContent(actual interface{}) ([]byte, error) {
	var buf bytes.Buffer
	transformed, err := m.Transformer.Transform(actual)
	if err != nil {
		return nil, fmt.Errorf("transform error: %w", err)
	}

	if err := m.Serializer.Serialize(&buf, transformed); err != nil {
		return nil, fmt.Errorf("serialize error: %w", err)
	}

	return buf.Bytes(), nil
}

func (m *Matcher) FailureMessage(actual interface{}) string {
	return m.getMessage(actual, "to")
}

func (m *Matcher) NegatedFailureMessage(actual interface{}) string {
	return m.getMessage(actual, "not to")
}
