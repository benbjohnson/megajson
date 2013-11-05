package generator

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go/ast"
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
	err := writeTypeDecoder(&b, f.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec))
	assert.NoError(t, err)
}
