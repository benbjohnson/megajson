package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/benbjohnson/megajson/generator/decoder"
	"github.com/benbjohnson/megajson/generator/encoder"
)

var extregexp = regexp.MustCompile(`\.go$`)

// Generator generates encoders and decoders for Go files matching a given path.
type Generator interface {
	Generate(path string) error
}

type generator struct {
	decoder decoder.Generator
	encoder encoder.Generator
}

func New() Generator {
	return &generator{decoder:decoder.NewGenerator(), encoder:encoder.NewGenerator()}
}

// Generate recursively iterates over a path and generates encoders and decoders.
func (g *generator) Generate(path string) error {
	return filepath.Walk(path, g.walk)
}

// walk iterates is the callback used by Generate() for iterating over files and directories.
func (g *generator) walk(path string, info os.FileInfo, err error) error {
	// Only go file are used for generation.
	if info == nil {
		return fmt.Errorf("file not found: %s", path)
	} else if info.IsDir() || filepath.Ext(path) != ".go" {
		return nil
	}

	// Parse Go file.
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		return err
	}

	if err := g.encode(file, extregexp.ReplaceAllString(path, "_encoder.go"), info.Mode()); err != nil {
		return err
	}
	if err := g.decode(file, extregexp.ReplaceAllString(path, "_decoder.go"), info.Mode()); err != nil {
		return err
	}

	return nil
}

// decode generates a decoder file from a given Go file.
func (g *generator) decode(file *ast.File, path string, mode os.FileMode) error {
	var b bytes.Buffer
	if err := g.decoder.Generate(&b, file); err != nil {
		return err
	}
	if b.Len() > 0 {
		if err := ioutil.WriteFile(path, b.Bytes(), mode); err != nil {
			return err
		}
	}
	return nil
}

// encode generates an encoder file from a given Go file.
func (g *generator) encode(file *ast.File, path string, mode os.FileMode) error {
	var b bytes.Buffer
	if err := g.encoder.Generate(&b, file); err != nil {
		return err
	}
	if b.Len() > 0 {
		if err := ioutil.WriteFile(path, b.Bytes(), mode); err != nil {
			return err
		}
	}
	return nil
}
