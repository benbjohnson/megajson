package generator

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var extregexp = regexp.MustCompile(`\.go$`)

// Generate recursively iterates over a path and generates encoders and decoders.
func Generate(root string, options *Options) error {
	if options == nil {
		options = NewOptions()
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// Only go file are used for generation.
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// Parse Go file.
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			return err
		}

		if options.GenerateEncoder {
			if err := generateEncoder(file, extregexp.ReplaceAllString(path, "_encoder.go"), info.Mode()); err != nil {
				return err
			}
		}

		if options.GenerateDecoder {
			if err := generateDecoder(file, extregexp.ReplaceAllString(path, "_decoder.go"), info.Mode()); err != nil {
				return err
			}
		}

		return nil
	})
}

func generateEncoder(file *ast.File, path string, mode os.FileMode) error {
	var b bytes.Buffer
	if err := writeFileEncoder(&b, file); err != nil {
		return err
	}
	if b.Len() > 0 {
		if err := ioutil.WriteFile(path, b.Bytes(), mode); err != nil {
			return err
		}
	}
	return nil
}

func generateDecoder(file *ast.File, path string, mode os.FileMode) error {
	var b bytes.Buffer
	if err := writeFileDecoder(&b, file); err != nil {
		return err
	}
	if b.Len() > 0 {
		if err := ioutil.WriteFile(path, b.Bytes(), mode); err != nil {
			return err
		}
	}
	return nil
}


