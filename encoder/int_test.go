package encoder

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensures that an int can be written.
func TestWriteInt(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	assert.NoError(t, e.WriteInt(-100))
	assert.NoError(t, e.Flush())
	assert.Equal(t, b.String(), `-100`)
}

// Ensures that a uint can be written.
func TestWriteUint(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	assert.NoError(t, e.WriteUint(uint(1230928137)))
	assert.NoError(t, e.Flush())
	assert.Equal(t, b.String(), `1230928137`)
}

func BenchmarkWriteInt(b *testing.B) {
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	v := -3
	for i := 0; i < b.N; i++ {
		if err := e.WriteInt(v); err != nil {
			b.Fatal("WriteInt:", err)
		}
	}
	e.Flush()
	b.SetBytes(int64(len("-3")))
}

func BenchmarkWriteUint(b *testing.B) {
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	v := uint(30)
	for i := 0; i < b.N; i++ {
		if err := e.WriteUint(v); err != nil {
			b.Fatal("WriteUint:", err)
		}
	}
	b.SetBytes(int64(len("30")))
}
