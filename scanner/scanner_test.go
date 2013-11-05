package scanner

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
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
