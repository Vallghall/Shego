package lex

// Node represents a handler in the Chain of Responsibility for lexing tokens.
type Node interface {
	// Handle attempts to parse a token at the current position in input.
	// Returns:
	//   - token: the parsed token if successful, nil otherwise
	//   - consumed: number of characters consumed from input
	//   - handled: true if this node handled the input, false to pass to next
	Handle(input string, pos int, line int, linePos int, file string) (token Token, consumed int, handled bool)

	// SetNext sets the next node in the chain.
	SetNext(Node)

	// Next returns the next node in the chain.
	Next() Node
}

// BaseNode provides common chain traversal logic for Node implementations.
type BaseNode struct {
	next Node
}

// SetNext sets the next node in the chain.
func (b *BaseNode) SetNext(n Node) {
	b.next = n
}

// Next returns the next node in the chain.
func (b *BaseNode) Next() Node {
	return b.next
}

// PassToNext delegates handling to the next node in the chain.
// Returns nil, 0, false if there is no next node.
func (b *BaseNode) PassToNext(input string, pos int, line int, linePos int, file string) (Token, int, bool) {
	if b.next != nil {
		return b.next.Handle(input, pos, line, linePos, file)
	}
	return nil, 0, false
}

// BuildChain constructs a chain from the given nodes in order.
// Returns the head of the chain.
func BuildChain(nodes ...Node) Node {
	if len(nodes) == 0 {
		return nil
	}
	for i := 0; i < len(nodes)-1; i++ {
		nodes[i].SetNext(nodes[i+1])
	}
	return nodes[0]
}
