package encoder

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

// Ensures a basic sanity check when generating the encoder.
func TestWriteTypeEncoder(t *testing.T) {
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

// Ensures that a simple struct can be encoded to JSON.
func TestGenerateEncodeSimple(t *testing.T) {
	out, err := execute("simple")
	assert.NoError(t, err)
	assert.Equal(t, out, `{"StringX":"foo","IntX":200,"Int64X":189273,"myuint":2392,"Uint64X":172389984,"Float32X":182.23,"Float64X":19380.1312,"BoolX":true}`)
}

// Ensures that a nested struct can be encoded to JSON.
func TestGenerateEncodeNested(t *testing.T) {
	out, err := execute("nested")
	assert.NoError(t, err)
	assert.Equal(t, out, `{"StringX":"foo","BX":{"Name":"John","Age":20},"BY":null,"Bn":[{"Name":"Jane","Age":60}],"Bn2":[]}`)
}

// execute generates an encoder against a fixture, executes the main prorgam, and returns the results.
func execute(name string) (ret string, err error) {
	test.Test(name, func(path string) {
		var file *ast.File
		file, err = parser.ParseFile(token.NewFileSet(), filepath.Join(path, "types.go"), nil, 0)
		if err != nil {
			return
		}

		// Generate decoder.
		f, _ := os.Create(filepath.Join(path, "encoder.go"))
		if err = NewGenerator().Generate(f, file); err != nil {
			fmt.Println("generate error:", err.Error())
			return
		}
		f.Close()

		// Execute fixture.
		out, _ := exec.Command("go", "run", filepath.Join(path, "encode.go"), filepath.Join(path, "encoder.go"), filepath.Join(path, "types.go")).CombinedOutput()
		ret = string(out)
	})
	return
}
