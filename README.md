# Gosick

Scheme implementation by Go, which is specified by [R5RS](http://www.schemers.org/Documents/Standards/R5RS/r5rs.pdf).  
This is started as [a programming project](https://github.com/k0kubun/gosick/blob/master/project.md) for newcomers of my laboratory.

## Specification

### Implemented syntax and functions

| Type | Name | Implemented |
|:-----|:-----|:-----------:|
| Number | number?, +, -, *, /, =, <, <=, >, >= | x |
| List | null?, pair?, list?, symbol?, car, cdr, cons, list, length, memq, last, append, set-car!, set-cdr! | x |
| Boolean | boolean?, not | x |
| String | string?, string-append, symbol->string, string->symbol, string->number, number->string | x |
| Procedure | procedure? | x |
| Comparison | eq?, neq?, equal? | x |
| Syntax | lambda, let, let*, letrec | x |
| Statement | if, cond, and, or, begin, do | x |
| Definition | set!, define, define-macro | x |
| Others | load | x |

### Optimization

| Name | Implemented |
|:-----|:-----------:|
| Tail Call Optimization | x |

## License

Gosick is released under the [MIT License](http://opensource.org/licenses/MIT).
