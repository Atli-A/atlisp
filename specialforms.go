package main

import (
	"errors"
	"fmt"
)

func sliceToConsQuote(exprs []*Expression) (Var, RuntimeError) {
	var res Cons // TODO start with nil somehow

	for i := len(exprs)-1; i >= 0; i-- {
		v, Err := Quote(*exprs[i])
		if Err.Exists() {
			return Var{}, RuntimeError{
				errors.New(fmt.Sprintf("Quote Failed!!!! on %v", exprs[i])), 0, 0,
			}
		}
		var prev Var
		if i != len(exprs) - 1 {
			prev =	Var{
				Data: res,
				Type: VarTypes.CONS,
			}
		} else {
			prev = Var{
				Data: nil,
				Type: VarTypes.NIL,
			}
		}
		res = Cons{
			v,
			prev,
		}
	}

	return Var{
		Type: VarTypes.CONS,
		Data: res,
	}, RuntimeError{}
}

func Quote(expr Expression) (Var, RuntimeError) { // TODO handle expression
	if expr.Children == nil {
		return expr.Value, RuntimeError{}
	}
	return sliceToConsQuote(expr.Children)
}

func Lambda(expr *Expression, commands []*Expression) (Var, RuntimeError) {
	Params := make([]string, len(expr.Children))
	for i := range expr.Children {
		// TODO assert that the children are all symbols else error
		if expr.Children[i].Value.Type == VarTypes.SYMBOL {
			Params[i] = expr.Children[i].Value.Data.(string) // the string of the name
		}
	}
	first := &Expression{
		Parent: nil, // TODO
		Children: nil,
		Value: Var{Data: "progn", Type: VarTypes.SPECIALFORM},
		Index: 0,
		Length: 0,
	}
	sliced_first := make([]*Expression, 0)
	sliced_first = append(sliced_first, first)
	Code := &Expression{
		Parent:   nil, //TODO
		Children: append(sliced_first, commands...),
		Value:    Var{},
		Index:    0, //TODO
		Length:   0, // TODO
	}
	return Var{
		Data: Function{
			Params: Params,
			Code:   Code,
		},
		Type: VarTypes.FN,
	}, RuntimeError{}

}

func Progn(commands []*Expression, local Stack) (Var, RuntimeError) {
	for i := 0; i < len(commands)-1; i++ {
		_, err := Eval(commands[i], local)
		if err.Exists() {
			return Var{}, err
		}
	}
	return Eval(commands[len(commands)-1], local)
}

func If(condition, case1, case2 *Expression, local Stack) (Var, RuntimeError) {
	result, err := Eval(condition, local)
	if err.Exists() {
		return Var{}, err
	}
	if result.Type != VarTypes.NIL { // if true
		return Eval(case1, local)
	}
	return Eval(case2, local) // if false
}

func Def(name Expression, expr Expression, local Stack) (Var, RuntimeError) {
	if name.Value.Type != VarTypes.SYMBOL {
		return Var{}, RuntimeError{
			errors.New("No symbol given for definition"), 0, 0,
		}
	}
	str := name.Value.Data.(string)
	if _, err := local.Lookup(str); err == nil {
		fmt.Printf("Warning redefining SYMBOL %s", name)
	}
	evalled, RE := Eval(&expr, local)
	local[len(local)-1][str] = evalled
	return evalled, RE
}

func Set(name Expression, expr Expression, local Stack) (Var, RuntimeError) {
	if name.Value.Type != VarTypes.SYMBOL {
		return Var{}, RuntimeError{
			errors.New("No symbol given for definition"), 0, 0,
		}
	}
	str := name.Value.Data.(string)
	if found, err := local.Lookup(str); err == nil {
		return Var{}, RuntimeError{
			errors.New("Cannot set variable that does not exist"), 0, 0,
		}
		_ = found
	}
	evalled, RE := Eval(&expr, local)
	local[len(local)-1][str] = evalled
	return evalled, RE
}
