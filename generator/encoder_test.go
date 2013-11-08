package generator

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/parser"
	"go/token"
	"testing"
)

// Ensures a basic sanity check when generating the encoder.
func TestWriteTypeEncoder(t *testing.T) {
	var b bytes.Buffer
	src := `
package foo
type Foo struct {
    Name string
    Age int
}
`
	f, _ := parser.ParseFile(token.NewFileSet(), "foo.go", src, 0)
	err := writeFileEncoder(&b, f)
	assert.NoError(t, err)
}

// Ensures that a simple struct can be encoded to JSON.
func TestGenerateEncodeSimple(t *testing.T) {
	out, err := runEncodingFixture("encode/simple")
	assert.NoError(t, err)
	assert.Equal(t, out, `{"StringX":"foo","IntX":200,"Int64X":189273,"myuint":2392,"Uint64X":172389984,"Float32X":182.23,"Float64X":19380.1312,"BoolX":true}`)
}

// Ensures that a nested struct can be encoded to JSON.
func TestGenerateEncodeNested(t *testing.T) {
	out, err := runEncodingFixture("encode/nested")
	assert.NoError(t, err)
	assert.Equal(t, out, `{"StringX":"foo","BX":{"Name":"John","Age":20},"BY":null,"Bn":[{"Name":"Jane","Age":60}],"Bn2":[]}`)
}

func runEncodingFixture(name string) (ret string, err error) {
	options := NewOptions()
	options.GenerateEncoder = true
	options.GenerateDecoder = false
	return runFixture(name, options)
}
