;;; Schego Core Library - Load Operation
;;; This file contains Scheme definition stub for load builtin.
;;; Actual implementation is in Go: pkg/core/builtin/load.go

;;; Load
;;; (load path) -> value
;;; Loads and evaluates a Scheme file. The path is searched:
;;; 1. In the current working directory
;;; 2. In the core/builtin/ library directory
;;; Returns the last result of evaluating the file, or void.
(define load
  "Implemented in Go: loads and evaluates Scheme files from path.")
