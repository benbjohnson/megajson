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
	assert.NoError(t, e.Flush())
	assert.Equal(t, `"foo\u0009\n\r\"大"`, b.String())
}

// Ensures that a large string can be escaped and encoded.
func TestWriteStringLarge(t *testing.T) {
	var input, expected string
	for i := 0; i < 10000; i++ {
		input += "\t"
		expected += `\u0009`
	}
	input += "X"
	expected = "\"" + expected + "X\""

	var b bytes.Buffer
	e := NewEncoder(&b)
	err := e.WriteString(input)
	assert.NoError(t, e.Flush())
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
	assert.NoError(t, e.Flush())
	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(b.String()))
	if err == nil && len(expected) == len(b.String()) {
		assert.Equal(t, expected, b.String())
	}
}

// Ensures that a multiple strings can be encoded sequentially and share the same buffer.
func TestWriteMultipleStrings(t *testing.T) {
	var b bytes.Buffer
	var expected string
	e := NewEncoder(&b)

	for i := 0; i < 10000; i++ {
		err := e.WriteString("foo\t\n\r\"大\t")
		assert.NoError(t, err)
		expected += `"foo\u0009\n\r\"大\u0009"`
	}
	assert.NoError(t, e.Flush())
	assert.Equal(t, len(expected), len(b.String()))
	if len(expected) == len(b.String()) {
		assert.Equal(t, expected, b.String())
	}
}

// Ensures that a blank string can be encoded.
func TestWriteBlankString(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	e.WriteString("")
	assert.NoError(t, e.Flush())
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

// Ensures that a single byte can be written to the encoder.
func TestWriteByte(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	assert.NoError(t, e.WriteByte(':'))
	assert.NoError(t, e.Flush())
	assert.Equal(t, b.String(), `:`)
}

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

// Ensures that a null value can be written.
func TestWriteNull(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	assert.NoError(t, e.WriteNull())
	assert.NoError(t, e.Flush())
	assert.Equal(t, b.String(), `null`)
}

func BenchmarkWriteNull(b *testing.B) {
	var buf bytes.Buffer
	e := NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		if err := e.WriteNull(); err != nil {
			b.Fatal("WriteNull:", err)
		}
	}
	e.Flush()
	b.SetBytes(int64(len(`true`)))
}
