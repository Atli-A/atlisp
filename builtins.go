package main

import (
	"errors"
	"fmt"
)

var (
	builtins = map[string](func(...Var) (Var, RuntimeError)){
		"+": add,
		"-": subtract,
		"*": multiply,
		"/": divide,
		"%": mod,
		"cons": cons,
		"car":  car,
		"cdr":  cdr,
		"print": print_,
		"eq": eq,
	}
)

func TypeMatches(v Var, vt VarType) (bool, RuntimeError) {
	if v.Type == vt {
		return true, RuntimeError{}
	}
	return false, RuntimeError{
		errors.New("Incorrect type expects TODO"), 0, 0,
	}
}

func add(vars ...Var) (Var, RuntimeError) {
	sum := int64(0)
	for _, v := range vars {
		if matches, err := TypeMatches(v, VarTypes.NUM); !matches {
			return Var{}, err
		}
		sum += v.Data.(int64)
	}
	return Var{
		Data: sum,
		Type: VarTypes.NUM,
	}, RuntimeError{}
}

func subtract(vars ...Var) (Var, RuntimeError) {
	if len(vars) == 0 {
		return Var{}, RuntimeError{errors.New("- requires one or more arguments"), 0, 0}
	} else if len(vars) == 1 {
		return Var{
			Type: VarTypes.NUM,
			Data: -(vars[0].Data.(int64)),
		}, RuntimeError{}
	}
	first := vars[0].Data.(int64)
	for _, v := range vars[1:] {
		if matches, err := TypeMatches(v, VarTypes.NUM); !matches {
			return Var{}, err
		}
		first -= v.Data.(int64)
	}
	return Var{
		Data: first,
		Type: VarTypes.NUM,
	}, RuntimeError{}
}

func multiply(vars ...Var) (Var, RuntimeError) {
	mul := int64(1)
	for _, v := range vars {
		if matches, err := TypeMatches(v, VarTypes.NUM); !matches {
			return Var{}, err
		}
		mul *= v.Data.(int64)
	}
	return Var{
		Data: mul,
		Type: VarTypes.NUM,
	}, RuntimeError{}
}

func divide(vars ...Var) (Var, RuntimeError) {
	if len(vars) == 0 {
		return Var{}, RuntimeError{
			errors.New("divide requires at least on value"), 0, 0,
		}
	} else if len(vars) == 1 {
		if matches, err := TypeMatches(vars[0], VarTypes.NUM); !matches {
			return Var{}, err
		}
		return Var{
			Data: int64(1)/vars[0].Data.(int64),
			Type: VarTypes.NUM,
		}, RuntimeError{}
	} else {
		for _, v := range vars {
			if matches, err := TypeMatches(v, VarTypes.NUM); !matches {
				return Var{}, err
			}
		}
		start := vars[0].Data.(int64)
		for _, v := range vars[1:] {
			start /= v.Data.(int64)
		}
		return Var{
			Data: start,
			Type: VarTypes.NUM,
		}, RuntimeError{}
	}
	
}

func mod(vars ...Var) (Var, RuntimeError) {
	if len(vars) != 2 {
		return Var{}, RuntimeError{
			errors.New("% requires exactly 2 arguments"), 0, 0,
		}
	}
	if matches, err := TypeMatches(vars[0], VarTypes.NUM); !matches {
		return Var{}, err
	}
	if matches, err := TypeMatches(vars[1], VarTypes.NUM); !matches {
		return Var{}, err
	}
	return Var{
		Data: vars[0].Data.(int64) % vars[1].Data.(int64),
		Type: VarTypes.NUM,
	}, RuntimeError{}

}

func cons(vars ...Var) (Var, RuntimeError) {
	if len(vars) != 2 {
		return Var{}, RuntimeError{
			errors.New("cons expects exactly 2 arguments"), 0, 0,
		}
	}
	return Var{
		Data: Cons{
			vars[0],
			vars[1],
		},
		Type: VarTypes.CONS,
	}, RuntimeError{}
}

func car(vars ...Var) (Var, RuntimeError) {
	if len(vars) != 1 {
		return Var{}, RuntimeError{
			errors.New("car expects exactly 1 argument"), 0, 0,
		}
	}
	if vars[0].Type != VarTypes.CONS {
		return Var{}, RuntimeError{
			errors.New("car expects a cons or list"), 0, 0,
		}
	}
	return vars[0].Data.(Cons).First, RuntimeError{}
}

func cdr(vars ...Var) (Var, RuntimeError) {
	if len(vars) != 1 {
		return Var{}, RuntimeError{
			errors.New("cdr expects exactly 1 argument"), 0, 0,
		}
	}
	if vars[0].Type != VarTypes.CONS {
		return Var{}, RuntimeError{
			errors.New("cdr expects a cons or list"), 0, 0,
		}
	}
	return vars[0].Data.(Cons).Rest, RuntimeError{}
}

func print_(vars ...Var) (Var, RuntimeError) {
	if len(vars) != 1 {
		return Var{}, RuntimeError{
			errors.New("print expects exactly 1 argument"), 0, 0,
		}
	}
	fmt.Println(vars[0].Data)
	return vars[0], RuntimeError{}
}

func eq(vars ...Var) (Var, RuntimeError) {
	// TODO check len >= 1
	first := vars[0]
	for i := 1; i < len(vars); i++ {
		if first != vars[i] {
			return Var{
				Data: nil,
				Type: VarTypes.NIL,
			}, RuntimeError{}
		}
	}
	return Var{
		Data: "true",
		Type: VarTypes.SYMBOL,
	}, RuntimeError{}
}

func GenerateBuiltins() map[string]Var {
	res := make(map[string]Var)
	for k, v := range builtins {
		res[k] = Var{
			Data: v,
			Type: VarTypes.BUILTIN,
		}
	}
	return res
}

