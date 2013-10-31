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

// Ensures that an int can be written.
func TestWriteInt(t *testing.T) {
	var b bytes.Buffer
	WriteInt(&b, -100)
	assert.Equal(t, b.String(), `-100`, "")
}

// Ensures that a uint can be written.
func TestWriteUint(t *testing.T) {
	var b bytes.Buffer
	WriteUint(&b, uint(1230928137))
	assert.Equal(t, b.String(), `1230928137`, "")
}
