;;; Schego Core Library - Mathematical Operations
;;; This file contains Scheme definition stubs for math builtins.
;;; Actual implementations are in Go: pkg/core/math/math.go

;;; ============================================================================
;;; Constants
;;; ============================================================================

;;; Pi (π)
;;; The ratio of a circle's circumference to its diameter.
(define pi
  "Implemented in Go: 3.141592653589793...")

;;; Euler's number (e)
;;; The base of natural logarithms.
(define e
  "Implemented in Go: 2.718281828459045...")

;;; ============================================================================
;;; Rounding Operations
;;; ============================================================================

;;; Floor
;;; (floor n) -> number
;;; Returns the largest integer less than or equal to n.
(define floor
  "Implemented in Go: returns largest integer <= n.")

;;; Ceiling
;;; (ceiling n) -> number
;;; Returns the smallest integer greater than or equal to n.
(define ceiling
  "Implemented in Go: returns smallest integer >= n.")

;;; Round
;;; (round n) -> number
;;; Returns the nearest integer (banker's rounding for ties).
(define round
  "Implemented in Go: returns nearest integer (banker's rounding).")

;;; Truncate
;;; (truncate n) -> number
;;; Returns the integer part, truncating toward zero.
(define truncate
  "Implemented in Go: returns integer part (toward zero).")

;;; ============================================================================
;;; Exponential and Logarithmic
;;; ============================================================================

;;; Square root
;;; (sqrt n) -> number
;;; Returns the square root of n. Error if n is negative.
(define sqrt
  "Implemented in Go: returns square root (error if negative).")

;;; Exponentiation
;;; (expt base exp) -> number
;;; Returns base raised to the power of exp.
(define expt
  "Implemented in Go: returns base^exp.")

;;; Exponential
;;; (exp n) -> number
;;; Returns e raised to the power of n.
(define exp
  "Implemented in Go: returns e^n.")

;;; Logarithm
;;; (log n) -> number (natural log)
;;; (log n base) -> number (log with specified base)
;;; Returns logarithm. With one argument, returns natural log.
;;; With two arguments, returns log base of n.
(define log
  "Implemented in Go: natural log or log with specified base.")

;;; ============================================================================
;;; Trigonometric Functions
;;; ============================================================================

;;; Sine
;;; (sin n) -> number
;;; Returns the sine of n (argument in radians).
(define sin
  "Implemented in Go: returns sine (radians).")

;;; Cosine
;;; (cos n) -> number
;;; Returns the cosine of n (argument in radians).
(define cos
  "Implemented in Go: returns cosine (radians).")

;;; Tangent
;;; (tan n) -> number
;;; Returns the tangent of n (argument in radians).
(define tan
  "Implemented in Go: returns tangent (radians).")

;;; Arcsine
;;; (asin n) -> number
;;; Returns the arcsine of n (result in radians).
;;; Argument must be in range [-1, 1].
(define asin
  "Implemented in Go: returns arcsine (radians), arg in [-1,1].")

;;; Arccosine
;;; (acos n) -> number
;;; Returns the arccosine of n (result in radians).
;;; Argument must be in range [-1, 1].
(define acos
  "Implemented in Go: returns arccosine (radians), arg in [-1,1].")

;;; Arctangent
;;; (atan n) -> number
;;; (atan y x) -> number
;;; With one argument: returns arctangent of n.
;;; With two arguments: returns atan2(y, x), using signs for quadrant.
(define atan
  "Implemented in Go: arctangent or atan2(y, x).")

;;; ============================================================================
;;; Predicates
;;; ============================================================================

;;; NaN predicate
;;; (nan? n) -> boolean
;;; Returns #t if n is NaN (Not a Number).
(define nan?
  "Implemented in Go: returns #t if NaN.")

;;; Infinite predicate
;;; (infinite? n) -> boolean
;;; Returns #t if n is infinite (+∞ or -∞).
(define infinite?
  "Implemented in Go: returns #t if infinite.")

;;; Finite predicate
;;; (finite? n) -> boolean
;;; Returns #t if n is finite (not NaN or infinite).
(define finite?
  "Implemented in Go: returns #t if finite.")

;;; Integer predicate
;;; (integer? n) -> boolean
;;; Returns #t if n is an integer (has no fractional part).
(define integer?
  "Implemented in Go: returns #t if integer.")

