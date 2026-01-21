;;; Schego Core Library - Eval Operation
;;; This file contains Scheme definition stub for eval builtin.
;;; Actual implementation is in Go: pkg/core/builtin/eval.go

;;; Eval
;;; (eval obj) -> value
;;; Evaluates any Scheme object. Converts the object to AST and evaluates it.
;;; The object can be a number, string, symbol, or list (s-expression).
(define eval
  "Implemented in Go: evaluates any Scheme object by converting to AST and evaluating.")
