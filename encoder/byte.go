package encoder

import (
	"io"
)

// WriteByte writes a single byte to the writer.
func WriteByte(w io.Writer, c byte) error {
	if bw, ok := w.(io.ByteWriter); ok {
		return bw.WriteByte(c)
	}
	_, err := w.Write([]byte{c})
	return err
}

// WriteBytes writes a byte array to the writer.
func WriteBytes(w io.Writer, b []byte) error {
	_, err := w.Write(b)
	return err
}
