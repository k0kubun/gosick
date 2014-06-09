(define master
  (actor
    (("hello")
      (print "hello")
    )
  )
)

(master start)
(master ! "hello")

(do () (#f))
