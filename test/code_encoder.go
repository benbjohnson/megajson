package test

import (
	"github.com/benbjohnson/megajson/encoding"
	"io"
)

type codeResponseJSONEncoder struct {
	w io.Writer
}

func NewcodeResponseJSONEncoder(w io.Writer) *codeResponseJSONEncoder {
	return &codeResponseJSONEncoder{w: w}
}

func (e *codeResponseJSONEncoder) Encode(v *codeResponse) error {
	if err := encoding.WriteByte(e.w, '{'); err != nil {
		return err
	}
	if err := encoding.WriteString(e.w, "Username"); err != nil {
		return err
	}
	if err := encoding.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoding.WriteString(e.w, v.Username); err != nil {
		return err
	}
	if err := encoding.WriteByte(e.w, '}'); err != nil {
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
	if err := encoding.WriteByte(e.w, '{'); err != nil {
		return err
	}
	if err := encoding.WriteString(e.w, "Name"); err != nil {
		return err
	}
	if err := encoding.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoding.WriteString(e.w, v.Name); err != nil {
		return err
	}
	if err := encoding.WriteByte(e.w, ','); err != nil {
		return err
	}
	if err := encoding.WriteString(e.w, "Touches"); err != nil {
		return err
	}
	if err := encoding.WriteByte(e.w, ':'); err != nil {
		return err
	}
	if err := encoding.WriteInt(e.w, v.Touches); err != nil {
		return err
	}
	if err := encoding.WriteByte(e.w, '}'); err != nil {
		return err
	}
	return nil
}
