package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Sets up a Go project using a given fixture directory.
func Test(name string, fn func(string)) {
	path, _ := ioutil.TempDir("", "")
	os.RemoveAll(path)
	defer os.RemoveAll(path)

	_, base, _, _ := runtime.Caller(0)
	run("cp", "-r", filepath.Join(filepath.Dir(base), ".fixtures", name), path)
	fn(path)
}

// Executes a command that is expected run successfully.
// On failure it dumps output and panics.
func run(name string, args ...string) {
	c := exec.Command(name, args...)
	if err := c.Run(); err != nil {
		fmt.Println(c.CombinedOutput())
		panic("fixture error: " + err.Error())
	}
}
