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
