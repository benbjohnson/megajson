package writer

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

type Writer struct {
	w   io.Writer
	buf [actualBufSize + 64]byte
	pos int
}

// NewWriter creates a new JSON writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// Flush writes all data in the buffer to the writer.
func (w *Writer) Flush() error {
	if w.pos > 0 {
		if _, err := w.w.Write(w.buf[0:w.pos]); err != nil {
			return err
		}
		w.pos = 0
	}
	return nil
}

// check verifies there is space in the buffer.
func (w *Writer) check() error {
	if w.pos > actualBufSize {
		return w.Flush()
	}
	return nil
}

// writeByte writes a single byte to the buffer and increments the position.
func (w *Writer) writeByte(c byte) {
	w.buf[w.pos] = c
	w.pos++
}

// writeString writes a string to the buffer and increments the position.
func (w *Writer) writeString(s string) {
	copy(w.buf[w.pos:], s)
	w.pos += len(s)
}

// WriteByte writes a single byte.
func (w *Writer) WriteByte(c byte) error {
	if err := w.check(); err != nil {
		return err
	}
	w.buf[w.pos] = c
	w.pos++
	return nil
}

// WriteString writes a JSON string to the writer. Parts of this function are
// borrowed from the encoding/json package.
func (w *Writer) WriteString(v string) error {
	bufsz := (actualBufSize - w.pos) / maxByteEncodeSize

	w.writeByte('"')
	for i := 0; i < len(v); i += bufsz {
		if i > 0 {
			bufsz = bufSize
			if err := w.Flush(); err != nil {
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
					w.writeString(sub[prev:j])
				}
				switch b {
				case '\\':
					w.writeByte('\\')
					w.writeByte('\\')
				case '"':
					w.writeByte('\\')
					w.writeByte('"')
				case '\n':
					w.writeByte('\\')
					w.writeByte('n')
				case '\r':
					w.writeByte('\\')
					w.writeByte('r')
				default:
					// This encodes bytes < 0x20 except for \n and \r,
					// as well as < and >. The latter are escaped because they
					// can lead to security holes when user-controlled strings
					// are rendered into JSON and served to some browsers.
					w.writeByte('\\')
					w.writeByte('u')
					w.writeByte('0')
					w.writeByte('0')
					w.writeByte(hex[b>>4])
					w.writeByte(hex[b&0xF])
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
			w.writeString(sub[prev:sublen])
		}
	}
	w.writeByte('"')
	return nil
}

// WriteInt encodes and writes an integer.
func (w *Writer) WriteInt(v int) error {
	return w.WriteInt64(int64(v))
}

// WriteInt64 encodes and writes a 64-bit integer.
func (w *Writer) WriteInt64(v int64) error {
	if err := w.check(); err != nil {
		return err
	}

	buf := strconv.AppendInt(w.buf[w.pos:w.pos], v, 10)
	w.pos += len(buf)
	return nil
}

// WriteUint encodes and writes an unsigned integer.
func (w *Writer) WriteUint(v uint) error {
	return w.WriteUint64(uint64(v))
}

// WriteUint encodes and writes an unsigned integer.
func (w *Writer) WriteUint64(v uint64) error {
	if err := w.check(); err != nil {
		return err
	}

	buf := strconv.AppendUint(w.buf[w.pos:w.pos], v, 10)
	w.pos += len(buf)
	return nil
}

// WriteFloat32 encodes and writes a 32-bit float.
func (w *Writer) WriteFloat32(v float32) error {
	if err := w.check(); err != nil {
		return err
	}
	buf := strconv.AppendFloat(w.buf[w.pos:w.pos], float64(v), 'g', -1, 32)
	w.pos += len(buf)
	return nil
}

// WriteFloat64 encodes and writes a 64-bit float.
func (w *Writer) WriteFloat64(v float64) error {
	if err := w.check(); err != nil {
		return err
	}
	buf := strconv.AppendFloat(w.buf[w.pos:w.pos], v, 'g', -1, 64)
	w.pos += len(buf)
	return nil
}

// WriteBool writes a boolean.
func (w *Writer) WriteBool(v bool) error {
	if err := w.check(); err != nil {
		return err
	}
	if v {
		w.buf[w.pos+0] = 't'
		w.buf[w.pos+1] = 'r'
		w.buf[w.pos+2] = 'u'
		w.buf[w.pos+3] = 'e'
		w.pos += 4
	} else {
		w.buf[w.pos+0] = 'f'
		w.buf[w.pos+1] = 'a'
		w.buf[w.pos+2] = 'l'
		w.buf[w.pos+3] = 's'
		w.buf[w.pos+4] = 'e'
		w.pos += 5
	}
	return nil
}

// WriteNull writes "null".
func (w *Writer) WriteNull() error {
	if err := w.check(); err != nil {
		return err
	}
	w.buf[w.pos+0] = 'n'
	w.buf[w.pos+1] = 'u'
	w.buf[w.pos+2] = 'l'
	w.buf[w.pos+3] = 'l'
	w.pos += 4
	return nil
}

// WriteMap writes a map.
func (w *Writer) WriteMap(v map[string]interface{}) error {
	if err := w.check(); err != nil {
		return err
	}

	w.buf[w.pos] = '{'
	w.pos++

	var index int
	for key, value := range v {
		if index > 0 {
			w.buf[w.pos] = ','
			w.pos++
		}

		// Write key and colon.
		if err := w.WriteString(key); err != nil {
			return err
		}
		w.buf[w.pos] = ':'
		w.pos++

		// Write value.
		if value == nil {
			if err := w.WriteNull(); err != nil {
				return err
			}

		} else {
			switch value := value.(type) {
			case string:
				if err := w.WriteString(value); err != nil {
					return err
				}
			case int:
				if err := w.WriteInt(value); err != nil {
					return err
				}
			case int64:
				if err := w.WriteInt64(value); err != nil {
					return err
				}
			case uint:
				if err := w.WriteUint(value); err != nil {
					return err
				}
			case uint64:
				if err := w.WriteUint64(value); err != nil {
					return err
				}
			case float32:
				if err := w.WriteFloat32(value); err != nil {
					return err
				}
			case float64:
				if err := w.WriteFloat64(value); err != nil {
					return err
				}
			case bool:
				if err := w.WriteBool(value); err != nil {
					return err
				}
			case map[string]interface{}:
				if err := w.WriteMap(value); err != nil {
					return err
				}
			}
		}

		index++
	}

	w.buf[w.pos] = '}'
	w.pos++

	return nil
}
