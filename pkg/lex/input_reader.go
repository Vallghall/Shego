package lex

// Reader provides methods for reading input and tracking position.
// It encapsulates the input string and current position state.
type Reader interface {
	// Position information
	File() string
	Position() int
	Line() int
	LinePosition() int

	// Reading operations
	EOF() bool
	Current() byte
	Peek(offset int) (byte, bool)
	PeekString(length int) (string, bool)

	// Advance moves the cursor forward by n characters, updating line tracking.
	Advance(n int)

	// SkipWhitespace advances past any whitespace characters.
	// Returns true if any whitespace was skipped.
	SkipWhitespace() bool

	// Fork creates a snapshot of current position for lookahead.
	// The returned values can be used to create tokens.
	Snapshot() (file string, pos int, line int, linePos int)
}

// inputReader is the concrete implementation of Reader.
type inputReader struct {
	input   string
	file    string
	pos     int
	line    int
	linePos int
}

// NewReader creates a new Reader for the given input.
func NewReader(input string, file string) Reader {
	return &inputReader{
		input:   input,
		file:    file,
		pos:     0,
		line:    1,
		linePos: 1,
	}
}

func (r *inputReader) File() string      { return r.file }
func (r *inputReader) Position() int     { return r.pos }
func (r *inputReader) Line() int         { return r.line }
func (r *inputReader) LinePosition() int { return r.linePos }

// EOF returns true if the reader has reached the end of input.
func (r *inputReader) EOF() bool {
	return r.pos >= len(r.input)
}

// Current returns the current character, or 0 if at EOF.
func (r *inputReader) Current() byte {
	if r.EOF() {
		return 0
	}
	return r.input[r.pos]
}

// Peek returns the character at offset positions ahead of current.
// Returns the character and true if valid, or 0 and false if out of bounds.
func (r *inputReader) Peek(offset int) (byte, bool) {
	idx := r.pos + offset
	if idx < 0 || idx >= len(r.input) {
		return 0, false
	}
	return r.input[idx], true
}

// PeekString returns the next n characters starting from current position.
// Returns the string and true if all characters are available.
func (r *inputReader) PeekString(length int) (string, bool) {
	if r.pos+length > len(r.input) {
		return "", false
	}
	return r.input[r.pos : r.pos+length], true
}

// Advance moves the cursor forward by n characters, updating line tracking.
func (r *inputReader) Advance(n int) {
	for i := 0; i < n && r.pos < len(r.input); i++ {
		if r.input[r.pos] == '\n' {
			r.line++
			r.linePos = 1
		} else {
			r.linePos++
		}
		r.pos++
	}
}

// SkipWhitespace advances past any whitespace characters.
func (r *inputReader) SkipWhitespace() bool {
	skipped := false
	for !r.EOF() {
		ch := r.Current()
		if ch == ' ' || ch == '\t' {
			r.Advance(1)
			skipped = true
		} else if ch == '\n' {
			r.Advance(1)
			skipped = true
		} else if ch == '\r' {
			r.Advance(1)
			// Handle \r\n as single newline
			if !r.EOF() && r.Current() == '\n' {
				r.pos++ // Don't call Advance to avoid double line increment
			}
			skipped = true
		} else {
			break
		}
	}
	return skipped
}

// Snapshot returns the current position state for token creation.
func (r *inputReader) Snapshot() (file string, pos int, line int, linePos int) {
	return r.file, r.pos, r.line, r.linePos
}
