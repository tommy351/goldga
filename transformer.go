package goldga

// nolint: gochecknoglobals
var (
	DefaultTransformer Transformer = &NopTransformer{}
)

type Transformer interface {
	Transform(input interface{}) (interface{}, error)
}

var _ Transformer = (*NopTransformer)(nil)

type NopTransformer struct{}

func (NopTransformer) Transform(input interface{}) (interface{}, error) {
	return input, nil
}
