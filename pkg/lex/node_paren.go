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
func (p *ParenNode) Handle(input string, pos int, line int, linePos int, file string) (Token, int, bool) {
	if pos >= len(input) {
		return nil, 0, false
	}

	ch := input[pos]
	switch ch {
	case '(':
		return NewToken(file, pos, line, linePos, "(", ParenOpen), 1, true
	case ')':
		return NewToken(file, pos, line, linePos, ")", ParenClose), 1, true
	default:
		return p.PassToNext(input, pos, line, linePos, file)
	}
}
