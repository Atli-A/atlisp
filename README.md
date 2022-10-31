# Atlisp
Atlisp is a minimal lisp that I created in my free time

The name is a portmanteau of atli and lisp.


## Build
Requires go 1.18 or above.

After cloning the repo you can run:
```bash
go build
```

## Usage
Prequisites: Build

To run a REPL of atlisp use ./atlisp

To run a file with atlisp use ./atlisp <filename>

## Tutorial 
To call a function use the following syntax `(function_name parameter1 parameter2)`. Operations such as addition are functions ex:
```lisp
 > (+ 1 2)
3
 > (* 10 9)
90
```

### Variables
Variables and functions share a namespace. To define a variable use the following syntax
```lisp
 > (def myvar "a string")
"a string"
 > myvar
"a string"
```

### Functions
Functions are defined with a similar syntax
```lisp
 > (def factorial
..   (fn (n)
..     (if (eq n 1)
..         1
..         (* n (factorial (- n 1))))))
```
Here the `fn` keyword is used to create a function. To call this function we would say:
```
 > (factorial 5)
 120
```

### Keywords / Special Forms
`eq`: checks for equality

`progn`: evaluate a series of expressions and returns the output of the last one.

`fn`: creates a function as seen earlier

`def`: define a variable as seen earlier

`quote`: return a symbol or variable unmodified or return an expression as a list of symbols
```lisp
 > (quote test)
test
```
`if`: takes an expression which it evaluates. If the expression is truthy it evaluates the next expression otherwise it doesn't evaluate that one and instead evaluates the other.
For example: 
```lisp
(if (eq 1 1) "this" "that")
```
evaluates to "this" and "that" is never evaluated.


### Future work
- Add `eval`
- Add special characters to simplify common operations such as ' , `

