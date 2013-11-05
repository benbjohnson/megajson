package generator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Ensures that a simple struct can be encoded to JSON.
func TestGeneratorWalkSimpleEncoder(t *testing.T) {
	out, err := runFixture("encode.simple")
	assert.NoError(t, err)
	assert.Equal(t, out, `{"StringX":"foo","IntX":200,"Int64X":189273,"myuint":2392,"Uint64X":172389984,"Float32X":182.23,"Float64X":19380.1312,"BoolX":true}`)
}

// Ensures that a nested struct can be encoded to JSON.
func TestGeneratorWalkNestedEncoder(t *testing.T) {
	out, err := runFixture("encode.nested")
	assert.NoError(t, err)
	assert.Equal(t, out, `{"StringX":"foo","BX":{"Name":"John","Age":20},"BY":null,"Bn":[{"Name":"Jane","Age":60}],"Bn2":[]}`)
}

// Ensures that a simple struct can be encoded to JSON.
func TestGeneratorWalkSimpleDecoder(t *testing.T) {
	/*
	out, err := runFixture("simple.decode")
	assert.NoError(t, err)
	assert.Equal(t, out, `{"StringX":"foo","IntX":200,"Int64X":189273,"myuint":2392,"Uint64X":172389984,"Float32X":182.23,"Float64X":19380.1312,"BoolX":true}`)
	*/
}

func runFixture(name string) (ret string, err error) {
	withFixture(name, func(path string) {
		// Generate encoder.
		if err = Walk(path); err != nil {
			return
		}

		// Shell to `go run encode.go`.
		files, _ := filepath.Glob(filepath.Join(path, "*"))
		args := []string{"run"}
		args = append(args, files...)
		c := exec.Command("go", args...)
		out, _ := c.CombinedOutput()
		ret = string(out)
	})
	return
}

// Sets up a Go project using a given fixture directory.
func withFixture(name string, fn func(string)) {
	path, _ := ioutil.TempDir("", "")
	os.RemoveAll(path)
	defer os.RemoveAll(path)

	src, _ := filepath.Abs("../test/.fixtures/" + name)
	mustRun("cp", "-r", src, path)
	fn(path)
}

// Executes a command that is expected run successfully. Otherwise dumps output and panics.
func mustRun(name string, args ...string) {
	c := exec.Command(name, args...)
	if err := c.Run(); err != nil {
		fmt.Println(c.CombinedOutput())
		panic("Fixture error: " + err.Error())
	}
}
