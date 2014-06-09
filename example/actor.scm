(define master
  (actor
    (("hello" name)
      (print name))
    (("test" name)
      (print name)
      (print name))))

(master start)
