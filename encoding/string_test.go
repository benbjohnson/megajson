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



func BenchmarkWriteRawBytes(b *testing.B) {
	s := "hello, world"
	var w bytes.Buffer
	for i := 0; i < b.N; i++ {
		if _, err := w.Write([]byte(s)); err != nil {
			b.Fatal("WriteRawBytes:", err)
		}
	}
	b.SetBytes(int64(len(s)))
}

func BenchmarkWriteString(b *testing.B) {
	s := "hello, world"
	var w bytes.Buffer
	for i := 0; i < b.N; i++ {
		if err := WriteString(&w, s); err != nil {
			b.Fatal("WriteString:", err)
		}
	}
	b.SetBytes(int64(len(s)))
}

func BenchmarkWriteStringReuseBuffer(b *testing.B) {
	s := "hello, world"
	var w bytes.Buffer
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		if err := WriteStringWithBuffer(&w, s, &buf); err != nil {
			b.Fatal("WriteStringReuseBuffer:", err)
		}
	}
	b.SetBytes(int64(len(s)))
}
