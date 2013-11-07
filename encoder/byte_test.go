package encoder

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensures that a single byte can be written to the encoder.
func TestWriteByte(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	assert.NoError(t, e.WriteByte(':'))
	assert.NoError(t, e.Flush())
	assert.Equal(t, b.String(), `:`)
}
