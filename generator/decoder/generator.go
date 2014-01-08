package decoder

import (
	"bytes"
	"go/ast"
	"go/format"
	"io"
)

// Generator writes a generated JSON decoder to a writer.
type Generator interface {
	Generate(io.Writer, *ast.File) error
}

type generator struct{
}

// NewGenerator creates a new Generator instance.
func NewGenerator() Generator {
	return &generator{}
}

// Generator writes the generated decoder to the writer.
func (g *generator) Generate(w io.Writer, f *ast.File) error {
	// Ignore files without type specs.
	if len(types(f)) == 0 {
		return nil
	}

	// Generate code and the format the source code.
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, f); err != nil {
		return err
	}
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}
