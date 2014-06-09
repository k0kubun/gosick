(define master
  (actor
    (("sum-upto" last)
      (print "hello")
      ((generate-child) ! "sum-range" 0 10 self)
    )
    (("add-result" num)
      (exit)
    )
  )
)

(define (generate-child)
  (actor
    (("sum-range" range-start range-end parent)
      (do ((sum 0) (i range-start))
        ((> i range-end)
          (parent ! "add-result" sum)
        )
        (set! sum (+ sum i))
        (set! i (+ i 1))
      )
    )
  )
)

(master start)
(master ! "sum-upto" 1000)

(do () (#f))
