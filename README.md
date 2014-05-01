# Gosick

Scheme implementation by Go, which is specified by [R5RS](http://www.schemers.org/Documents/Standards/R5RS/r5rs.pdf).  
This is started as [a programming project](https://github.com/k0kubun/gosick/blob/master/project.md) for newcomers of my laboratory.

## Specification

### Implemented syntax and functions

| Type | List | Implemented |
|:-----|:-----|:-----------:|
| Number | number?, +, -, *, /, =, <, <=, >, >= | + |
| List | null?, pair?, list?, symbol?, car, cdr, cons, list, length, memq, last, append, set-car!, set-cdr! |  |
| Boolean | boolean?, not |  |
| String | string?, string-append, symbol->string, string->symbol, string->number, number->string |  |
| Procedure | procedure? |  |
| Comparison | eq?, neq?, equal? |  |
| Syntax | lambda, let, let*, letrec |  |
| Statement | if, cond, and, or, begin, do |  |
| Definition | set!, define, define-macro |  |
| Others | load |  |

## License

Gosick is released under the [MIT License](http://opensource.org/licenses/MIT).