;;; Schego Core Library - List Operations
;;; This file contains Scheme definition stubs for list builtins.
;;; Actual implementations are in Go: pkg/core/builtin/list.go

;;; =============================================================================
;;; Basic Pair Operations
;;; =============================================================================

;;; Cons
;;; (cons a b) -> pair
;;; Creates a new pair with a as car and b as cdr.
(define cons
  "Implemented in Go: creates a new pair (a . b).")

;;; Car
;;; (car pair) -> object
;;; Returns the first element of a pair.
(define car
  "Implemented in Go: returns first element of pair.")

;;; Cdr
;;; (cdr pair) -> object
;;; Returns the second element (rest) of a pair.
(define cdr
  "Implemented in Go: returns second element of pair.")

;;; =============================================================================
;;; Two-level Compositions
;;; =============================================================================

;;; Caar
;;; (caar x) -> object
;;; Equivalent to (car (car x)).
(define caar
  "Implemented in Go: (car (car x)).")

;;; Cadr
;;; (cadr x) -> object
;;; Returns the second element. Equivalent to (car (cdr x)).
(define cadr
  "Implemented in Go: (car (cdr x)) - second element.")

;;; Cdar
;;; (cdar x) -> object
;;; Equivalent to (cdr (car x)).
(define cdar
  "Implemented in Go: (cdr (car x)).")

;;; Cddr
;;; (cddr x) -> object
;;; Returns the rest after second. Equivalent to (cdr (cdr x)).
(define cddr
  "Implemented in Go: (cdr (cdr x)) - rest after second.")

;;; =============================================================================
;;; Three-level Compositions
;;; =============================================================================

;;; (caaar x) = (car (car (car x)))
(define caaar
  "Implemented in Go: (car (car (car x))).")

;;; (caadr x) = (car (car (cdr x)))
(define caadr
  "Implemented in Go: (car (car (cdr x))).")

;;; (cadar x) = (car (cdr (car x)))
(define cadar
  "Implemented in Go: (car (cdr (car x))).")

;;; (caddr x) = (car (cdr (cdr x))) - third element
(define caddr
  "Implemented in Go: (car (cdr (cdr x))) - third element.")

;;; (cdaar x) = (cdr (car (car x)))
(define cdaar
  "Implemented in Go: (cdr (car (car x))).")

;;; (cdadr x) = (cdr (car (cdr x)))
(define cdadr
  "Implemented in Go: (cdr (car (cdr x))).")

;;; (cddar x) = (cdr (cdr (car x)))
(define cddar
  "Implemented in Go: (cdr (cdr (car x))).")

;;; (cdddr x) = (cdr (cdr (cdr x))) - rest after third
(define cdddr
  "Implemented in Go: (cdr (cdr (cdr x))) - rest after third.")

;;; =============================================================================
;;; List Utilities
;;; =============================================================================

;;; List
;;; (list a b c ...) -> list
;;; Creates a proper list from the arguments.
(define list
  "Implemented in Go: creates a proper list from arguments.")

;;; Length
;;; (length lst) -> number
;;; Returns the number of elements in a proper list.
(define length
  "Implemented in Go: returns length of list.")

;;; Append
;;; (append lst1 lst2 ...) -> list
;;; Concatenates lists. Last argument may be any object.
(define append
  "Implemented in Go: concatenates lists.")

;;; Reverse
;;; (reverse lst) -> list
;;; Returns a reversed copy of the list.
(define reverse
  "Implemented in Go: returns reversed copy of list.")

;;; =============================================================================
;;; Predicates
;;; =============================================================================

;;; Null?
;;; (null? obj) -> boolean
;;; Returns #t if obj is the empty list '().
(define null?
  "Implemented in Go: returns #t if empty list.")

;;; Pair?
;;; (pair? obj) -> boolean
;;; Returns #t if obj is a pair.
(define pair?
  "Implemented in Go: returns #t if pair.")

;;; List?
;;; (list? obj) -> boolean
;;; Returns #t if obj is a proper list (nil-terminated).
(define list?
  "Implemented in Go: returns #t if proper list.")

