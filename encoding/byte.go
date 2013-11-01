package encoding

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
