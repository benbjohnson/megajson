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
func TestGenerateEncodeSimple(t *testing.T) {
	out, err := runEncodingFixture("encode/simple")
	assert.NoError(t, err)
	assert.Equal(t, out, `{"StringX":"foo","IntX":200,"Int64X":189273,"myuint":2392,"Uint64X":172389984,"Float32X":182.23,"Float64X":19380.1312,"BoolX":true}`)
}

// Ensures that a nested struct can be encoded to JSON.
func TestGenerateEncodeNested(t *testing.T) {
	out, err := runEncodingFixture("encode/nested")
	assert.NoError(t, err)
	assert.Equal(t, out, `{"StringX":"foo","BX":{"Name":"John","Age":20},"BY":null,"Bn":[{"Name":"Jane","Age":60}],"Bn2":[]}`)
}

// Ensures that a simple struct can be encoded to JSON.
func TestGenerateDecodeSimple(t *testing.T) {
	out, err := runDecodingFixture("decode/simple")
	assert.NoError(t, err)
	assert.Equal(t, out, `|foo|200|189273|2392|172389984|182.23|19380.1312|true|`)
}

func runEncodingFixture(name string) (ret string, err error) {
	options := NewOptions()
	options.GenerateEncoder = true
	options.GenerateDecoder = false
	return runFixture(name, options)
}

func runDecodingFixture(name string) (ret string, err error) {
	options := NewOptions()
	options.GenerateEncoder = false
	options.GenerateDecoder = true
	return runFixture(name, options)
}

func runFixture(name string, options *Options) (ret string, err error) {
	withFixture(name, func(path string) {
		// Generate encoder.
		if err = Generate(path, options); err != nil {
			fmt.Println("Generate error:", err.Error())
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
