;;; Schego Core Library - Apply Operation
;;; This file contains Scheme definition stub for apply builtin.
;;; Actual implementation is in Go: pkg/core/builtin/apply.go

;;; Apply
;;; (apply proc arg1 ... argN lst) -> value
;;; Applies procedure proc to arguments. The last argument must be a list.
;;; All variadic arguments (arg1 ... argN) are prepended to the list,
;;; and the combined arguments are passed to the procedure.
(define apply
  "Implemented in Go: applies procedure with variadic args + list.")
