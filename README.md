# Gosick

Scheme interpreter implemented by Golang.  
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
$ gosick source.scm

# One liner
$ gosick -e "(+ 1 2)"

# Dump AST of input source code
$ gosick -a

# Show help
$ gosick -h
```

## Implemented syntax and functions
- +, -, *, /, =, <, <=, >, >=
- cons, car, cdr, list, length, last, append, set-car!, set-cdr!
- if, cond, and, or, not, begin, do
- memq, eq?, neq?, equal?
- null?, number?, boolean?, procedure?, pair?, list?, symbol?, string?
- string-append, symbol->string, string->symbol, string->number, number->string
- let, let*, letrec, lambda, define, set!, quote
- write, print, load

## License

Gosick is released under the [MIT License](http://opensource.org/licenses/MIT).
