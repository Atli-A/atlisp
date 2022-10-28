# Notes

Now that we have a basic lisp lets go through a few ideas.
## Slices
Replace (cons 1 2) -> (1 . 2) to (slice 1 2 3) -> (1 2 3) slices have constant time access. They are in implementations pointers to an array with length and cap. Similar to golang's implementation but append will actually change the position if needed. so you don't have to weirdly assign it. 

Code can be represented as slices of types??

## Types
integers: (u8, u16, u32, u64, i8, i16, i32, i64)
floats: (f16?, f32, f64)

### Sum Types
examples of sum types 
integers as seen above
floats as seen above

addable: (integers, floats)
any: literally anything i guess

### Function Types
functions signatures are now like this:
have the parameters, their types, and the return type. All of these types can be sum types

### Slice Types
cons is replaced with slices probably
slice variable signatures are like this ??? 

### Strings
????

### Reflection
var.type returns a type 
how does this work for functions?
thats a TODO 

### Typed Macros
TODO?? Macros rearrange AND evaluate code
Some common-lisp-ish based on this idea
```lisp
(defmacro do3 (some-code slice-any n u8) do-stuff)
```

## Generics
### Generic Functions
Generic functions can be achieved in a few ways. I am undecided on this one.
Say we have a functions as follows
func (n addable) returns addable
One cool way that avoids significantly changing the syntax would be to get the type of n and use it. therefore we can return the same type.
perhaps return types can be determined from the parameter types. However if I try to add i64 and u16 what should occur?
maybe:
func (n addable, m n.type) returns n.type
this is bad because we can't immediately tell that n m and the return type are all the same.
this brings us to the obvious alternative
```
func (T addable) (n T, m T) return T
```
How do we make this work with lisp??
A zig style generic in lisp form could work:
```
fn max(comptime T: type, a: T, b: T) T {...}
```
This has potential and introduces a "type" type. Could be elegant?
To support currying we would need to ensure that type params are first. Otherwise you could curry a variable that conflicts with a type. Even if its solvable it overcomplicates it.




## Exports
```lisp
(export (def x 1))  ; exports variable x to other files. 
; This syntax is nice but has little success as expressions are evaluated
; what would this even do with any input?
; instead we could
(def x 1)
(export x) ; this solves our problem i think
```

### Imports
```lisp
(compile "filename.atlisp" import-name) ; must be string so that filenames are properly handled
; should import name be used? or should it instead be inferred form filename? i would argue for import-name but idk
(print import-name::x) ; still deciding between the :: syntax and the . syntax. "." will probably be used for type access so it poses an interesting idea.
; or 
(print import-name.x)
```

## Consts
To make constants which is clearly valuable we use the def syntax but use const instead

