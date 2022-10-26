package main

import (
	"errors"
	"fmt"
)

type Expression struct {
	Parent   *Expression
	Children []*Expression
	Value    Var
	Index    uint64
	Length   uint16
}

func ParseLoneToken(token Token) Expression {
	res := Expression{}
	res.Index = token.Index
	res.Length = token.Length
	switch token.Type {
	case STRING:
		res.Value.Type = VarTypes.STRING
	case NUMBER:
		res.Value.Type = VarTypes.NUM
	case SYMBOL:
		res.Value.Type = VarTypes.SYMBOL
	default:
		//		fmt.Println(token)
		panic(errors.New("Unrecognized type"))
	}
	return res
}

func Parse(tokens []Token, runes []rune) *Expression {
	fmt.Println(tokens)
	res := Expression{}
	if len(tokens) == 0 {
		panic(errors.New("0 tokens given"))
	}
	str := string(runes[tokens[0].Index : tokens[len(tokens)-1].Index+uint64(tokens[len(tokens)-1].Length)])
	fmt.Printf("Parsing: %s\n", str)
	if tokens[0].Type != LPAREN {
		if len(tokens) != 1 {
			panic(errors.New("No Left Paren but mulitple tokens found"))
		}
		res = ParseLoneToken(tokens[0])
		res.Value.Data = tokens[0].Data
	} else {
		if tokens[len(tokens)-1].Type != RPAREN {
			panic(errors.New("No matching RPAREN token given"))
		}
		res.Children = make([]*Expression, 0)
		for i := 1; i < len(tokens)-1; i++ {
			tokens[i].Print(runes)
			if tokens[i].Type == LPAREN {
				start := i
				counts := 0
				for ; ; i++ {
					//					fmt.Println(counts)
					//					tokens[i].Print(runes)
					if tokens[i].Type == LPAREN {
						counts++
					} else if tokens[i].Type == RPAREN {
						counts--
						if counts <= 1 {
							res.Children = append(res.Children, Parse(tokens[start:i+1], runes))
							break
						}
					}
				}
			} else {
				res.Children = append(res.Children, Parse(tokens[i:i+1], runes))
			}
		}
		for i, _ := range res.Children {
			res.Children[i].Parent = &res
		}
	}
	fmt.Printf("Finished %s\n", str)
	return &res
}
