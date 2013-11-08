package scanner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
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
}

type scanner struct {
	r *bufio.Reader
	c rune
	pos int
	tmp struct{
		tok int
		b   []byte
		err error
	}
}

// NewScanner initializes a new scanner with a given reader.
func NewScanner(r io.Reader) Scanner {
	s := &scanner{r:bufio.NewReader(r)}
	return s
}

// Pos returns the current character position of the scanner.
func (s *scanner) Pos() int {
	return s.pos
}

// read retrieves the next rune from the reader.
func (s *scanner) read() error {
	var err error
	s.c, _, err = s.r.ReadRune()
	if err == nil {
		s.pos++
	}
	return err
}

// unread places the current rune back on the reader.
func (s *scanner) unread() error {
	return s.r.UnreadRune()
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
	var b bytes.Buffer
	if s.c == '-' {
		b.WriteRune('-')
		if err := s.read(); err != nil {
			return 0, nil, err
		}
	}

	// Read whole number.
	if _, err := s.scanDigits(&b); err == io.EOF {
		return TNUMBER, b.Bytes(), nil
	} else if err != nil {
		return 0, nil, err
	}

	// Read period.
	if err := s.read(); err != nil {
		return 0, nil, err
	} else if s.c != '.' {
		if err := s.unread(); err != nil {
			return 0, nil, err
		}
		return TNUMBER, b.Bytes(), nil
	}
	b.WriteByte('.')

	if err := s.read(); err != nil {
		return 0, nil, err
	}

	// Read fraction.
	if _, err := s.scanDigits(&b); err == io.EOF {
		return TNUMBER, b.Bytes(), nil
	} else if err != nil {
		return 0, nil, err
	}

	return TNUMBER, b.Bytes(), nil
}

// scanDigits reads a series of digits from the reader.
func (s *scanner) scanDigits(b *bytes.Buffer) (int, error) {
	count := 0
	for {
		if s.c >= '0' && s.c <= '9' {
			b.WriteRune(s.c)
			if err := s.read(); err == io.EOF {
				return count, err
			} else if err != nil {
				return 0, err
			}
			count++
		} else {
			if err := s.unread(); err != nil {
				return 0, err
			}
			return count-1, nil
		}
	}
}

// scanString reads a quoted JSON string from the reader.
func (s *scanner) scanString() (int, []byte, error) {
	var b bytes.Buffer

	for {
		if err := s.read(); err != nil {
			return 0, nil, err
		}
		switch s.c {
		case '\\':
			if err := s.scanEscape(&b); err != nil {
				return 0, nil, err
			}
		case '"':
			return TSTRING, b.Bytes(), nil
		default:
			b.WriteRune(s.c)
		}
	}
}

// scanEscape reads an escaped string character.
func (s *scanner) scanEscape(b *bytes.Buffer) error {
	for {
		if err := s.read(); err != nil {
			return err
		}
		switch s.c {
		case '"':
			return b.WriteByte('"')
		case '\\':
			return b.WriteByte('\\')
		case '/':
			return b.WriteByte('/')
		case 'b':
			return b.WriteByte('\b')
		case 'f':
			return b.WriteByte('\f')
		case 'n':
			return b.WriteByte('\n')
		case 'r':
			return b.WriteByte('\r')
		case 't':
			return b.WriteByte('\t')
		case 'u':
			// TODO: \u0000
			return nil
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

