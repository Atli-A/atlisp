package main

import (
	"errors"
	"fmt"
)

type RuntimeError struct {
	Err    error
	Index  uint64
	Length uint16
}

func (re RuntimeError) Exists() bool {
	return re.Err != nil
}

type VarType int32

var VarTypes = struct {
	FN     VarType
	NUM    VarType
	STRING VarType
	SYMBOL VarType
	MACRO  VarType
	CONS   VarType
}{0, 1, 2, 3, 4, 5}

type Var struct {
	Data any
	Type VarType
}

type Function struct {
	Params []string // if the first param is |a| then params is ["a"]
	Code   *Expression
}

var (
	SpecialForms = map[string]any{
		"quote": func(v Var) (Var, RuntimeError) {
			return v, RuntimeError{}
		},
		"lambda": func(expr *Expression, commands []*Expression) (Var, RuntimeError) {
			Params := make([]string, len(expr.Children))
			for i := range expr.Children {
				// TODO assert that the children are all symbols else error
				Params[i] = expr.Children[i].Value.Data.(string) // the string of the name
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

		},
		"progn": func(commands []*Expression, local Stack) (Var, RuntimeError) {
			for i := 0; i < len(commands)-1; i++ {
				_, err := Eval(commands[i], local)
				if err.Exists() {
					return Var{}, err
				}
			}
			return Eval(commands[len(commands)-1], local)
		},
	}
)

type Stack []map[string]Var // the last map is the localest one

func (s Stack) Copy() Stack {
	res := make(Stack, len(s))
	copy(res, s)
	return res
}

func (s Stack) AddLayer(m map[string]Var) {
	s = append(s, m)
}

func (s Stack) Lookup(name string) (Var, error) {
	for i := len(s) - 1; i >= 0; i++ {
		_, ok := s[i][name]
		if ok {
			return s[i][name], nil
		}
	}
	return Var{}, errors.New(fmt.Sprintf("Variable/Function %s not found", name))
}

var (
	GlobalStack *Stack
)

func Init() {
	// Add builtins
	// Add other stdlib

}

func Eval(expr *Expression, local Stack) (Var, RuntimeError) {
	fmt.Println(expr)
	local = local.Copy()
	if expr.Children != nil {
		first, err := Eval(expr.Children[0], local)
		if err.Exists() {
			return Var{}, err
		}
		switch first.Type {
		case VarTypes.MACRO:
			fmt.Println("Macros are unsupported!")
		case VarTypes.FN:
			fn := first.Data.(Function)
			if len(fn.Params) != len(expr.Children)-1 {
				return Var{}, RuntimeError{
					errors.New("Wrong number of parameters for function"), 0, 0,
				}
			}
			param_layer := make(map[string]Var)
			for i, name := range fn.Params {
				evalled, err := Eval(expr.Children[i+1], local)
				if err.Exists() {
					return Var{}, err
				}
				param_layer[name] = evalled
			}
			pass_stack := local.Copy()
			pass_stack.AddLayer(param_layer)
			return Eval(fn.Code, pass_stack)
		default:
			return Var{}, RuntimeError{
				errors.New("Cannot use value of not macro or function to call"), 0, 0,
			}
		}

	} else {
		switch expr.Value.Type {
		case VarTypes.SYMBOL:
			evals_to, err := local.Lookup(expr.Value.Data.(string))
			fmt.Println("----")
			fmt.Println(evals_to)
			if err != nil {
				return Var{}, RuntimeError{err, 0, 0}
			}
			return evals_to, RuntimeError{}
		default:
			return expr.Value, RuntimeError{}
		}
	}

	return Var{}, RuntimeError{}
}
