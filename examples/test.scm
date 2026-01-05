; Test file for Schego interpreter

; Basic arithmetic
(define x (+ 1 2 3))
(define y (* 4 5))

; Function definition
(define (square n)
  (* n n))

; Result
(+ x y (square 3))

