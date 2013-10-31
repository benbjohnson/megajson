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



// Ensures that an int can be written.
func TestWriteInt(t *testing.T) {
	var b bytes.Buffer
	WriteInt(&b, -100)
	assert.Equal(t, b.String(), `-100`, "")
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



// Ensures that a uint can be written.
func TestWriteUint(t *testing.T) {
	var b bytes.Buffer
	WriteUint(&b, uint(1230928137))
	assert.Equal(t, b.String(), `1230928137`, "")
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



// Ensures that a float32 can be written.
func TestWriteFloat32(t *testing.T) {
	var b bytes.Buffer
	WriteFloat32(&b, float32(2319.1921))
	assert.Equal(t, b.String(), `2319.1921`, "")
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


// Ensures that a float64 can be written.
func TestWriteFloat64(t *testing.T) {
	var b bytes.Buffer
	WriteFloat64(&b, 2319123.1921918273)
	assert.Equal(t, b.String(), `2.319123192191827e+06`, "")
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


// Ensures that a true boolean value can be written.
func TestWriteTrue(t *testing.T) {
	var b bytes.Buffer
	WriteBool(&b, true)
	assert.Equal(t, b.String(), `true`, "")
}

// Ensures that a false boolean value can be written.
func TestWriteFalse(t *testing.T) {
	var b bytes.Buffer
	WriteBool(&b, false)
	assert.Equal(t, b.String(), `false`, "")
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
