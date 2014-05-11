# Gosick

Scheme implementation by Go, which is specified by [R5RS](http://www.schemers.org/Documents/Standards/R5RS/r5rs.pdf).  
This is started as [a programming project](https://github.com/k0kubun/gosick/blob/master/project.md) for newcomers of my laboratory.

## Installation

```bash
$ go get github.com/k0kubun/gosick
```

## Usage

```bash
# Invoke interactive shell
$ gosick

# Excecute scheme source code
$ gosick -f source.scm

# One liner
$ gosick -e "(+ 1 2)"

# Dump AST of input source code
$ gosick -a

# Show help
$ gosick -h
```

## Specification

### Implemented syntax and functions

| Type | To be done | Implemented |
|:-----|:-----|:-----------:|
| Number | | number?, +, -, *, /, =, <, <=, >, >= |
| List | | cons, car, cdr, list, length, memq, last, append, set-car!, set-cdr! |
| Boolean | | not |
| String | | string-append, symbol->string, string->symbol, string->number, number->string |
| Type | | null?, boolean?, procedure?, pair?, list?, symbol?, string? |
| Comparison | | eq?, neq?, equal? |
| Syntax | | lambda, let, let*, letrec |
| Statement | cond, and, or, begin, do | if |
| Definition | define-macro | define, set! |
| Others | | load |

## License

Gosick is released under the [MIT License](http://opensource.org/licenses/MIT).
