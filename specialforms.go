package main

import (
	"errors"
	"fmt"
)

func sliceToConsQuote(exprs []*Expression) (Var, RuntimeError) {
	var res Cons // TODO start with nil somehow

	for i := range exprs {
		v, Err := Quote(*exprs[i])
		if Err.Exists() {
			return Var{}, RuntimeError{
				errors.New(fmt.Sprintf("Quote Failed!!!! on %v", exprs[i])), 0, 0,
			}
		}

		res = Cons{
			v,
			Var{
				Data: res,
				Type: VarTypes.CONS,
			},
		}
	}

	return Var{
		Type: VarTypes.CONS,
		Data: res,
	}, RuntimeError{}
}

func Quote(expr Expression) (Var, RuntimeError) { // TODO handle expression
	fmt.Println("----")
	fmt.Println(expr)
	fmt.Println("----")
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
	Code := &Expression{
		Parent:   nil, //TODO
		Children: commands,
		Value:    Var{Data: "progn", Type: VarTypes.SYMBOL},
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
	local[len(local) - 1][str] = evalled
	return evalled, RE
}
