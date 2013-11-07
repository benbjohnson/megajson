package bench

import (
	"github.com/benbjohnson/megajson/encoder"
	"io"
)

type codeResponseJSONEncoder struct {
	w io.Writer
}

func NewcodeResponseJSONEncoder(w io.Writer) *codeResponseJSONEncoder {
	return &codeResponseJSONEncoder{w: w}
}

func (e *codeResponseJSONEncoder) Encode(v *codeResponse) error {
	if v == nil {
		return encoder.WriteNull(e.w)
	}

	if err := encoder.WriteByte(e.w, '{'); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, "tree"); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := NewcodeNodeJSONEncoder(e.w).Encode(v.Tree); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ','); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, "username"); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, v.Username); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, '}'); err != nil {
		return err
	}
	return nil
}

type codeNodeJSONEncoder struct {
	w io.Writer
}

func NewcodeNodeJSONEncoder(w io.Writer) *codeNodeJSONEncoder {
	return &codeNodeJSONEncoder{w: w}
}

func (e *codeNodeJSONEncoder) Encode(v *codeNode) error {
	if v == nil {
		return encoder.WriteNull(e.w)
	}

	if err := encoder.WriteByte(e.w, '{'); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, "name"); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, v.Name); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ','); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, "kids"); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, '['); err != nil {
		return err
	}
	for index, elem := range v.Kids {
		if index > 0 {
			if err := encoder.WriteByte(e.w, ','); err != nil {
				return err
			}
		}
		if err := NewcodeNodeJSONEncoder(e.w).Encode(elem); err != nil {
			return err
		}
	}
	if err := encoder.WriteByte(e.w, ']'); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ','); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, "cl_weight"); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoder.WriteFloat64(e.w, v.CLWeight); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ','); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, "touches"); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoder.WriteInt(e.w, v.Touches); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ','); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, "min_t"); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoder.WriteInt64(e.w, v.MinT); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ','); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, "max_t"); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoder.WriteInt64(e.w, v.MaxT); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ','); err != nil {
		return err
	}
	if err := encoder.WriteString(e.w, "mean_t"); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoder.WriteInt64(e.w, v.MeanT); err != nil {
		return err
	}
	if err := encoder.WriteByte(e.w, '}'); err != nil {
		return err
	}
	return nil
}
