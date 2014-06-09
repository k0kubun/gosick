(define concurrency 1)

(define master
  (actor
    (("sum-upto" last)
      (define result 0)
      (define count 0)
      (define per-proc (/ last concurrency))
      (do ((i 0))
        ((= i concurrency))
        (sum-range (+ (* i per-proc) 1) (* (+ i 1) per-proc) self)
        (set! i (+ i 1))))

    (("add-result" num)
      (set! result (+ result num))
      (set! count (+ count 1))
      (if (= count concurrency)
        (begin (print result) (exit))))))

(define (sum-range start end master)
  (let ((child (generate-child)))
    (child start)
    (child ! "sum-range" start end master)))

(define (generate-child)
  (actor
    (("sum-range" range-start range-end parent)
      (do ((sum 0) (i range-start))
        ((> i range-end)
          (parent ! "add-result" sum))
        (set! sum (+ sum i))
        (set! i (+ i 1))))))

(master start)
(master ! "sum-upto" 1200)

(do () (#f))
