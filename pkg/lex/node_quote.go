package lex

// QuoteLexNode handles quote-related tokens: ', `, ,, and ,@
type QuoteLexNode struct {
	BaseNode
}

// NewQuoteLexNode creates a new QuoteLexNode.
func NewQuoteLexNode() *QuoteLexNode {
	return &QuoteLexNode{}
}

// Handle checks if the current character is a quote symbol and returns the appropriate token.
func (q *QuoteLexNode) Handle(r Reader) (Token, bool) {
	if r.EOF() {
		return nil, false
	}

	file, pos, line, linePos := r.Snapshot()
	ch := r.Current()

	switch ch {
	case '\'':
		r.Advance(1)
		return NewToken(file, pos, line, linePos, "'", Quote), true
	case '`':
		r.Advance(1)
		return NewToken(file, pos, line, linePos, "`", Backtick), true
	case ',':
		r.Advance(1)
		// Check for ,@ (unquote-splicing)
		if !r.EOF() && r.Current() == '@' {
			r.Advance(1)
			return NewToken(file, pos, line, linePos, ",@", CommaAt), true
		}
		return NewToken(file, pos, line, linePos, ",", Comma), true
	default:
		return q.PassToNext(r)
	}
}

