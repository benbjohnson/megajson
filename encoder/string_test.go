package encoder

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensures that a string can be escaped and encoded.
func TestWriteString(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	e.WriteString("foo\t\n\r\"大")
	e.Flush()
	assert.Equal(t, `"foo\u0009\n\r\"大"`, b.String())
}

// Ensures that a large string can be escaped and encoded.
func TestWriteStringLarge(t *testing.T) {
	var input, expected string
	for i := 0; i < 10000; i++ {
		input += "\t"
		expected += `\u0009`
	}
	expected = "\"" + expected + "\""

	var b bytes.Buffer
	e := NewEncoder(&b)
	err := e.WriteString(input)
	e.Flush()
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(b.String()))
	if err == nil && len(expected) == len(b.String()) {
		assert.Equal(t, expected, b.String())
	}
}

// Ensures that a large unicode string can be escaped and encoded.
func TestWriteStringLargeUnicode(t *testing.T) {
	var input, expected string
	for i := 0; i < 10000; i++ {
		input += "大"
		expected += "大"
	}
	expected = "\"" + expected + "\""

	var b bytes.Buffer
	e := NewEncoder(&b)
	err := e.WriteString(input)
	e.Flush()
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(b.String()))
	//if err == nil && len(expected) == len(b.String()) {
		assert.Equal(t, expected, b.String())
	//}
}

// Ensures that a blank string can be encoded.
func TestWriteBlankString(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	e.WriteString("")
	e.Flush()
	assert.Equal(t, b.String(), `""`)
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
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	s := "hello, world"
	for i := 0; i < b.N; i++ {
		if err := e.WriteString(s); err != nil {
			b.Fatal("WriteString:", err)
		}
	}
	e.Flush()

	b.SetBytes(int64(len(s)))
}
