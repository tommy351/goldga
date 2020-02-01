package goldga

import (
	"encoding/json"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/davecgh/go-spew/spew"
	yaml "gopkg.in/yaml.v2"
)

// DefaultSerializer is the default serializer.
// nolint: gochecknoglobals
var DefaultSerializer Serializer = &DumpSerializer{
	Config: getDefaultDumpConfig(),
}

// Serializer serializes input and writes output to a writer.
type Serializer interface {
	Serialize(w io.Writer, input interface{}) error
}

// DumpSerializer serializes data using [go-spew](https://github.com/davecgh/go-spew).
type DumpSerializer struct {
	Config *spew.ConfigState
}

// Serialize implements Serializer
func (d *DumpSerializer) Serialize(w io.Writer, input interface{}) error {
	d.Config.Fdump(w, input)
	return nil
}

func getDefaultDumpConfig() *spew.ConfigState {
	conf := spew.NewDefaultConfig()
	conf.SortKeys = true
	conf.DisableCapacities = true

	return conf
}

// YAMLSerializer serializes data into YAML format.
type YAMLSerializer struct{}

// Serialize implements Serializer.
func (y *YAMLSerializer) Serialize(w io.Writer, input interface{}) error {
	enc := yaml.NewEncoder(w)
	defer enc.Close()

	return enc.Encode(input)
}

// JSONSerializer serializes data into JSON format.
type JSONSerializer struct {
	EscapeHTML   bool
	IndentPrefix string
	Indent       string
}

// Serialize implements Serializer.
func (j *JSONSerializer) Serialize(w io.Writer, input interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(j.EscapeHTML)
	enc.SetIndent(j.IndentPrefix, j.Indent)
	return enc.Encode(input)
}

// TOMLSerializer serializes data into TOML format.
type TOMLSerializer struct {
	Indent string
}

// Serialize implements Serializer.
func (t *TOMLSerializer) Serialize(w io.Writer, input interface{}) error {
	enc := toml.NewEncoder(w)
	enc.Indent = t.Indent
	return enc.Encode(input)
}
