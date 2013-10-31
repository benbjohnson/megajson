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

// Ensures that a float32 can be written.
func TestWriteFloat32(t *testing.T) {
	var b bytes.Buffer
	WriteFloat32(&b, float32(2319.1921))
	assert.Equal(t, b.String(), `2319.1921`, "")
}

// Ensures that a float64 can be written.
func TestWriteFloat64(t *testing.T) {
	var b bytes.Buffer
	WriteFloat64(&b, 2319123.1921918273)
	assert.Equal(t, b.String(), `2.319123192191827e+06`, "")
}

// Ensures that a true boolean value can be written.
func TestWriteTrue(t *testing.T) {
	var b bytes.Buffer
	WriteBool(&b, true)
	assert.Equal(t, b.String(), `true`, "")
}

// Ensures that a false boolean value can be written.
func TestWriteFalse(t *testing.T) {
	var b bytes.Buffer
	WriteBool(&b, false)
	assert.Equal(t, b.String(), `false`, "")
}
