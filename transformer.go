package goldga

// DefaultTransformer is the default transformer.
// nolint: gochecknoglobals
var DefaultTransformer = &NopTransformer{}

// Transformer transforms input data.
type Transformer interface {
	Transform(input interface{}) (interface{}, error)
}

// NopTransformer doesn't transforms data.
type NopTransformer struct{}

// Transform implements Transformer.
func (n *NopTransformer) Transform(input interface{}) (interface{}, error) {
	return input, nil
}
