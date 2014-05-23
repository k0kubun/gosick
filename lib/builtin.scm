(define cadr
  (lambda (x)
    (car (cdr x))
  )
)
(define cddr
  (lambda (x)
    (cdr (cdr x))
  )
)

(define not
  (lambda (x)
    (eq? x #f)
  )
)
