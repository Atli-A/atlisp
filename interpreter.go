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
	FN          VarType
	NUM         VarType
	STRING      VarType
	SYMBOL      VarType
	MACRO       VarType
	CONS        VarType
	BUILTIN     VarType
	SPECIALFORM VarType
	NIL         VarType
}{0, 1, 2, 3, 4, 5, 6, 7, 8}

type Var struct {
	Data any
	Type VarType
}

func (v Var) String() string {
	return fmt.Sprintf("%v", v.Data)
}

type Function struct {
	Params []string // if the first param is |a| then params is ["a"]
	Code   *Expression
}

type Cons struct {
	First Var
	Rest  Var
}

func (c Cons) String() string {
	return fmt.Sprintf("(%v . %v)", c.First, c.Rest)
}

var (
	SpecialFormNames = []string{
		"quote",
		"fn",
		"progn",
		"def",
		"if",
		"eval",
	}
)

type Stack []map[string]Var // the last map is the localest one

func (s Stack) Copy() Stack {
	res := make(Stack, len(s))
	copy(res, s)
	return res
}

func (s *Stack) AddLayer(m map[string]Var) {
	*s = append(*s, m)
}

func (s Stack) Lookup(name string) (Var, error) {
	if name == "nil" || name == "T" {
		return Var{Type: VarTypes.SYMBOL, Data: name}, nil
	}
	for i := len(s) - 1; i >= 0; i-- {
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

func Eval(expr *Expression, local Stack) (Var, RuntimeError) {
//	fmt.Println(expr)
	local = local.Copy()
	if expr.Children != nil {
		first, err := Eval(expr.Children[0], local)
		if err.Exists() {
			return Var{}, err
		}
		switch first.Type {
		case VarTypes.MACRO:
			fmt.Println("Macros are unsupported!")
		case VarTypes.BUILTIN:
			params := make([]Var, 0, len(expr.Children)-1)
			for i, _ := range expr.Children[1:] {
				evalled, err := Eval(expr.Children[i+1], local)
				if err.Exists() {
					return Var{}, err
				}
				params = append(params, evalled)
			}
			return first.Data.(func(...Var) (Var, RuntimeError))(params...)
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
			return Eval(fn.Code, append(pass_stack, param_layer))
		case VarTypes.SPECIALFORM:
			switch first.Data.(string) {
			case "quote":
				return Quote(*expr.Children[1])
				// TODO ensure right number of args
			case "progn":
				// TODO ensure except 0 args
				return Progn(expr.Children[1:], local)
			case "fn":
				// TODO ensure right number of args
				return Lambda(expr.Children[1], expr.Children[2:])
			case "def":
				// TODO ensure right number of args
				return Def(*expr.Children[1], *expr.Children[2], local)
			case "if":
				// TODO ensure right number of args
				return If(expr.Children[1], expr.Children[2], 
							expr.Children[3], local)

			}

		default:
			return Var{}, RuntimeError{
				errors.New("Cannot use value of not macro or function to call"), 0, 0,
			}
		}
	} else {
		switch expr.Value.Type {
		case VarTypes.SYMBOL:
			if Contains[string]([]string{"nil", "true"}, expr.Value.Data.(string)) {
				return expr.Value, RuntimeError{}
			}
			evals_to, err := local.Lookup(expr.Value.Data.(string))
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
