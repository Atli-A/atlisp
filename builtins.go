package main

import (
	"errors"
	"fmt"
)

var (
	builtins = map[string](func(...Var) (Var, RuntimeError)){
		"+": Add,
		"-": Subtract,
		"*": Multiply,
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

func Add(vars ...Var) (Var, RuntimeError) {
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

func Subtract(vars ...Var) (Var, RuntimeError) {
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

func Multiply(vars ...Var) (Var, RuntimeError) {
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
func GenerateBuiltins() map[string]Var {
	res := make(map[string]Var)
	for k, v := range builtins {
		res[k] = Var{
			Data: v,
			Type: VarTypes.BUILTIN,
		}
		fmt.Println("hihih")
	}
	return res
}
