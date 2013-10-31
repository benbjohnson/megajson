package main

import (
	"io/ioutil"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

// Ensures that a simple struct can be encoded to JSON.
func TestSimpleEncoder(t *testing.T) {
	withFixture("simple", func(path string) {
		// TODO: Generate encoder.
		// TODO: Shell to `go run encode.go`.
		// TODO: Verify output.
	})
}


// Sets up a Go project using a given fixture directory.
func withFixture(name string, fn func(string)) {
	path, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(path)

	// TODO: Copy over fixture files.
	fn(path)
}
