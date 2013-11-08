package encoder

import (
	"encoding/json"
	"io"
	"strconv"
	"unicode/utf8"
)

const (
	// The maximum size that a byte can be encoded as.
	maxByteEncodeSize = 6

	// The size, in bytes, that can be encoded at a time.
	bufSize = 16384

	// The max size, in bytes, that an encoded value can be.
	actualBufSize = (bufSize * maxByteEncodeSize)
)

var hex = "0123456789abcdef"

// Encoder is an interface for the low-level JSON encoder.
type Encoder interface {
	Flush() error
	WriteString(string) error
	WriteByte(byte) error
	WriteInt(int) error
	WriteInt64(int64) error
	WriteUint(uint) error
	WriteUint64(uint64) error
	WriteFloat32(float32) error
	WriteFloat64(float64) error
	WriteBool(bool) error
	WriteNull() error
}

type encoder struct {
	w io.Writer
	buf [actualBufSize + 64]byte
	pos int
}

// NewEncoder creates a new encoder.
func NewEncoder(w io.Writer) Encoder {
	return &encoder{w: w}
}

// Flush writes all data in the buffer to the writer.
func (e *encoder) Flush() error {
	if e.pos > 0 {
		if _, err := e.w.Write(e.buf[0:e.pos]); err != nil {
			return err
		}
		e.pos = 0
	}
	return nil
}

// check verifies there is space in the buffer.
func (e *encoder) check() error {
	if e.pos > actualBufSize {
		return e.Flush();
	}
	return nil
}

// writeByte writes a single byte to the buffer and increments the position.
func (e *encoder) writeByte(c byte) {
	e.buf[e.pos] = c
	e.pos++
}

// writeString writes a string to the buffer and increments the position.
func (e *encoder) writeString(s string) {
	copy(e.buf[e.pos:], s)
	e.pos += len(s)
}

// WriteByte writes a single byte.
func (e *encoder) WriteByte(c byte) error {
	if err := e.check(); err != nil {
		return err
	}
	e.buf[e.pos] = c
	e.pos++
	return nil
}

// WriteString writes a JSON string to the writer. Parts of this function are
// borrowed from the encoding/json package.
func (e *encoder) WriteString(v string) error {
	bufsz := (actualBufSize - e.pos) / maxByteEncodeSize

	e.writeByte('"')
	for i := 0; i < len(v); i += bufsz {
		if i > 0 {
			bufsz = bufSize
			if err := e.Flush(); err != nil {
				return err
			}
		}

		// Extract substring.
		end := i + bufsz
		if end > len(v) {
			end = len(v)
		}
		bufend := end + utf8.UTFMax
		if bufend > len(v) {
			bufend = len(v)
		}
		sub := v[i:bufend]
		sublen := end - i
	
		prev := 0
		for j := 0; j < sublen; {
			if b := sub[j]; b < utf8.RuneSelf {
				if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' {
					j++
					continue
				}
				if prev < j {
					e.writeString(sub[prev:j])
				}
				switch b {
				case '\\':
					e.writeByte('\\')
					e.writeByte('\\')
				case '"':
					e.writeByte('\\')
					e.writeByte('"')
				case '\n':
					e.writeByte('\\')
					e.writeByte('n')
				case '\r':
					e.writeByte('\\')
					e.writeByte('r')
				default:
					// This encodes bytes < 0x20 except for \n and \r,
					// as well as < and >. The latter are escaped because they
					// can lead to security holes when user-controlled strings
					// are rendered into JSON and served to some browsers.
					e.writeByte('\\')
					e.writeByte('u')
					e.writeByte('0')
					e.writeByte('0')
					e.writeByte(hex[b>>4])
					e.writeByte(hex[b&0xF])
				}
				j++
				prev = j
				continue
			}
			c, size := utf8.DecodeRuneInString(sub[j:])
			if c == utf8.RuneError && size == 1 {
				return &json.InvalidUTF8Error{S: v}
			}
			j += size

			// If we cross the buffer end then adjust the outer loop
			if j > bufsz {
				i += j - bufsz
				sublen += j - bufsz
			}
		}
		if prev < sublen {
			e.writeString(sub[prev:sublen])
		}
	}
	e.writeByte('"')
	return nil
}

// WriteInt encodes and writes an integer.
func (e *encoder) WriteInt(v int) error {
	return e.WriteInt64(int64(v))
}

// WriteInt64 encodes and writes a 64-bit integer.
func (e *encoder) WriteInt64(v int64) error {
	if err := e.check(); err != nil {
		return err
	}

	buf := strconv.AppendInt(e.buf[e.pos:e.pos], v, 10)
	e.pos += len(buf)
	return nil
}

// WriteUint encodes and writes an unsigned integer.
func (e *encoder) WriteUint(v uint) error {
	return e.WriteUint64(uint64(v))
}

// WriteUint encodes and writes an unsigned integer.
func (e *encoder) WriteUint64(v uint64) error {
	if err := e.check(); err != nil {
		return err
	}

	buf := strconv.AppendUint(e.buf[e.pos:e.pos], v, 10)
	e.pos += len(buf)
	return nil
}

// WriteFloat32 encodes and writes a 32-bit float.
func (e *encoder) WriteFloat32(v float32) error {
	if err := e.check(); err != nil {
		return err
	}
	buf := strconv.AppendFloat(e.buf[e.pos:e.pos], float64(v), 'g', -1, 32)
	e.pos += len(buf)
	return nil
}

// WriteFloat64 encodes and writes a 64-bit float.
func (e *encoder) WriteFloat64(v float64) error {
	if err := e.check(); err != nil {
		return err
	}
	buf := strconv.AppendFloat(e.buf[e.pos:e.pos], v, 'g', -1, 64)
	e.pos += len(buf)
	return nil
}

// WriteBool writes a boolean.
func (e *encoder) WriteBool(v bool) error {
	if err := e.check(); err != nil {
		return err
	}
	if v {
		e.buf[e.pos+0] = 't'
		e.buf[e.pos+1] = 'r'
		e.buf[e.pos+2] = 'u'
		e.buf[e.pos+3] = 'e'
		e.pos += 4
	} else {
		e.buf[e.pos+0] = 'f'
		e.buf[e.pos+1] = 'a'
		e.buf[e.pos+2] = 'l'
		e.buf[e.pos+3] = 's'
		e.buf[e.pos+4] = 'e'
		e.pos += 5
	}
	return nil
}

// WriteNull writes "null".
func (e *encoder) WriteNull() error {
	if err := e.check(); err != nil {
		return err
	}
	e.buf[e.pos+0] = 'n'
	e.buf[e.pos+1] = 'u'
	e.buf[e.pos+2] = 'l'
	e.buf[e.pos+3] = 'l'
	e.pos += 4
	return nil
}
