package encoding

import (
	"io"
	"strconv"
)

// WriteInt encodes and writes an integer to a writer.
func WriteInt(w io.Writer, v int) error {
	return WriteInt64(w, int64(v))
}

// WriteInt64 encodes and writes a 64-bit integer to a writer.
func WriteInt64(w io.Writer, v int64) error {
	var b [64]byte
	_, err := w.Write(strconv.AppendInt(b[:0], v, 10))
	return err
}

// WriteUint encodes and writes an unsigned integer to a writer.
func WriteUint(w io.Writer, v uint) error {
	return WriteUint64(w, uint64(v))
}

// WriteUint encodes and writes an unsigned integer to a writer.
func WriteUint64(w io.Writer, v uint64) error {
	var b [64]byte
	_, err := w.Write(strconv.AppendUint(b[:0], v, 10))
	return err
}


