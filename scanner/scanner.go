package scanner

import (
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

const (
	// The size, in bytes, that is read from the reader at a time.
	bufSize = 4096
)

// Scanner is a tokenizer for JSON input from an io.Reader.
type Scanner interface {
	Pos() int
	Scan() (int, []byte, error)
	Unscan(tok int, b []byte)
	ReadString(target *string) error
	ReadInt(target *int) error
	ReadInt64(target *int64) error
	ReadUint(target *uint) error
	ReadUint64(target *uint64) error
	ReadFloat32(target *float32) error
	ReadFloat64(target *float64) error
	ReadBool(target *bool) error
	ReadMap(target *map[string]interface{}) error
}

type scanner struct {
	r       io.Reader
	c       rune
	scratch [bufSize]byte
	buf     [bufSize]byte
	buflen  int
	idx     int
	pos     int
	tmpc    rune
	tmp     struct {
		tok int
		b   []byte
		err error
	}
}

// NewScanner initializes a new scanner with a given reader.
func NewScanner(r io.Reader) Scanner {
	s := &scanner{r: r, buflen: -1}
	return s
}

// Pos returns the current rune position of the scanner.
func (s *scanner) Pos() int {
	return s.pos
}

// read retrieves the next rune from the reader.
func (s *scanner) read() error {
	if s.tmpc > 0 {
		s.c = s.tmpc
		s.tmpc = 0
		return nil
	}

	// Read from the reader if the buffer is empty.
	if s.idx >= s.buflen {
		var err error
		if s.buflen, err = s.r.Read(s.buf[0:]); err != nil {
			return err
		}
		s.idx = 0
	}

	// Read a single byte and then determine if utf8 decoding is needed.
	b := s.buf[s.idx]
	if b < utf8.RuneSelf {
		s.c = rune(b)
		s.idx++
	} else {
		// Read a new buffer if we don't have at least the max size of a UTF8 character.
		if s.idx+utf8.UTFMax >= s.buflen {
			s.buf[0] = b
			var err error
			if s.buflen, err = s.r.Read(s.buf[1:]); err != nil {
				return err
			}
			s.buflen += 1
		}

		var size int
		s.c, size = utf8.DecodeRune(s.buf[s.idx:])
		s.idx += size
	}

	s.pos++
	return nil
}

// unread places the current rune back on the reader.
func (s *scanner) unread() {
	s.tmpc = s.c
}

// expect reads the next rune and checks that it matches.
func (s *scanner) expect(c rune) error {
	if err := s.read(); err != nil {
		return err
	} else if s.c != c {
		return fmt.Errorf("Unexpected char: %q", s.c)
	}
	return nil
}

// Scan returns the next JSON token from the reader.
func (s *scanner) Scan() (int, []byte, error) {
	if s.tmp.tok != 0 {
		tok, b := s.tmp.tok, s.tmp.b
		s.tmp.tok, s.tmp.b = 0, nil
		return tok, b, nil
	}

	for {
		if err := s.read(); err != nil {
			return 0, nil, err
		}

		switch s.c {
		case '{':
			return TLBRACE, []byte{'{'}, nil
		case '}':
			return TRBRACE, []byte{'}'}, nil
		case '[':
			return TLBRACKET, []byte{'['}, nil
		case ']':
			return TRBRACKET, []byte{']'}, nil
		case ':':
			return TCOLON, []byte{':'}, nil
		case ',':
			return TCOMMA, []byte{','}, nil
		case '"':
			return s.scanString()
		case 't':
			return s.scanTrue()
		case 'f':
			return s.scanFalse()
		case 'n':
			return s.scanNull()
		}

		if (s.c >= '0' && s.c <= '9') || s.c == '-' {
			return s.scanNumber()
		}
	}
}

// Unscan adds a token and byte array back onto the buffer to be read
// on the next call to Scan().
func (s *scanner) Unscan(tok int, b []byte) {
	s.tmp.tok = tok
	s.tmp.b = b
	s.pos--
}

// scanNumber reads a JSON number from the reader.
func (s *scanner) scanNumber() (int, []byte, error) {
	var n int

	if s.c == '-' {
		s.scratch[n] = '-'
		n++
		if err := s.read(); err != nil {
			return 0, nil, err
		}
	}

	// Read whole number.
	if err := s.scanDigits(&n); err == io.EOF {
		return TNUMBER, s.scratch[0:n], nil
	} else if err != nil {
		return 0, nil, err
	}
	n++

	// Read period.
	if err := s.read(); err == io.EOF {
		return TNUMBER, s.scratch[0:n], nil
	} else if err != nil {
		return 0, nil, err
	} else if s.c != '.' {
		s.unread()
		return TNUMBER, s.scratch[0:n], nil
	}
	s.scratch[n] = '.'
	n++

	if err := s.read(); err != nil {
		return 0, nil, err
	}

	// Read fraction.
	if err := s.scanDigits(&n); err == io.EOF {
		return TNUMBER, s.scratch[0:n], nil
	} else if err != nil {
		return 0, nil, err
	}
	n++

	return TNUMBER, s.scratch[0:n], nil
}

// scanDigits reads a series of digits from the reader.
func (s *scanner) scanDigits(n *int) error {
	for {
		if s.c >= '0' && s.c <= '9' {
			s.scratch[*n] = byte(s.c)
			(*n)++
			if err := s.read(); err != nil {
				return err
			}
		} else {
			s.unread()
			(*n)--
			return nil
		}
	}
}

// scanString reads a quoted JSON string from the reader.
func (s *scanner) scanString() (int, []byte, error) {
	// TODO: Support large strings (e.g. >bufSize).
	var overflow []byte

	var n int
	for {
		if err := s.read(); err != nil {
			return 0, nil, err
		}
		switch s.c {
		case '\\':
			if err := s.read(); err != nil {
				return 0, nil, err
			}
			switch s.c {
			case '"':
				s.scratch[n] = '"'
				n++
			case '\\':
				s.scratch[n] = '\\'
				n++
			case '/':
				s.scratch[n] = '/'
				n++
			case 'b':
				s.scratch[n] = '\b'
				n++
			case 'f':
				s.scratch[n] = '\f'
				n++
			case 'n':
				s.scratch[n] = '\n'
				n++
			case 'r':
				s.scratch[n] = '\r'
				n++
			case 't':
				s.scratch[n] = '\t'
				n++
			case 'u':
				numeric := make([]byte, 4)
				numericCount := 0
			unicode_loop:
				for {
					if err := s.read(); err != nil {
						return 0, nil, err
					}
					switch {
					case s.c >= '0' && s.c <= '9' || s.c >= 'a' && s.c <= 'f' || s.c >= 'A' && s.c <= 'F':
						numeric[numericCount] = byte(s.c)
						numericCount++
						if numericCount == 4 {
							var i int64
							var err error
							if i, err = strconv.ParseInt(string(numeric), 16, 32); err != nil {
								return 0, nil, err
							}
							if i < utf8.RuneSelf {
								s.scratch[n] = byte(i)
								n++
							} else {
								encoded := utf8.EncodeRune(s.scratch[n:], rune(i))
								n += encoded
							}
							break unicode_loop
						}
					default:
						s.unread()
						return 0, nil, fmt.Errorf("Unexpected symbol in unicode escape: %c", s.c)
					}
				}
			default:
				return 0, nil, fmt.Errorf("Invalid escape character: \\%c", s.c)
			}

		case '"':
			if len(overflow) == 0 {
				return TSTRING, s.scratch[0:n], nil
			}
			overflow = append(overflow, s.scratch[0:n]...)
			return TSTRING, overflow, nil

		default:
			if s.c < utf8.RuneSelf {
				if n == bufSize {
					overflow = append(overflow, s.scratch[0:n]...)
					n = 0
				}
				s.scratch[n] = byte(s.c)
				n++
			} else {
				n += utf8.EncodeRune(s.scratch[n:], s.c)
			}
		}
	}
}

// scanTrue reads the "true" token.
func (s *scanner) scanTrue() (int, []byte, error) {
	if err := s.expect('r'); err != nil {
		return 0, nil, err
	}
	if err := s.expect('u'); err != nil {
		return 0, nil, err
	}
	if err := s.expect('e'); err != nil {
		return 0, nil, err
	}
	return TTRUE, nil, nil
}

// scanFalse reads the "false" token.
func (s *scanner) scanFalse() (int, []byte, error) {
	if err := s.expect('a'); err != nil {
		return 0, nil, err
	}
	if err := s.expect('l'); err != nil {
		return 0, nil, err
	}
	if err := s.expect('s'); err != nil {
		return 0, nil, err
	}
	if err := s.expect('e'); err != nil {
		return 0, nil, err
	}
	return TFALSE, nil, nil
}

// scanNull reads the "null" token.
func (s *scanner) scanNull() (int, []byte, error) {
	if err := s.expect('u'); err != nil {
		return 0, nil, err
	}
	if err := s.expect('l'); err != nil {
		return 0, nil, err
	}
	if err := s.expect('l'); err != nil {
		return 0, nil, err
	}
	return TNULL, nil, nil
}

// ReadString reads a token into a string variable.
func (s *scanner) ReadString(target *string) error {
	tok, b, err := s.Scan()
	if err != nil {
		return err
	}
	switch tok {
	case TSTRING:
		*target = string(b)
	case TNUMBER, TTRUE, TFALSE, TNULL:
		*target = ""
	default:
		return fmt.Errorf("Unexpected %s at %d: %s; expected string", TokenName(tok), s.pos, string(b))
	}
	return nil
}

// ReadInt reads a token into an int variable.
func (s *scanner) ReadInt(target *int) error {
	tok, b, err := s.Scan()
	if err != nil {
		return err
	}
	switch tok {
	case TNUMBER:
		n, _ := strconv.ParseInt(string(b), 10, 64)
		*target = int(n)
	case TSTRING, TTRUE, TFALSE, TNULL:
		*target = 0
	default:
		return fmt.Errorf("Unexpected %s at %d: %s; expected number", TokenName(tok), s.pos, string(b))
	}
	return nil
}

// ReadInt64 reads a token into an int64 variable.
func (s *scanner) ReadInt64(target *int64) error {
	tok, b, err := s.Scan()
	if err != nil {
		return err
	}
	switch tok {
	case TNUMBER:
		n, _ := strconv.ParseInt(string(b), 10, 64)
		*target = n
	case TSTRING, TTRUE, TFALSE, TNULL:
		*target = 0
	default:
		return fmt.Errorf("Unexpected %s at %d: %s; expected number", TokenName(tok), s.pos, string(b))
	}
	return nil
}

// ReadUint reads a token into an uint variable.
func (s *scanner) ReadUint(target *uint) error {
	tok, b, err := s.Scan()
	if err != nil {
		return err
	}
	switch tok {
	case TNUMBER:
		n, _ := strconv.ParseUint(string(b), 10, 64)
		*target = uint(n)
	case TSTRING, TTRUE, TFALSE, TNULL:
		*target = 0
	default:
		return fmt.Errorf("Unexpected %s at %d: %s; expected number", TokenName(tok), s.pos, string(b))
	}
	return nil
}

// ReadUint64 reads a token into an uint64 variable.
func (s *scanner) ReadUint64(target *uint64) error {
	tok, b, err := s.Scan()
	if err != nil {
		return err
	}
	switch tok {
	case TNUMBER:
		n, _ := strconv.ParseUint(string(b), 10, 64)
		*target = n
	case TSTRING, TTRUE, TFALSE, TNULL:
		*target = 0
	default:
		return fmt.Errorf("Unexpected %s at %d: %s; expected number", TokenName(tok), s.pos, string(b))
	}
	return nil
}

// ReadFloat32 reads a token into a float32 variable.
func (s *scanner) ReadFloat32(target *float32) error {
	tok, b, err := s.Scan()
	if err != nil {
		return err
	}
	switch tok {
	case TNUMBER:
		n, _ := strconv.ParseFloat(string(b), 32)
		*target = float32(n)
	case TSTRING, TTRUE, TFALSE, TNULL:
		*target = 0
	default:
		return fmt.Errorf("Unexpected %s at %d: %s; expected number", TokenName(tok), s.pos, string(b))
	}
	return nil
}

// ReadFloat64 reads a token into a float64 variable.
func (s *scanner) ReadFloat64(target *float64) error {
	tok, b, err := s.Scan()
	if err != nil {
		return err
	}
	switch tok {
	case TNUMBER:
		n, _ := strconv.ParseFloat(string(b), 64)
		*target = n
	case TSTRING, TTRUE, TFALSE, TNULL:
		*target = 0
	default:
		return fmt.Errorf("Unexpected %s at %d: %s; expected number", TokenName(tok), s.pos, string(b))
	}
	return nil
}

// ReadBool reads a token into a boolean variable.
func (s *scanner) ReadBool(target *bool) error {
	tok, b, err := s.Scan()
	if err != nil {
		return err
	}
	switch tok {
	case TTRUE:
		*target = true
	case TFALSE, TSTRING, TNUMBER, TNULL:
		*target = false
	default:
		return fmt.Errorf("Unexpected %s at %d: %s; expected number", TokenName(tok), s.pos, string(b))
	}
	return nil
}

// ReadMap reads the next value into a map variable.
func (s *scanner) ReadMap(target *map[string]interface{}) error {
	if tok, b, err := s.Scan(); err != nil {
		return err
	} else if tok == TNULL {
		*target = nil
		return nil
	} else if tok != TLBRACE {
		return fmt.Errorf("Unexpected %s at %d: %s; expected '{'", TokenName(tok), s.Pos(), string(b))
	}

	// Create a new map.
	*target = make(map[string]interface{})
	v := *target

	// Loop over key/value pairs.
	index := 0
	for {
		// Read in key.
		var key string
		tok, b, err := s.Scan()
		if err != nil {
			return err
		} else if tok == TRBRACE {
			return nil
		} else if tok == TCOMMA {
			if index == 0 {
				return fmt.Errorf("Unexpected comma at %d", s.Pos())
			}
			if tok, b, err = s.Scan(); err != nil {
				return err
			}
		}

		if tok != TSTRING {
			return fmt.Errorf("Unexpected %s at %d: %s; expected '{' or string", TokenName(tok), s.Pos(), string(b))
		} else {
			key = string(b)
		}

		// Read in the colon.
		if tok, b, err := s.Scan(); err != nil {
			return err
		} else if tok != TCOLON {
			return fmt.Errorf("Unexpected %s at %d: %s; expected colon", TokenName(tok), s.Pos(), string(b))
		}

		// Read the next token.
		tok, b, err = s.Scan()
		if err != nil {
			return err
		}
		switch tok {
		case TSTRING:
			v[key] = string(b)
		case TNUMBER:
			v[key], _ = strconv.ParseFloat(string(b), 64)
		case TTRUE:
			v[key] = true
		case TFALSE:
			v[key] = false
		case TNULL:
			v[key] = nil
		case TLBRACE:
			s.Unscan(tok, b)
			m := make(map[string]interface{})
			if err := s.ReadMap(&m); err != nil {
				return err
			}
			v[key] = m
		case TLBRACKET:
			// TODO: Read arrays.
			panic("Not supported yet")
		default:
			return fmt.Errorf("Unexpected %s at %d: %s", TokenName(tok), s.Pos(), string(b))
		}

		index++
	}

	return nil
}
