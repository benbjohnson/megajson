package scanner

const (
	TSTRING = iota + 1
	TNUMBER
	TTRUE
	TFALSE
	TNULL
	TLBRACE
	TRBRACE
	TLBRACKET
	TRBRACKET
	TCOLON
	TCOMMA
)

var names = map[int]string{
	TSTRING:   "string",
	TNUMBER:   "number",
	TTRUE:     "true",
	TFALSE:    "false",
	TNULL:     "null",
	TLBRACE:   "left brace",
	TRBRACE:   "right brace",
	TLBRACKET: "left bracket",
	TRBRACKET: "right bracket",
	TCOLON:    "colon",
	TCOMMA:    "comma",
}

// TokenName returns a human readable version of the token.
func TokenName(tok int) string {
	return names[tok]
}
