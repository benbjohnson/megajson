package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensures that a single byte can be written to the encoder.
func TestWriteByte(t *testing.T) {
	var b bytes.Buffer
	WriteByte(&b, ':')
	assert.Equal(t, b.String(), `:`)
}

// Ensures that a byte array can be written to the encoder.
func TestWriteBytes(t *testing.T) {
	var b bytes.Buffer
	WriteBytes(&b, []byte(`null`))
	assert.Equal(t, b.String(), `null`)
}
