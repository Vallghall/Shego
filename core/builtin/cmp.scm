;;; Schego Core Library - Comparison Operations
;;; This file contains Scheme definition stubs for comparison builtins.
;;; Actual implementations are in Go: pkg/core/builtin/cmp.go

;;; Numeric equality
;;; (= n1 n2 ...) -> boolean
;;; Returns #t if all numbers are equal.
(define =
  "Implemented in Go: returns #t if all numbers are equal.")

;;; Less than
;;; (< n1 n2 ...) -> boolean
;;; Returns #t if arguments are strictly increasing.
(define <
  "Implemented in Go: returns #t if strictly increasing.")

;;; Greater than
;;; (> n1 n2 ...) -> boolean
;;; Returns #t if arguments are strictly decreasing.
(define >
  "Implemented in Go: returns #t if strictly decreasing.")

;;; Less than or equal
;;; (<= n1 n2 ...) -> boolean
;;; Returns #t if arguments are non-decreasing.
(define <=
  "Implemented in Go: returns #t if non-decreasing.")

;;; Greater than or equal
;;; (>= n1 n2 ...) -> boolean
;;; Returns #t if arguments are non-increasing.
(define >=
  "Implemented in Go: returns #t if non-increasing.")

;;; Zero predicate
;;; (zero? n) -> boolean
;;; Returns #t if n is zero.
(define zero?
  "Implemented in Go: returns #t if argument is zero.")

;;; Positive predicate
;;; (positive? n) -> boolean
;;; Returns #t if n is positive.
(define positive?
  "Implemented in Go: returns #t if argument is positive.")

;;; Negative predicate
;;; (negative? n) -> boolean
;;; Returns #t if n is negative.
(define negative?
  "Implemented in Go: returns #t if argument is negative.")

;;; Identity comparison
;;; (eq? a b) -> boolean
;;; Returns #t if a and b are the same object (identity comparison).
;;; For symbols, compares by ID. For primitives, compares by value.
(define eq?
  "Implemented in Go: identity comparison (pointer/symbol ID).")

;;; Structural equality
;;; (equal? a b) -> boolean
;;; Returns #t if a and b are structurally equal.
(define equal?
  "Implemented in Go: structural equality comparison.")

;;; Boolean negation
;;; (not x) -> boolean
;;; Returns #t if x is #f, otherwise returns #f.
(define not
  "Implemented in Go: boolean negation.")

