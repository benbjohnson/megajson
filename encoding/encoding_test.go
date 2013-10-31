package encoding

import (
	"bytes"
	"testing"
	"github.com/stretchr/testify/assert"
)

// Ensures that a string can be escaped and encoded.
func TestWriteString(t *testing.T) {
	var b bytes.Buffer
	WriteString(&b, "foo\n\"")
	assert.Equal(t, b.String(), `"foo\n\""`, "")
}

// Ensures that a blank string can be encoded.
func TestWriteBlankString(t *testing.T) {
	var b bytes.Buffer
	WriteString(&b, "")
	assert.Equal(t, b.String(), `""`, "")
}
