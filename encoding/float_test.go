package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensures that a float32 can be written.
func TestWriteFloat32(t *testing.T) {
	var b bytes.Buffer
	WriteFloat32(&b, float32(2319.1921))
	assert.Equal(t, b.String(), `2319.1921`)
}

// Ensures that a float64 can be written.
func TestWriteFloat64(t *testing.T) {
	var b bytes.Buffer
	WriteFloat64(&b, 2319123.1921918273)
	assert.Equal(t, b.String(), `2.319123192191827e+06`)
}


func BenchmarkWriteFloat32(b *testing.B) {
	v := float32(2319.1921)
	var w bytes.Buffer
	for i := 0; i < b.N; i++ {
		if err := WriteFloat32(&w, v); err != nil {
			b.Fatal("WriteFloat32:", err)
		}
	}
	b.SetBytes(int64(len("2319.1921")))
}

func BenchmarkWriteFloat64(b *testing.B) {
	v := 2319123.1921918273
	var w bytes.Buffer
	for i := 0; i < b.N; i++ {
		if err := WriteFloat64(&w, v); err != nil {
			b.Fatal("WriteFloat64:", err)
		}
	}
	b.SetBytes(int64(len(`2.319123192191827e+06`)))
}
