;;; Schego Core Library - Functional Operations
;;; This file contains Scheme definitions for map, filter, and fold.
;;; These are implemented in Scheme itself, not in Go.

;;; Map
;;; (map proc lst) -> list
;;; Applies proc to each element of lst and returns a list of results.
(define (map proc lst)
  (if (null? lst)
      '()
      (cons (proc (car lst))
            (map proc (cdr lst)))))

;;; Filter
;;; (filter pred lst) -> list
;;; Returns a list containing only elements of lst for which pred returns #t.
(define (filter pred lst)
  (if (null? lst)
      '()
      (if (pred (car lst))
          (cons (car lst) (filter pred (cdr lst)))
          (filter pred (cdr lst)))))

;;; Fold
;;; (fold proc init lst) -> value
;;; Folds lst using proc, starting with init.
;;; proc is called as (proc accumulator element) for each element.
(define (fold proc init lst)
  (if (null? lst)
      init
      (fold proc (proc init (car lst)) (cdr lst))))
