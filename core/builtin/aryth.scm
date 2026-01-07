;;; Schego Core Library - Arithmetic Operations
;;; This file contains Scheme definition stubs for arithmetic builtins.
;;; Actual implementations are in Go: pkg/core/builtin/aryth.go

;;; Addition
;;; (+ n1 n2 ...) -> number
;;; Returns the sum of all arguments. With no arguments, returns 0.
(define +
  "Implemented in Go: returns sum of all arguments, 0 if none.")

;;; Subtraction
;;; (- n) -> number (negation)
;;; (- n1 n2 ...) -> number (subtraction)
;;; With one argument: returns negation.
;;; With multiple arguments: returns n1 - n2 - n3 - ...
(define -
  "Implemented in Go: negation with 1 arg, subtraction with 2+.")

;;; Multiplication
;;; (* n1 n2 ...) -> number
;;; Returns the product of all arguments. With no arguments, returns 1.
(define *
  "Implemented in Go: returns product of all arguments, 1 if none.")

;;; Division
;;; (/ n) -> number (reciprocal)
;;; (/ n1 n2 ...) -> number (division)
;;; With one argument: returns reciprocal.
;;; With multiple arguments: returns n1 / n2 / n3 / ...
(define /
  "Implemented in Go: reciprocal with 1 arg, division with 2+.")

;;; Modulo
;;; (modulo n1 n2) -> number
;;; Returns the remainder of integer division.
;;; Result has the same sign as the divisor (Scheme semantics).
(define modulo
  "Implemented in Go: integer remainder with Scheme sign semantics.")

;;; Absolute value
;;; (abs n) -> number
;;; Returns the absolute value of n.
(define abs
  "Implemented in Go: returns absolute value.")

;;; Minimum
;;; (min n1 n2 ...) -> number
;;; Returns the minimum of all arguments.
(define min
  "Implemented in Go: returns minimum value.")

;;; Maximum
;;; (max n1 n2 ...) -> number
;;; Returns the maximum of all arguments.
(define max
  "Implemented in Go: returns maximum value.")

