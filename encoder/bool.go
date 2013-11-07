package encoder

import (
	"io"
)

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
