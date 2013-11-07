package encoder

import (
	"bytes"
	"encoding/json"
	"io"
	"unicode/utf8"
)

var hex = "0123456789abcdef"

// WriteString writes a JSON string to the writer. JSON encoding used
// from the encoding/json package.
func WriteString(w io.Writer, v string) error {
	var buf bytes.Buffer
	return WriteStringWithBuffer(w, v, &buf)
}

// WriteString writes a JSON string to the writer. JSON encoding used
// from the encoding/json package.
func WriteStringWithBuffer(w io.Writer, v string, buf *bytes.Buffer) error {
	buf.Reset()
	buf.WriteByte('"')
	prev := 0
	for i := 0; i < len(v); {
		if b := v[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' {
				i++
				continue
			}
			if prev < i {
				buf.WriteString(v[prev:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			default:
				// This encodes bytes < 0x20 except for \n and \r,
				// as well as < and >. The latter are escaped because they
				// can lead to security holes when user-controlled strings
				// are rendered into JSON and served to some browsers.
				buf.WriteString(`\u00`)
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			}
			i++
			prev = i
			continue
		}
		c, size := utf8.DecodeRuneInString(v[i:])
		if c == utf8.RuneError && size == 1 {
			return &json.InvalidUTF8Error{S: v}
		}
		i += size
	}
	if prev < len(v) {
		buf.WriteString(v[prev:])
	}
	buf.WriteByte('"')
	w.Write(buf.Bytes())
	return nil
}
