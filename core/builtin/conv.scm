;;; Schego Core Library - Conversion Operations
;;; This file contains Scheme definition stubs for conversion builtins.
;;; Actual implementations are in Go: pkg/core/builtin/conv.go

;;; Number to string
;;; (number->string n) -> string
;;; Converts a number to its string representation.
(define number->string
  "Implemented in Go: converts number to string representation.")

;;; String to number
;;; (string->number s) -> number or #f
;;; Parses a string as a number. Returns #f if unparseable.
(define string->number
  "Implemented in Go: parses string as number, returns #f if invalid.")

;;; Symbol to string
;;; (symbol->string sym) -> string
;;; Returns the symbol's name as a string.
(define symbol->string
  "Implemented in Go: returns symbol's name as string.")

;;; String to list
;;; (string->list s) -> list
;;; Converts a string to a list of single-character strings.
(define string->list
  "Implemented in Go: converts string to list of characters.")

;;; List to string
;;; (list->string lst) -> string
;;; Converts a list of strings to a concatenated string.
(define list->string
  "Implemented in Go: concatenates list of strings.")

;;; Integer to character
;;; (integer->char n) -> string
;;; Converts a Unicode code point to a single-character string.
(define integer->char
  "Implemented in Go: converts Unicode code point to character.")

;;; Character to integer
;;; (char->integer c) -> number
;;; Converts a single-character string to its Unicode code point.
(define char->integer
  "Implemented in Go: converts character to Unicode code point.")

