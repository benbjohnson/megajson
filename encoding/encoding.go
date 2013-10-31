package encoding

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
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
         return &json.InvalidUTF8Error{S:v}
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

// WriteInt encodes and writes an integer to a writer.
func WriteInt(w io.Writer, v int) error {
   var b [64]byte
   _, err := w.Write(strconv.AppendInt(b[:0], int64(v), 10))
   return err
}

// WriteUint encodes and writes an unsigned integer to a writer.
func WriteUint(w io.Writer, v uint) error {
   var b [64]byte
   _, err := w.Write(strconv.AppendUint(b[:0], uint64(v), 10))
   return err
}

// WriteFloat32 encodes and writes a 32-bit float to a writer.
func WriteFloat32(w io.Writer, v float32) error {
   var b [64]byte
   _, err := w.Write(strconv.AppendFloat(b[:0], float64(v), 'g', -1, 32))
   return err
}

// WriteFloat64 encodes and writes a 64-bit float to a writer.
func WriteFloat64(w io.Writer, v float64) error {
   var b [64]byte
   _, err := w.Write(strconv.AppendFloat(b[:0], float64(v), 'g', -1, 64))
   return err
}

// WriteBool encodes and writes a boolean value to a writer.
func WriteBool(w io.Writer, v bool) error {
   if v {
      _, err := w.Write([]byte("true"))
      return err
   } else {
      _, err := w.Write([]byte("false"))
      return err
   }
}

// writeByte writes a single byte to the writer.
func writeByte(w io.Writer, c byte) error {
	if bw, ok := w.(io.ByteWriter); ok {
		return bw.WriteByte(c)
	}
	_, err := w.Write([]byte{c})
	return err
}
