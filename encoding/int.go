package encoding

import (
	"io"
	"strconv"
)

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


