package test

import (
	"errors"
	"fmt"
	"github.com/benbjohnson/megajson/scanner"
	"io"
)

type codeResponseJSONDecoder struct {
	s scanner.Scanner
}

func NewcodeResponseJSONDecoder(r io.Reader) *codeResponseJSONDecoder {
	return &codeResponseJSONDecoder{s: scanner.NewScanner(r)}
}

func NewcodeResponseJSONScanDecoder(s scanner.Scanner) *codeResponseJSONDecoder {
	return &codeResponseJSONDecoder{s: s}
}

func (e *codeResponseJSONDecoder) Decode(ptr **codeResponse) error {
	s := e.s
	if tok, tokval, err := s.Scan(); err != nil {
		return err
	} else if tok == scanner.TNULL {
		*ptr = nil
		return nil
	} else if tok != scanner.TLBRACE {
		return fmt.Errorf("Unexpected %s at %d: %s; expected '{'", scanner.TokenName(tok), s.Pos(), string(tokval))
	}

	// Create the object if it doesn't exist.
	if *ptr == nil {
		*ptr = &codeResponse{}
	}
	v := *ptr

	// Loop over key/value pairs.
	index := 0
	for {
		// Read in key.
		var key string
		tok, tokval, err := s.Scan()
		if err != nil {
			return err
		} else if tok == scanner.TRBRACE {
			return nil
		} else if tok == scanner.TCOMMA {
			if index == 0 {
				return fmt.Errorf("Unexpected comma at %d", s.Pos())
			}
			if tok, tokval, err = s.Scan(); err != nil {
				return err
			}
		}

		if tok != scanner.TSTRING {
			return fmt.Errorf("Unexpected %s at %d: %s; expected '{' or string", scanner.TokenName(tok), s.Pos(), string(tokval))
		} else {
			key = string(tokval)
		}

		// Read in the colon.
		if tok, tokval, err := s.Scan(); err != nil {
			return err
		} else if tok != scanner.TCOLON {
			return fmt.Errorf("Unexpected %s at %d: %s; expected colon", scanner.TokenName(tok), s.Pos(), string(tokval))
		}

		switch key {

		case "tree":
			v := &v.Tree

			if err := NewcodeNodeJSONScanDecoder(s).Decode(v); err != nil {
				return err
			}

		case "username":
			v := &v.Username

			if err := s.ReadString(v); err != nil {
				return err
			}

		}

		index++
	}

	return nil
}

func (e *codeResponseJSONDecoder) DecodeArray(ptr *[]*codeResponse) error {
	s := e.s
	if tok, _, err := s.Scan(); err != nil {
		return err
	} else if tok != scanner.TLBRACKET {
		return errors.New("Expected '['")
	}

	slice := make([]*codeResponse, 0)

	// Loop over items.
	index := 0
	for {
		tok, tokval, err := s.Scan()
		if err != nil {
			return err
		} else if tok == scanner.TRBRACKET {
			*ptr = slice
			return nil
		} else if tok == scanner.TCOMMA {
			if index == 0 {
				return fmt.Errorf("Unexpected comma in array at %d", s.Pos())
			}
			if tok, tokval, err = s.Scan(); err != nil {
				return err
			}
		}
		s.Unscan(tok, tokval)

		item := &codeResponse{}
		if err := e.Decode(&item); err != nil {
			return err
		}
		slice = append(slice, item)

		index++
	}
}

type codeNodeJSONDecoder struct {
	s scanner.Scanner
}

func NewcodeNodeJSONDecoder(r io.Reader) *codeNodeJSONDecoder {
	return &codeNodeJSONDecoder{s: scanner.NewScanner(r)}
}

func NewcodeNodeJSONScanDecoder(s scanner.Scanner) *codeNodeJSONDecoder {
	return &codeNodeJSONDecoder{s: s}
}

func (e *codeNodeJSONDecoder) Decode(ptr **codeNode) error {
	s := e.s
	if tok, tokval, err := s.Scan(); err != nil {
		return err
	} else if tok == scanner.TNULL {
		*ptr = nil
		return nil
	} else if tok != scanner.TLBRACE {
		return fmt.Errorf("Unexpected %s at %d: %s; expected '{'", scanner.TokenName(tok), s.Pos(), string(tokval))
	}

	// Create the object if it doesn't exist.
	if *ptr == nil {
		*ptr = &codeNode{}
	}
	v := *ptr

	// Loop over key/value pairs.
	index := 0
	for {
		// Read in key.
		var key string
		tok, tokval, err := s.Scan()
		if err != nil {
			return err
		} else if tok == scanner.TRBRACE {
			return nil
		} else if tok == scanner.TCOMMA {
			if index == 0 {
				return fmt.Errorf("Unexpected comma at %d", s.Pos())
			}
			if tok, tokval, err = s.Scan(); err != nil {
				return err
			}
		}

		if tok != scanner.TSTRING {
			return fmt.Errorf("Unexpected %s at %d: %s; expected '{' or string", scanner.TokenName(tok), s.Pos(), string(tokval))
		} else {
			key = string(tokval)
		}

		// Read in the colon.
		if tok, tokval, err := s.Scan(); err != nil {
			return err
		} else if tok != scanner.TCOLON {
			return fmt.Errorf("Unexpected %s at %d: %s; expected colon", scanner.TokenName(tok), s.Pos(), string(tokval))
		}

		switch key {

		case "name":
			v := &v.Name

			if err := s.ReadString(v); err != nil {
				return err
			}

		case "kids":
			v := &v.Kids

			if err := NewcodeNodeJSONScanDecoder(s).DecodeArray(v); err != nil {
				return err
			}

		case "cl_weight":
			v := &v.CLWeight

			if err := s.ReadFloat64(v); err != nil {
				return err
			}

		case "touches":
			v := &v.Touches

			if err := s.ReadInt(v); err != nil {
				return err
			}

		case "min_t":
			v := &v.MinT

			if err := s.ReadInt64(v); err != nil {
				return err
			}

		case "max_t":
			v := &v.MaxT

			if err := s.ReadInt64(v); err != nil {
				return err
			}

		case "mean_t":
			v := &v.MeanT

			if err := s.ReadInt64(v); err != nil {
				return err
			}

		}

		index++
	}

	return nil
}

func (e *codeNodeJSONDecoder) DecodeArray(ptr *[]*codeNode) error {
	s := e.s
	if tok, _, err := s.Scan(); err != nil {
		return err
	} else if tok != scanner.TLBRACKET {
		return errors.New("Expected '['")
	}

	slice := make([]*codeNode, 0)

	// Loop over items.
	index := 0
	for {
		tok, tokval, err := s.Scan()
		if err != nil {
			return err
		} else if tok == scanner.TRBRACKET {
			*ptr = slice
			return nil
		} else if tok == scanner.TCOMMA {
			if index == 0 {
				return fmt.Errorf("Unexpected comma in array at %d", s.Pos())
			}
			if tok, tokval, err = s.Scan(); err != nil {
				return err
			}
		}
		s.Unscan(tok, tokval)

		item := &codeNode{}
		if err := e.Decode(&item); err != nil {
			return err
		}
		slice = append(slice, item)

		index++
	}
}
