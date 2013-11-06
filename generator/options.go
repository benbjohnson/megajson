package generator

// Options is a list of configuration options passed into the generator.
type Options struct {
	GenerateEncoder bool
	GenerateDecoder bool
}

// NewOptions constructs a default options object.
func NewOptions() *Options {
	return &Options{
		GenerateEncoder: true,
		GenerateDecoder: true,
	}
}
