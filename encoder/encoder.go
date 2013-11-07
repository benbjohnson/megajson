package encoder

import (
	"encoding/json"
	"io"
	"unicode/utf8"
)

const (
	// The maximum size that a byte can be encoded as.
	maxByteEncodeSize = 6

	// The size, in bytes, that can be encoded at a time.
	bufSize = 4096

	// The max size, in bytes, that an encoded value can be.
	actualBufSize = (4096 * maxByteEncodeSize)
)

// Encoder is an interface for the low-level JSON encoder.
type Encoder interface {
	Flush()
	WriteString(string) error
}

type encoder struct {
	w io.Writer
	scratch [64]byte
	buf [actualBufSize + 64]byte
	pos int
}

// NewEncoder creates a new encoder.
func NewEncoder(w io.Writer) Encoder {
	return &encoder{w: w}
}

// Flush writes all data in the buffer to the writer.
func (e *encoder) Flush() {
	if e.pos > 0 {
		e.w.Write(e.buf[0:e.pos])
		e.pos = 0
	}
}

// writeByte writes a single byte to the buffer and increments the position.
func (e *encoder) writeByte(c byte) {
	e.buf[e.pos] = c
	e.pos++
}

// writeBytes writes a byte array to the buffer and increments the position.
func (e *encoder) writeBytes(b []byte) {
	copy(e.buf[e.pos:], b)
	e.pos += len(b)
}

// writeBytes writes a string to the buffer and increments the position.
func (e *encoder) writeString(s string) {
	copy(e.buf[e.pos:], s)
	e.pos += len(s)
}

// WriteString writes a JSON string to the writer. Parts of this function are
// borrowed from the encoding/json package.
func (e *encoder) WriteString(v string) error {
	bufsz := (actualBufSize - e.pos) / maxByteEncodeSize

	e.writeByte('"')
	for i := 0; i < len(v); i += bufsz {
		if i > 0 {
			bufsz = bufSize
			e.Flush()
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
