package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var extregexp = regexp.MustCompile(`\.go$`)

// Walk recursively iterates over a path and generates encoders and decoders.
func Walk(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Only go file are used for generation.
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// Parse Go file.
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			fmt.Println("Parse error: ", err.Error())
			return err
		}

		// Write header.
		var b bytes.Buffer
		b.WriteString(header(file.Name.Name))

		// Loop over each spec and create a struct.
		encoded := false
		for _, decl := range file.Decls {
			if decl, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range decl.Specs {
					if spec, ok := spec.(*ast.TypeSpec); ok {
						err = writeTypeEncoder(&b, spec)
						if err != nil {
							return err
						}
						encoded = true
					}
				}
			}
		}

		// If no types were found to encode then skip this file.
		if !encoded {
			return nil
		}

		// fmt.Println(">>>>\n", b.String(), "\n<<<<<")

		// Format source.
		bfmt, err := format.Source(b.Bytes())
		if err != nil {
			return err
		}

		fmt.Println(">>>>\n", string(bfmt), "\n<<<<<")

		// Write to output file.
		if err := ioutil.WriteFile(extregexp.ReplaceAllString(path, "_encoder.go"), bfmt, info.Mode()); err != nil {
			return err
		}

		return nil
	})
}

// header is the string that is placed at the top of files.
func header(pkg string) string {
	s := fmt.Sprintf("package %s\n", pkg)
	s += "import (\n"
	s += "\"io\"\n"
	s += "\"github.com/benbjohnson/megajson/encoding\"\n"
	s += ")\n"
	return s
}
