package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

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
