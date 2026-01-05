package lex

// ParenNode handles parentheses tokens: ( and )
type ParenNode struct {
	BaseNode
}

// NewParenNode creates a new ParenNode.
func NewParenNode() *ParenNode {
	return &ParenNode{}
}

// Handle checks if the current character is a parenthesis and returns the appropriate token.
func (p *ParenNode) Handle(r Reader) (Token, bool) {
	if r.EOF() {
		return nil, false
	}

	file, pos, line, linePos := r.Snapshot()
	ch := r.Current()

	switch ch {
	case '(':
		r.Advance(1)
		return NewToken(file, pos, line, linePos, "(", ParenOpen), true
	case ')':
		r.Advance(1)
		return NewToken(file, pos, line, linePos, ")", ParenClose), true
	default:
		return p.PassToNext(r)
	}
}
