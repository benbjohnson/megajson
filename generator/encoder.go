package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"
)

func writeFileEncoder(w io.Writer, file *ast.File) error {
	// Ignore files without type specs.
	if types := getTypeSpecs(file); len(types) == 0 {
		return nil
	}

	var b bytes.Buffer
	if err := encoderTemplate.Execute(&b, file); err != nil {
		return err
	}

	bfmt, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Println("+++++\n", b.String(), "\n+++++")
		return err
	}

	// fmt.Println("+++++\n", string(bfmt), "\n+++++")

	if _, err := w.Write(bfmt); err != nil {
		return err
	}

	return nil
}


