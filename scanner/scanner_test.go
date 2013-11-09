package scanner

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensures that a positive number can be scanned.
func TestScanPositiveNumber(t *testing.T) {
	tok, b, err := NewScanner(strings.NewReader("100")).Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TNUMBER)
	assert.Equal(t, string(b), "100")
}

// Ensures that a negative number can be scanned.
func TestScanNegativeNumber(t *testing.T) {
	tok, b, err := NewScanner(strings.NewReader("-1")).Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TNUMBER)
	assert.Equal(t, string(b), "-1")
}

// Ensures that a fractional number can be scanned.
func TestScanFloat(t *testing.T) {
	tok, b, err := NewScanner(strings.NewReader("120.12931")).Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TNUMBER)
	assert.Equal(t, string(b), "120.12931")
}

// Ensures that a quoted string can be scanned.
func TestScanString(t *testing.T) {
	tok, b, err := NewScanner(strings.NewReader(`"hello world"`)).Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TSTRING)
	assert.Equal(t, string(b), "hello world")
}

// Ensures that a quoted string with escaped characters can be scanned.
func TestScanEscapedString(t *testing.T) {
	tok, b, err := NewScanner(strings.NewReader(`"\"\\\/\b\f\n\r\t"`)).Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TSTRING)
	assert.Equal(t, string(b), "\"\\/\b\f\n\r\t")
}

// Ensures that a true value can be scanned.
func TestScanTrue(t *testing.T) {
	tok, _, err := NewScanner(strings.NewReader(`true`)).Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TTRUE)
}

// Ensures that a false value can be scanned.
func TestScanFalse(t *testing.T) {
	tok, _, err := NewScanner(strings.NewReader(`false`)).Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TFALSE)
}

// Ensures that a null value can be scanned.
func TestScanNull(t *testing.T) {
	tok, _, err := NewScanner(strings.NewReader(`null`)).Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TNULL)
}

// Ensures that an EOF gets returned.
func TestScanEOF(t *testing.T) {
	_, _, err := NewScanner(strings.NewReader(``)).Scan()
	assert.Equal(t, err, io.EOF)
}

// Ensures that a string can be read into a field.
func TestReadString(t *testing.T) {
	var v string
	err := NewScanner(strings.NewReader(`"foo"`)).ReadString(&v)
	assert.NoError(t, err)
	assert.Equal(t, v, "foo")
}

// Ensures that a non-string value is read into a string field as blank.
func TestReadNonStringAsString(t *testing.T) {
	var v string
	err := NewScanner(strings.NewReader(`12`)).ReadString(&v)
	assert.NoError(t, err)
	assert.Equal(t, v, "")
}

// Ensures that a non-value returns a read error.
func TestReadNonValueAsString(t *testing.T) {
	var v string
	err := NewScanner(strings.NewReader(`{`)).ReadString(&v)
	assert.Error(t, err)
}

// Ensures that an int can be read into a field.
func TestReadInt(t *testing.T) {
	var v int
	err := NewScanner(strings.NewReader(`100`)).ReadInt(&v)
	assert.NoError(t, err)
	assert.Equal(t, v, 100)
}

// Ensures that a non-number value is read into an int field as zero.
func TestReadNonNumberAsInt(t *testing.T) {
	var v int
	err := NewScanner(strings.NewReader(`"foo"`)).ReadInt(&v)
	assert.NoError(t, err)
	assert.Equal(t, v, 0)
}

// Ensures that an int64 can be read into a field.
func TestReadInt64(t *testing.T) {
	var v int64
	err := NewScanner(strings.NewReader(`-100`)).ReadInt64(&v)
	assert.NoError(t, err)
	assert.Equal(t, v, -100)
}

// Ensures that a uint can be read into a field.
func TestReadUint(t *testing.T) {
	var v uint
	err := NewScanner(strings.NewReader(`100`)).ReadUint(&v)
	assert.NoError(t, err)
	assert.Equal(t, v, uint(100))
}

// Ensures that an uint64 can be read into a field.
func TestReadUint64(t *testing.T) {
	var v uint64
	err := NewScanner(strings.NewReader(`1024`)).ReadUint64(&v)
	assert.NoError(t, err)
	assert.Equal(t, v, uint(1024))
}

// Ensures that a float32 can be read into a field.
func TestReadFloat32(t *testing.T) {
	var v float32
	err := NewScanner(strings.NewReader(`1293.123`)).ReadFloat32(&v)
	assert.NoError(t, err)
	assert.Equal(t, v, float32(1293.123))
}

// Ensures that a float64 can be read into a field.
func TestReadFloat64(t *testing.T) {
	var v float64
	err := NewScanner(strings.NewReader(`9871293.414123`)).ReadFloat64(&v)
	assert.NoError(t, err)
	assert.Equal(t, v, 9871293.414123)
}

// Ensures that a boolean can be read into a field.
func TestReadBoolTrue(t *testing.T) {
	var v bool
	err := NewScanner(strings.NewReader(`true`)).ReadBool(&v)
	assert.NoError(t, err)
	assert.Equal(t, v, true)
}

// Ensures whitespace between tokens are ignored.
func TestScanIgnoreWhitespace(t *testing.T) {
	s := NewScanner(strings.NewReader(" 100 true false "))

	tok, _, err := s.Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TNUMBER)

	tok, _, err = s.Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TTRUE)

	tok, _, err = s.Scan()
	assert.NoError(t, err)
	assert.Equal(t, tok, TFALSE)

	tok, _, err = s.Scan()
	assert.Equal(t, err, io.EOF)
	assert.Equal(t, tok, 0)
}

// Ensures that a map can be read into a field.
func TestReadMap(t *testing.T) {
	var v map[string]interface{}
	err := NewScanner(strings.NewReader(`{"foo":"bar", "bat":1293,"truex":true,"falsex":false,"nullx":null,"nested":{"xxx":"yyy"}}`)).ReadMap(&v)
	assert.NoError(t, err)
	assert.Equal(t, v["foo"], "bar")
	assert.Equal(t, v["bat"], float64(1293))
	assert.Equal(t, v["truex"], true)
	assert.Equal(t, v["falsex"], false)
	_, exists := v["nullx"]
	assert.Equal(t, v["nullx"], nil)
	assert.True(t, exists)
	assert.NotNil(t, v["nested"])
	nested := v["nested"].(map[string]interface{})
	assert.Equal(t, nested["xxx"], "yyy")
}

func BenchmarkScanNumber(b *testing.B) {
	withBuffer(b, "100", func(buf []byte) {
		s := NewScanner(bytes.NewBuffer(buf))
		for i := 0; i < b.N; i++ {
			if _, _, err := s.Scan(); err == io.EOF {
				s = NewScanner(bytes.NewBuffer(buf))
			} else if err != nil {
				b.Fatal("scan error:", err)
			}
		}
	})
}

func BenchmarkScanString(b *testing.B) {
	withBuffer(b, `"01234567"`, func(buf []byte) {
		s := NewScanner(bytes.NewBuffer(buf))
		for i := 0; i < b.N; i++ {
			if _, _, err := s.Scan(); err == io.EOF {
				s = NewScanner(bytes.NewBuffer(buf))
			} else if err != nil {
				b.Fatal("scan error:", err)
			}
		}
	})
}

func BenchmarkScanLongString(b *testing.B) {
	withBuffer(b, `"foo foo foo foo foo foo foo foo foo foo foo foo foo foo"`, func(buf []byte) {
		s := NewScanner(bytes.NewBuffer(buf))
		for i := 0; i < b.N; i++ {
			if _, _, err := s.Scan(); err == io.EOF {
				s = NewScanner(bytes.NewBuffer(buf))
			} else if err != nil {
				b.Fatal("scan error:", err)
			}
		}
	})
}

func BenchmarkScanEscapedString(b *testing.B) {
	withBuffer(b, `"\"\\\/\b\f\n\r\t"`, func(buf []byte) {
		s := NewScanner(bytes.NewBuffer(buf))
		for i := 0; i < b.N; i++ {
			if _, _, err := s.Scan(); err == io.EOF {
				s = NewScanner(bytes.NewBuffer(buf))
			} else if err != nil {
				b.Fatal("scan error:", err)
			}
		}
	})
}

func BenchmarkReadString(b *testing.B) {
	withBuffer(b, `"01234567"`, func(buf []byte) {
		var v string
		s := NewScanner(bytes.NewBuffer(buf))
		for i := 0; i < b.N; i++ {
			if err := s.ReadString(&v); err == io.EOF {
				s = NewScanner(bytes.NewBuffer(buf))
			} else if err != nil {
				b.Fatal("scan error:", err)
			}
		}
	})
}

func BenchmarkReadLongString(b *testing.B) {
	withBuffer(b, `"foo foo foo foo foo foo foo foo foo foo foo foo foo foo"`, func(buf []byte) {
		var v string
		s := NewScanner(bytes.NewBuffer(buf))
		for i := 0; i < b.N; i++ {
			if err := s.ReadString(&v); err == io.EOF {
				s = NewScanner(bytes.NewBuffer(buf))
			} else if err != nil {
				b.Fatal("scan error:", err)
			}
		}
	})
}

func BenchmarkReadInt(b *testing.B) {
	withBuffer(b, `"100"`, func(buf []byte) {
		var v int
		s := NewScanner(bytes.NewBuffer(buf))
		for i := 0; i < b.N; i++ {
			if err := s.ReadInt(&v); err == io.EOF {
				s = NewScanner(bytes.NewBuffer(buf))
			} else if err != nil {
				b.Fatal("scan error:", err)
			}
		}
	})
}

func BenchmarkReadFloat64(b *testing.B) {
	withBuffer(b, `"9871293.414123"`, func(buf []byte) {
		var v float64
		s := NewScanner(bytes.NewBuffer(buf))
		for i := 0; i < b.N; i++ {
			if err := s.ReadFloat64(&v); err == io.EOF {
				s = NewScanner(bytes.NewBuffer(buf))
			} else if err != nil {
				b.Fatal("scan error:", err)
			}
		}
	})
}

func BenchmarkReadBool(b *testing.B) {
	withBuffer(b, `true`, func(buf []byte) {
		var v bool
		s := NewScanner(bytes.NewBuffer(buf))
		for i := 0; i < b.N; i++ {
			if err := s.ReadBool(&v); err == io.EOF {
				s = NewScanner(bytes.NewBuffer(buf))
			} else if err != nil {
				b.Fatal("scan error:", err)
			}
		}
	})
}

func withBuffer(b *testing.B, value string, fn func([]byte)) {
	b.StopTimer()
	var str string
	for i := 0; i < 1000; i++ {
		str += value + " "
	}
	b.StartTimer()

	fn([]byte(str))

	b.SetBytes(int64(len(value)))
}
