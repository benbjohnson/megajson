package encoding

import (
	"io"
	"strconv"
)

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
