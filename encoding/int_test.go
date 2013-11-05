package encoding

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensures that an int can be written.
func TestWriteInt(t *testing.T) {
	var b bytes.Buffer
	WriteInt(&b, -100)
	assert.Equal(t, b.String(), `-100`)
}

// Ensures that a uint can be written.
func TestWriteUint(t *testing.T) {
	var b bytes.Buffer
	WriteUint(&b, uint(1230928137))
	assert.Equal(t, b.String(), `1230928137`)
}

func BenchmarkWriteInt(b *testing.B) {
	v := -3
	var w bytes.Buffer
	for i := 0; i < b.N; i++ {
		if err := WriteInt(&w, v); err != nil {
			b.Fatal("WriteInt:", err)
		}
	}
	b.SetBytes(int64(len("-3")))
}

func BenchmarkWriteUint(b *testing.B) {
	v := uint(30)
	var w bytes.Buffer
	for i := 0; i < b.N; i++ {
		if err := WriteUint(&w, v); err != nil {
			b.Fatal("WriteUint:", err)
		}
	}
	b.SetBytes(int64(len("30")))
}
