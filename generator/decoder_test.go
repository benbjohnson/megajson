package generator

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/parser"
	"go/token"
	"testing"
)

// Ensures a basic sanity check when generating the decoder.
func TestWriteTypeDecoder(t *testing.T) {
	var b bytes.Buffer
	src := `
package foo
type Foo struct {
    Name string
    Age int
}
`
	f, _ := parser.ParseFile(token.NewFileSet(), "foo.go", src, 0)
	err := writeFileDecoder(&b, f)
	assert.NoError(t, err)
}

// Ensures that a simple struct can be decoded from JSON.
func TestGenerateDecodeSimple(t *testing.T) {
	out, err := runDecodingFixture("decode/simple")
	assert.NoError(t, err)
	assert.Equal(t, out, `|foo|200|189273|2392|172389984|182.23|19380.1312|true|`)
}

// Ensures that a complex nested struct can be decoded from JSON.
func TestGenerateDecodeNested(t *testing.T) {
	out, err := runDecodingFixture("decode/nested")
	assert.NoError(t, err)
	assert.Equal(t, out, `|foo|John|20|`)
}

func runDecodingFixture(name string) (ret string, err error) {
	options := NewOptions()
	options.GenerateEncoder = false
	options.GenerateDecoder = true
	return runFixture(name, options)
}

