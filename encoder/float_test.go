package encoder

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensures that a float32 can be written.
func TestWriteFloat32(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	assert.NoError(t, e.WriteFloat32(float32(2319.1921)))
	assert.NoError(t, e.Flush())
	assert.Equal(t, b.String(), `2319.1921`)
}

// Ensures that a float64 can be written.
func TestWriteFloat64(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	assert.NoError(t, e.WriteFloat64(2319123.1921918273))
	assert.NoError(t, e.Flush())
	assert.Equal(t, b.String(), `2.319123192191827e+06`)
}

func BenchmarkWriteFloat32(b *testing.B) {
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	v := float32(2319.1921)
	for i := 0; i < b.N; i++ {
		if err := e.WriteFloat32(v); err != nil {
			b.Fatal("WriteFloat32:", err)
		}
	}
	e.Flush()
	b.SetBytes(int64(len("2319.1921")))
}

func BenchmarkWriteFloat64(b *testing.B) {
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	v := 2319123.1921918273
	for i := 0; i < b.N; i++ {
		if err := e.WriteFloat64(v); err != nil {
			b.Fatal("WriteFloat64:", err)
		}
	}
	e.Flush()
	b.SetBytes(int64(len(`2.319123192191827e+06`)))
}
