# Atlisp
atlisp is a minimal lisp that I created in my free time

The name is portmanteua of atli and lisp

## Usage
To run a REPL of atlisp use ./atlisp.
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
`progn` evaluate a series of expressions 

`fn` create a function

`def` define a variable

`quote` return a symbol or variable unmodified or return an expression as a list of symbols

### Future work
- Add `eval`
- Add special characters to simplify common operations such as ' , `

