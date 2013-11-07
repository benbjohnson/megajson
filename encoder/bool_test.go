package encoder

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Ensures that a true boolean value can be written.
func TestWriteTrue(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	assert.NoError(t, e.WriteBool(true))
	assert.NoError(t, e.Flush())
	assert.Equal(t, b.String(), `true`)
}

// Ensures that a false boolean value can be written.
func TestWriteFalse(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	assert.NoError(t, e.WriteBool(false))
	assert.NoError(t, e.Flush())
	assert.Equal(t, b.String(), `false`)
}

func BenchmarkWriteBool(b *testing.B) {
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		if err := e.WriteBool(true); err != nil {
			b.Fatal("WriteBool:", err)
		}
	}
	e.Flush()
	b.SetBytes(int64(len(`true`)))
}
