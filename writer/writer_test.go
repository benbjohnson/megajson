package writer

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensures that a string can be escaped and encoded.
func TestWriteString(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	w.WriteString("foo\t\n\r\"大")
	assert.NoError(t, w.Flush())
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
	w := NewWriter(&b)
	err := w.WriteString(input)
	assert.NoError(t, w.Flush())
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
	w := NewWriter(&b)
	err := w.WriteString(input)
	assert.NoError(t, w.Flush())
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
	w := NewWriter(&b)

	for i := 0; i < 10000; i++ {
		err := w.WriteString("foo\t\n\r\"大\t")
		assert.NoError(t, err)
		expected += `"foo\u0009\n\r\"大\u0009"`
	}
	assert.NoError(t, w.Flush())
	assert.Equal(t, len(expected), len(b.String()))
	if len(expected) == len(b.String()) {
		assert.Equal(t, expected, b.String())
	}
}

// Ensures that a blank string can be encoded.
func TestWriteBlankString(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	w.WriteString("")
	assert.NoError(t, w.Flush())
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
	w := NewWriter(&buf)
	s := "hello, world"
	for i := 0; i < b.N; i++ {
		if err := w.WriteString(s); err != nil {
			b.Fatal("WriteString:", err)
		}
	}
	w.Flush()

	b.SetBytes(int64(len(s)))
}

// Ensures that an int can be written.
func TestWriteInt(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	assert.NoError(t, w.WriteInt(-100))
	assert.NoError(t, w.Flush())
	assert.Equal(t, b.String(), `-100`)
}

// Ensures that a uint can be written.
func TestWriteUint(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	assert.NoError(t, w.WriteUint(uint(1230928137)))
	assert.NoError(t, w.Flush())
	assert.Equal(t, b.String(), `1230928137`)
}

func BenchmarkWriteInt(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	v := -3
	for i := 0; i < b.N; i++ {
		if err := w.WriteInt(v); err != nil {
			b.Fatal("WriteInt:", err)
		}
	}
	w.Flush()
	b.SetBytes(int64(len("-3")))
}

func BenchmarkWriteUint(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	v := uint(30)
	for i := 0; i < b.N; i++ {
		if err := w.WriteUint(v); err != nil {
			b.Fatal("WriteUint:", err)
		}
	}
	b.SetBytes(int64(len("30")))
}

// Ensures that a float32 can be written.
func TestWriteFloat32(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	assert.NoError(t, w.WriteFloat32(float32(2319.1921)))
	assert.NoError(t, w.Flush())
	assert.Equal(t, b.String(), `2319.1921`)
}

// Ensures that a float64 can be written.
func TestWriteFloat64(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	assert.NoError(t, w.WriteFloat64(2319123.1921918273))
	assert.NoError(t, w.Flush())
	assert.Equal(t, b.String(), `2.319123192191827e+06`)
}

// Ensures that a simple map can be written.
func TestWriteSimpleMap(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	m := map[string]interface{}{
		"foo": "bar",
		"bat": "baz",
	}
	assert.NoError(t, w.WriteMap(m))
	assert.NoError(t, w.Flush())
	if b.String() != `{"foo":"bar","bat":"baz"}` && b.String() != `{"foo":"bar","bat":"baz"}` {
		t.Fatal("Invalid map encoding:", b.String())
	}
}

// Ensures that a more complex map can be written.
func TestWriteMap(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	m := map[string]interface{}{
		"stringx":  "foo",
		"intx":     100,
		"int64x":   int64(1023),
		"uintx":    uint(100),
		"uint64x":  uint64(1023),
		"float32x": float32(312.311),
		"float64x": float64(812731.19812),
		"truex":    true,
		"falsex":   false,
		"nullx":    nil,
	}
	assert.NoError(t, w.WriteMap(m))
	assert.NoError(t, w.Flush())
	assert.Contains(t, b.String(), `"intx":100`)
	assert.Contains(t, b.String(), `"int64x":1023`)
	assert.Contains(t, b.String(), `"uint64x":1023`)
	assert.Contains(t, b.String(), `"float32x":312.311`)
	assert.Contains(t, b.String(), `"float64x":812731.19812`)
	assert.Contains(t, b.String(), `"falsex":false`)
	assert.Contains(t, b.String(), `"nullx":null`)
	assert.Contains(t, b.String(), `"falsex":false`)
	assert.Contains(t, b.String(), `"stringx":"foo"`)
	assert.Contains(t, b.String(), `"uintx":100`)
	assert.Contains(t, b.String(), `"truex":true`)
}

// Ensures that a nested map can be written.
func TestWriteNestedMap(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	m := map[string]interface{}{
		"foo": map[string]interface{}{"bar": "bat"},
	}
	assert.NoError(t, w.WriteMap(m))
	assert.NoError(t, w.Flush())
	assert.Equal(t, b.String(), `{"foo":{"bar":"bat"}}`)
}

func BenchmarkWriteFloat32(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	v := float32(2319.1921)
	for i := 0; i < b.N; i++ {
		if err := w.WriteFloat32(v); err != nil {
			b.Fatal("WriteFloat32:", err)
		}
	}
	w.Flush()
	b.SetBytes(int64(len("2319.1921")))
}

func BenchmarkWriteFloat64(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	v := 2319123.1921918273
	for i := 0; i < b.N; i++ {
		if err := w.WriteFloat64(v); err != nil {
			b.Fatal("WriteFloat64:", err)
		}
	}
	w.Flush()
	b.SetBytes(int64(len(`2.319123192191827e+06`)))
}

// Ensures that a single byte can be written to the writer.
func TestWriteByte(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	assert.NoError(t, w.WriteByte(':'))
	assert.NoError(t, w.Flush())
	assert.Equal(t, b.String(), `:`)
}

// Ensures that a true boolean value can be written.
func TestWriteTrue(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	assert.NoError(t, w.WriteBool(true))
	assert.NoError(t, w.Flush())
	assert.Equal(t, b.String(), `true`)
}

// Ensures that a false boolean value can be written.
func TestWriteFalse(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	assert.NoError(t, w.WriteBool(false))
	assert.NoError(t, w.Flush())
	assert.Equal(t, b.String(), `false`)
}

func BenchmarkWriteBool(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	for i := 0; i < b.N; i++ {
		if err := w.WriteBool(true); err != nil {
			b.Fatal("WriteBool:", err)
		}
	}
	w.Flush()
	b.SetBytes(int64(len(`true`)))
}

// Ensures that a null value can be written.
func TestWriteNull(t *testing.T) {
	var b bytes.Buffer
	w := NewWriter(&b)
	assert.NoError(t, w.WriteNull())
	assert.NoError(t, w.Flush())
	assert.Equal(t, b.String(), `null`)
}

func BenchmarkWriteNull(b *testing.B) {
	var buf bytes.Buffer
	w := NewWriter(&buf)
	for i := 0; i < b.N; i++ {
		if err := w.WriteNull(); err != nil {
			b.Fatal("WriteNull:", err)
		}
	}
	w.Flush()
	b.SetBytes(int64(len(`true`)))
}
