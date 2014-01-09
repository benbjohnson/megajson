package bench

import (
	"github.com/benbjohnson/megajson/writer"
	"io"
)

type codeResponseJSONEncoder struct {
	w *writer.Writer
}

func NewcodeResponseJSONEncoder(w io.Writer) *codeResponseJSONEncoder {
	return &codeResponseJSONEncoder{w: writer.NewWriter(w)}
}

func NewcodeResponseJSONRawEncoder(w *writer.Writer) *codeResponseJSONEncoder {
	return &codeResponseJSONEncoder{w: w}
}

func (e *codeResponseJSONEncoder) Encode(v *codeResponse) error {
	if err := e.RawEncode(v); err != nil {
		return err
	}
	if err := e.w.Flush(); err != nil {
		return err
	}
	return nil
}

func (e *codeResponseJSONEncoder) RawEncode(v *codeResponse) error {
	if v == nil {
		return e.w.WriteNull()
	}

	if err := e.w.WriteByte('{'); err != nil {
		return err
	}

	// Write key and colon.
	if err := e.w.WriteString("tree"); err != nil {
		return err
	}
	if err := e.w.WriteByte(':'); err != nil {
		return err
	}

	// Write value.
	{
		v := v.Tree

		if err := NewcodeNodeJSONRawEncoder(e.w).RawEncode(v); err != nil {
			return err
		}

	}

	if err := e.w.WriteByte(','); err != nil {
		return err
	}

	// Write key and colon.
	if err := e.w.WriteString("username"); err != nil {
		return err
	}
	if err := e.w.WriteByte(':'); err != nil {
		return err
	}

	// Write value.
	{
		v := v.Username

		if err := e.w.WriteString(v); err != nil {
			return err
		}

	}

	if err := e.w.WriteByte('}'); err != nil {
		return err
	}
	return nil
}

type codeNodeJSONEncoder struct {
	w *writer.Writer
}

func NewcodeNodeJSONEncoder(w io.Writer) *codeNodeJSONEncoder {
	return &codeNodeJSONEncoder{w: writer.NewWriter(w)}
}

func NewcodeNodeJSONRawEncoder(w *writer.Writer) *codeNodeJSONEncoder {
	return &codeNodeJSONEncoder{w: w}
}

func (e *codeNodeJSONEncoder) Encode(v *codeNode) error {
	if err := e.RawEncode(v); err != nil {
		return err
	}
	if err := e.w.Flush(); err != nil {
		return err
	}
	return nil
}

func (e *codeNodeJSONEncoder) RawEncode(v *codeNode) error {
	if v == nil {
		return e.w.WriteNull()
	}

	if err := e.w.WriteByte('{'); err != nil {
		return err
	}

	// Write key and colon.
	if err := e.w.WriteString("name"); err != nil {
		return err
	}
	if err := e.w.WriteByte(':'); err != nil {
		return err
	}

	// Write value.
	{
		v := v.Name

		if err := e.w.WriteString(v); err != nil {
			return err
		}

	}

	if err := e.w.WriteByte(','); err != nil {
		return err
	}

	// Write key and colon.
	if err := e.w.WriteString("kids"); err != nil {
		return err
	}
	if err := e.w.WriteByte(':'); err != nil {
		return err
	}

	// Write value.
	{
		v := v.Kids

		if err := e.w.WriteByte('['); err != nil {
			return err
		}

		for index, v := range v {
			if index > 0 {
				if err := e.w.WriteByte(','); err != nil {
					return err
				}
			}
			if err := NewcodeNodeJSONRawEncoder(e.w).RawEncode(v); err != nil {
				return err
			}
		}

		if err := e.w.WriteByte(']'); err != nil {
			return err
		}

	}

	if err := e.w.WriteByte(','); err != nil {
		return err
	}

	// Write key and colon.
	if err := e.w.WriteString("cl_weight"); err != nil {
		return err
	}
	if err := e.w.WriteByte(':'); err != nil {
		return err
	}

	// Write value.
	{
		v := v.CLWeight

		if err := e.w.WriteFloat64(v); err != nil {
			return err
		}

	}

	if err := e.w.WriteByte(','); err != nil {
		return err
	}

	// Write key and colon.
	if err := e.w.WriteString("touches"); err != nil {
		return err
	}
	if err := e.w.WriteByte(':'); err != nil {
		return err
	}

	// Write value.
	{
		v := v.Touches

		if err := e.w.WriteInt(v); err != nil {
			return err
		}

	}

	if err := e.w.WriteByte(','); err != nil {
		return err
	}

	// Write key and colon.
	if err := e.w.WriteString("min_t"); err != nil {
		return err
	}
	if err := e.w.WriteByte(':'); err != nil {
		return err
	}

	// Write value.
	{
		v := v.MinT

		if err := e.w.WriteInt64(v); err != nil {
			return err
		}

	}

	if err := e.w.WriteByte(','); err != nil {
		return err
	}

	// Write key and colon.
	if err := e.w.WriteString("max_t"); err != nil {
		return err
	}
	if err := e.w.WriteByte(':'); err != nil {
		return err
	}

	// Write value.
	{
		v := v.MaxT

		if err := e.w.WriteInt64(v); err != nil {
			return err
		}

	}

	if err := e.w.WriteByte(','); err != nil {
		return err
	}

	// Write key and colon.
	if err := e.w.WriteString("mean_t"); err != nil {
		return err
	}
	if err := e.w.WriteByte(':'); err != nil {
		return err
	}

	// Write value.
	{
		v := v.MeanT

		if err := e.w.WriteInt64(v); err != nil {
			return err
		}

	}

	if err := e.w.WriteByte('}'); err != nil {
		return err
	}
	return nil
}
