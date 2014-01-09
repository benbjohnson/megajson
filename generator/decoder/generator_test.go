package decoder

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/benbjohnson/megajson/generator/test"
	"github.com/stretchr/testify/assert"
)

// Ensures a basic sanity check when generating the decoder.
func TestWriteTypeGenerator(t *testing.T) {
	src := `
package foo
type Foo struct {
    Name string
    Age int
}
`
	f, _ := parser.ParseFile(token.NewFileSet(), "foo.go", src, 0)
	err := NewGenerator().Generate(bytes.NewBufferString(src), f)
	assert.NoError(t, err)
}

// Ensures that a simple struct can be decoded from JSON.
func TestGenerateSimple(t *testing.T) {
	out, err := execute("simple")
	assert.NoError(t, err)
	assert.Equal(t, out, `|foo|200|189273|2392|172389984|182.23|19380.1312|true|`)
}

// Ensures that a complex nested struct can be decoded from JSON.
func TestGenerateDecodeNested(t *testing.T) {
	out, err := execute("nested")
	assert.NoError(t, err)
	assert.Equal(t, out, `|foo|John|20|<nil>|2|Jane|60|Jack|-13|`)
}

// execute generates a decoder against a fixture, executes the main prorgam, and returns the results.
func execute(name string) (ret string, err error) {
	test.Test(name, func(path string) {
		var file *ast.File
		file, err = parser.ParseFile(token.NewFileSet(), filepath.Join(path, "types.go"), nil, 0)
		if err != nil {
			return
		}

		// Generate decoder.
		f, _ := os.Create(filepath.Join(path, "decoder.go"))
		if err = NewGenerator().Generate(f, file); err != nil {
			fmt.Println("generate error:", err.Error())
			return
		}
		f.Close()

		// Execute fixture.
		out, _ := exec.Command("go", "run", filepath.Join(path, "decode.go"), filepath.Join(path, "decoder.go"), filepath.Join(path, "types.go")).CombinedOutput()
		ret = string(out)
	})
	return
}
