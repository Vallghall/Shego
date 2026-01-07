package lex

// TKind represents the type of a token.
type TKind int

const (
	ParenOpen TKind = iota
	ParenClose
	String
	Number
	Atom
	Quote    // '
	Backtick // `
	Comma    // ,
	CommaAt  // ,@
)

// String returns a human-readable name for the token kind.
func (k TKind) String() string {
	switch k {
	case ParenOpen:
		return "ParenOpen"
	case ParenClose:
		return "ParenClose"
	case String:
		return "String"
	case Number:
		return "Number"
	case Atom:
		return "Atom"
	case Quote:
		return "Quote"
	case Backtick:
		return "Backtick"
	case Comma:
		return "Comma"
	case CommaAt:
		return "CommaAt"
	default:
		return "Unknown"
	}
}

// Token represents a lexical token from the source code.
type Token interface {
	File() string
	Position() int     // absolute position in input
	Line() int         // 1-based line number
	LinePosition() int // 1-based column within line
	Raw() string       // original string the token was parsed from
	Kind() TKind
}

// token is the concrete implementation of the Token interface.
type token struct {
	file    string
	pos     int
	line    int
	linePos int
	raw     string
	kind    TKind
}

func (t *token) File() string      { return t.file }
func (t *token) Position() int     { return t.pos }
func (t *token) Line() int         { return t.line }
func (t *token) LinePosition() int { return t.linePos }
func (t *token) Raw() string       { return t.raw }
func (t *token) Kind() TKind       { return t.kind }

// NewToken creates a new token with the given properties.
func NewToken(file string, pos, line, linePos int, raw string, kind TKind) Token {
	return &token{
		file:    file,
		pos:     pos,
		line:    line,
		linePos: linePos,
		raw:     raw,
		kind:    kind,
	}
}
