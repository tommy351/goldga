package goldga

import (
	"encoding/json"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/davecgh/go-spew/spew"
	yaml "gopkg.in/yaml.v2"
)

// nolint: gochecknoglobals
var (
	DefaultSerializer Serializer = &DumpSerializer{
		Config: newDefaultDumpConfig(),
	}
)

type Serializer interface {
	Serialize(w io.Writer, input interface{}) error
}

type DumpSerializer struct {
	Config *spew.ConfigState
}

func (d *DumpSerializer) Serialize(w io.Writer, input interface{}) error {
	d.Config.Fdump(w, input)
	return nil
}

func newDefaultDumpConfig() *spew.ConfigState {
	conf := spew.NewDefaultConfig()
	conf.SortKeys = true
	conf.DisableCapacities = true

	return conf
}

type YAMLSerializer struct{}

func (y *YAMLSerializer) Serialize(w io.Writer, input interface{}) error {
	enc := yaml.NewEncoder(w)
	defer enc.Close()

	return enc.Encode(input)
}

type JSONSerializer struct {
	EscapeHTML   bool
	IndentPrefix string
	Indent       string
}

func (j *JSONSerializer) Serialize(w io.Writer, input interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(j.EscapeHTML)
	enc.SetIndent(j.IndentPrefix, j.Indent)
	return enc.Encode(input)
}

type TOMLSerializer struct {
	Indent string
}

func (t *TOMLSerializer) Serialize(w io.Writer, input interface{}) error {
	enc := toml.NewEncoder(w)
	enc.Indent = t.Indent
	return enc.Encode(input)
}
