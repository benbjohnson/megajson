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
	withFixture("simple", func(path string) {
		// Generate encoder.
		err := Walk(path)
		assert.Nil(t, err, "")

		// Shell to `go run encode.go`.
		fmt.Println(path)
		files, _ := filepath.Glob(filepath.Join(path, "*"))
		args := []string{"run"}
		args = append(args, files...)
		c := exec.Command("go", args...)
		out, _ := c.Output()

		// Verify output.
		assert.Nil(t, err, "")
		assert.Equal(t, string(out), `{"Name":"","Age":0}`, "")
	})
}

// Sets up a Go project using a given fixture directory.
func withFixture(name string, fn func(string)) {
	path, _ := ioutil.TempDir("", "")
	os.RemoveAll(path)
	// defer os.RemoveAll(path)

	src, _ := filepath.Abs("../.fixtures/" + name)
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
