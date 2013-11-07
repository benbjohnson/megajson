package bench

import (
	"github.com/benbjohnson/megajson/encoder"
	"io"
)

type codeResponseJSONEncoder struct {
	enc encoder.Encoder
}

func NewcodeResponseJSONEncoder(w io.Writer) *codeResponseJSONEncoder {
	return &codeResponseJSONEncoder{enc: encoder.NewEncoder(w)}
}

func NewcodeResponseJSONRawEncoder(enc encoder.Encoder) *codeResponseJSONEncoder {
	return &codeResponseJSONEncoder{enc: enc}
}

func (e *codeResponseJSONEncoder) Encode(v *codeResponse) error {
	if err := e.RawEncode(v); err != nil {
		return err
	}
	if err := e.enc.Flush(); err != nil {
		return err
	}
	return nil
}
func (e *codeResponseJSONEncoder) RawEncode(v *codeResponse) error {
	if v == nil {
		return e.enc.WriteNull()
	}

	if err := e.enc.WriteByte('{'); err != nil {
		return err
	}
	if err := e.enc.WriteString("tree"); err != nil {
		return err
	}
	if err := e.enc.WriteByte(':'); err != nil {
		return err
	}
	if err := NewcodeNodeJSONRawEncoder(e.enc).RawEncode(v.Tree); err != nil {
		return err
	}
	if err := e.enc.WriteByte(','); err != nil {
		return err
	}
	if err := e.enc.WriteString("username"); err != nil {
		return err
	}
	if err := e.enc.WriteByte(':'); err != nil {
		return err
	}
	if err := e.enc.WriteString(v.Username); err != nil {
		return err
	}
	if err := e.enc.WriteByte('}'); err != nil {
		return err
	}
	return nil
}

type codeNodeJSONEncoder struct {
	enc encoder.Encoder
}

func NewcodeNodeJSONEncoder(w io.Writer) *codeNodeJSONEncoder {
	return &codeNodeJSONEncoder{enc: encoder.NewEncoder(w)}
}

func NewcodeNodeJSONRawEncoder(enc encoder.Encoder) *codeNodeJSONEncoder {
	return &codeNodeJSONEncoder{enc: enc}
}

func (e *codeNodeJSONEncoder) Encode(v *codeNode) error {
	if err := e.RawEncode(v); err != nil {
		return err
	}
	if err := e.enc.Flush(); err != nil {
		return err
	}
	return nil
}
func (e *codeNodeJSONEncoder) RawEncode(v *codeNode) error {
	if v == nil {
		return e.enc.WriteNull()
	}

	if err := e.enc.WriteByte('{'); err != nil {
		return err
	}
	if err := e.enc.WriteString("name"); err != nil {
		return err
	}
	if err := e.enc.WriteByte(':'); err != nil {
		return err
	}
	if err := e.enc.WriteString(v.Name); err != nil {
		return err
	}
	if err := e.enc.WriteByte(','); err != nil {
		return err
	}
	if err := e.enc.WriteString("kids"); err != nil {
		return err
	}
	if err := e.enc.WriteByte(':'); err != nil {
		return err
	}
	if err := e.enc.WriteByte('['); err != nil {
		return err
	}
	for index, elem := range v.Kids {
		if index > 0 {
			if err := e.enc.WriteByte(','); err != nil {
				return err
			}
		}
		if err := NewcodeNodeJSONRawEncoder(e.enc).RawEncode(elem); err != nil {
			return err
		}
	}
	if err := e.enc.WriteByte(']'); err != nil {
		return err
	}
	if err := e.enc.WriteByte(','); err != nil {
		return err
	}
	if err := e.enc.WriteString("cl_weight"); err != nil {
		return err
	}
	if err := e.enc.WriteByte(':'); err != nil {
		return err
	}
	if err := e.enc.WriteFloat64(v.CLWeight); err != nil {
		return err
	}
	if err := e.enc.WriteByte(','); err != nil {
		return err
	}
	if err := e.enc.WriteString("touches"); err != nil {
		return err
	}
	if err := e.enc.WriteByte(':'); err != nil {
		return err
	}
	if err := e.enc.WriteInt(v.Touches); err != nil {
		return err
	}
	if err := e.enc.WriteByte(','); err != nil {
		return err
	}
	if err := e.enc.WriteString("min_t"); err != nil {
		return err
	}
	if err := e.enc.WriteByte(':'); err != nil {
		return err
	}
	if err := e.enc.WriteInt64(v.MinT); err != nil {
		return err
	}
	if err := e.enc.WriteByte(','); err != nil {
		return err
	}
	if err := e.enc.WriteString("max_t"); err != nil {
		return err
	}
	if err := e.enc.WriteByte(':'); err != nil {
		return err
	}
	if err := e.enc.WriteInt64(v.MaxT); err != nil {
		return err
	}
	if err := e.enc.WriteByte(','); err != nil {
		return err
	}
	if err := e.enc.WriteString("mean_t"); err != nil {
		return err
	}
	if err := e.enc.WriteByte(':'); err != nil {
		return err
	}
	if err := e.enc.WriteInt64(v.MeanT); err != nil {
		return err
	}
	if err := e.enc.WriteByte('}'); err != nil {
		return err
	}
	return nil
}
