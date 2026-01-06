package builtin

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/Vallghall/schego/pkg/mem"
)

// ConvDefinitions returns all conversion operation definitions.
func ConvDefinitions() []mem.Definition {
	return []mem.Definition{
		{Name: "number->string", Value: mem.NewPrimitive("number->string", 1, numberToString)},
		{Name: "string->number", Value: mem.NewPrimitive("string->number", 1, stringToNumber)},
		{Name: "symbol->string", Value: mem.NewPrimitive("symbol->string", 1, symbolToString)},
		{Name: "string->list", Value: mem.NewPrimitive("string->list", 1, stringToList)},
		{Name: "list->string", Value: mem.NewPrimitive("list->string", 1, listToString)},
		{Name: "integer->char", Value: mem.NewPrimitive("integer->char", 1, integerToChar)},
		{Name: "char->integer", Value: mem.NewPrimitive("char->integer", 1, charToInteger)},
	}
}

// numberToString implements (number->string n) - converts a number to its string representation.
func numberToString(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of number->string")
	}

	// Format as integer if possible
	if n.IsInteger() {
		return mem.NewString(strconv.FormatInt(n.Int64(), 10)), nil
	}
	return mem.NewString(strconv.FormatFloat(n.Value(), 'g', -1, 64)), nil
}

// stringToNumber implements (string->number s) - parses a string as a number.
// Returns #f if the string cannot be parsed as a number.
func stringToNumber(args []mem.Object) (mem.Object, error) {
	s, err := mem.AsString(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeString, args[0].Type(),
			"argument 1 of string->number")
	}

	// Try to parse as number
	f, err := strconv.ParseFloat(s.Value(), 64)
	if err != nil {
		// Return #f for unparseable strings (Scheme convention)
		return mem.False, nil
	}
	return mem.NewNumber(f), nil
}

// symbolToString implements (symbol->string sym) - returns the symbol's name as a string.
func symbolToString(args []mem.Object) (mem.Object, error) {
	sym, err := mem.AsSymbol(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeSymbol, args[0].Type(),
			"argument 1 of symbol->string")
	}

	return mem.NewString(sym.Name()), nil
}

// stringToList implements (string->list s) - converts a string to a list of single-character strings.
func stringToList(args []mem.Object) (mem.Object, error) {
	s, err := mem.AsString(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeString, args[0].Type(),
			"argument 1 of string->list")
	}

	str := s.Value()
	chars := make([]mem.Object, 0, utf8.RuneCountInString(str))
	for _, r := range str {
		chars = append(chars, mem.NewString(string(r)))
	}

	return mem.SliceToList(chars), nil
}

// listToString implements (list->string lst) - converts a list of single-character strings to a string.
func listToString(args []mem.Object) (mem.Object, error) {
	elements, err := mem.ListToSlice(args[0])
	if err != nil {
		return nil, fmt.Errorf("list->string: %w", err)
	}

	result := ""
	for i, elem := range elements {
		s, err := mem.AsString(elem)
		if err != nil {
			return nil, mem.NewTypeError(mem.TypeString, elem.Type(),
				fmt.Sprintf("element %d of list in list->string", i+1))
		}
		result += s.Value()
	}

	return mem.NewString(result), nil
}

// integerToChar implements (integer->char n) - converts an integer to a single-character string.
// The integer is treated as a Unicode code point.
func integerToChar(args []mem.Object) (mem.Object, error) {
	n, err := mem.AsNumber(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeNumber, args[0].Type(),
			"argument 1 of integer->char")
	}

	if !n.IsInteger() {
		return nil, fmt.Errorf("integer->char: expected integer, got %s", n.String())
	}

	codePoint := n.Int64()
	if codePoint < 0 || codePoint > 0x10FFFF {
		return nil, fmt.Errorf("integer->char: code point out of range: %d", codePoint)
	}

	return mem.NewString(string(rune(codePoint))), nil
}

// charToInteger implements (char->integer c) - converts a single-character string to its Unicode code point.
func charToInteger(args []mem.Object) (mem.Object, error) {
	s, err := mem.AsString(args[0])
	if err != nil {
		return nil, mem.NewTypeError(mem.TypeString, args[0].Type(),
			"argument 1 of char->integer")
	}

	str := s.Value()
	runeCount := utf8.RuneCountInString(str)
	if runeCount != 1 {
		return nil, fmt.Errorf("char->integer: expected single character, got string of length %d", runeCount)
	}

	r, _ := utf8.DecodeRuneInString(str)
	return mem.NewNumber(float64(r)), nil
}
