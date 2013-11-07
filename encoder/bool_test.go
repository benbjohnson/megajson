package encoder

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Ensures that a true boolean value can be written.
func TestWriteTrue(t *testing.T) {
	var b bytes.Buffer
	WriteBool(&b, true)
	assert.Equal(t, b.String(), `true`)
}

// Ensures that a false boolean value can be written.
func TestWriteFalse(t *testing.T) {
	var b bytes.Buffer
	WriteBool(&b, false)
	assert.Equal(t, b.String(), `false`)
}

func BenchmarkWriteBool(b *testing.B) {
	var w bytes.Buffer
	for i := 0; i < b.N; i++ {
		if err := WriteBool(&w, true); err != nil {
			b.Fatal("WriteBool:", err)
		}
	}
	b.SetBytes(int64(len(`true`)))
}
