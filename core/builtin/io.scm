;;; Schego Core Library - I/O Operations
;;; This file contains Scheme definition stubs for I/O builtins.
;;; Actual implementations are in Go: pkg/core/builtin/io.go

;;; Display
;;; (display obj) -> void
;;; Prints the string representation of obj to stdout.
;;; For strings, prints without quotes.
;;; For other types, prints the standard representation.
(define display
  "Implemented in Go: prints object to stdout, strings without quotes.")

;;; Write
;;; (write obj) -> void
;;; Prints the Scheme-readable representation of obj to stdout.
;;; For strings, preserves the surrounding quotes.
;;; Useful for outputting data that can be read back by the Scheme reader.
(define write
  "Implemented in Go: prints object to stdout, strings with quotes.")

;;; Newline
;;; (newline) -> void
;;; Prints a newline character to stdout.
(define newline
  "Implemented in Go: prints newline to stdout.")

