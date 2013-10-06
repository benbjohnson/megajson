package generator

import (
    "bytes"
    "go/ast"
    "go/parser"
    "go/token"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestGeneratorGenerate(t *testing.T) {
    var b bytes.Buffer
    src := `
    package foo
    type Foo struct {
        Name string
    }
    `
    f, _ := parser.ParseFile(token.NewFileSet(), "foo.go", src, 0)
    err := GenerateEncoder("foo", f.Decls[0].(*ast.GenDecl).Specs[0].(*ast.TypeSpec), &b)
    assert.Nil(t, err)
    assert.Equal(t, b.String(), "XXX")
}
