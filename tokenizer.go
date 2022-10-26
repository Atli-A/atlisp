package main

import (
	"errors"
	"strconv"
	"unicode"

	"fmt"
)

type TokenType int32

const (
	NONE      TokenType = iota
	LPAREN    TokenType = iota
	RPAREN    TokenType = iota
	QUOTE     TokenType = iota
	BACKQUOTE TokenType = iota
	COMMA     TokenType = iota
	STRING    TokenType = iota
	SYMBOL    TokenType = iota
	NUMBER    TokenType = iota
)

type Token struct {
	Type   TokenType
	Data   any
	Index  uint64
	Length uint16
}

func (T Token) Print(runes []rune) {
	fmt.Println(string(runes[T.Index : T.Index+uint64(T.Length)]))
}

type TokenizationError struct {
	Err    error
	Start  uint64
	Length uint16
}

func (TE TokenizationError) Print(runes []rune) {
	fmt.Println(string(runes[TE.Start : TE.Start+uint64(TE.Length)]))
}

func IdentifySingleCharTokens(c rune, index uint64) (Token, bool) {
	res := Token{
		Type:   NONE,
		Data:   nil,
		Index:  index,
		Length: 1,
	}
	switch c {
	case '(':
		res.Type = LPAREN
	case ')':
		res.Type = RPAREN
	case '\'':
		res.Type = QUOTE
	case '`':
		res.Type = BACKQUOTE
	case ',':
		res.Type = COMMA
	default:
		return res, false
	}
	return res, true
}

func IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func Contains[T comparable](slice []T, item T) bool {
	for i, _ := range slice {
		if slice[i] == item {
			return true
		}
	}
	return false
}

// returns the index where it was first not in list
func ReadUntilFalse(str []rune, f func(rune) bool) int {
	for i := 0; i < len(str); i++ {
		if !f(str[i]) {
			return i
		}
	}
	return -1
}

func Tokenize(runes []rune) ([]Token, TokenizationError) {
	res := make([]Token, 0)
	for i := uint64(0); i < uint64(len(runes)); i++ {
		c := runes[i]
		T, found := IdentifySingleCharTokens(c, i)
		switch {
		case c == ';':
			for runes[i] != '\n' {
				i++
			}
		case IsWhitespace(c):

		case found:
			res = append(res, T)
		// TODO bounds check

		case unicode.IsDigit(c) || (c == '-' && unicode.IsDigit(runes[i+1])):
			var strint int
			var str string
			if c == '-' {
				strint = ReadUntilFalse(runes[i+1:], unicode.IsDigit)
				str = string(runes[i : i+uint64(strint)+1])
			} else {
				strint = ReadUntilFalse(runes[i:], unicode.IsDigit)
				str = string(runes[i : i+uint64(strint)])
			}

			num, err := strconv.ParseInt(str, 10, 64)
			if err == nil {
				res = append(res, Token{NUMBER, num, i, uint16(strint)})
				i += uint64(strint)
			} else {
				fmt.Println("illegal number starting symbol?")
			}

		case c == '"':
			n := ReadUntilFalse(runes[i:],
				func(r rune) bool {
					return r == '"'
				})
			if n == -1 {
				return nil, TokenizationError{
					Err:    errors.New("No matching \" found"),
					Start:  i,
					Length: uint16(len(runes[i:])),
				}
			}
			res = append(res, Token{STRING, string(runes[i+1 : i+1+uint64(n)]), i, uint16(n + 1)})
			i += uint64(n)

		default: // SYMBOL
			n := ReadUntilFalse(runes[i:], func(r rune) bool {
				return unicode.IsDigit(r) || unicode.IsLetter(r)
			})
			res = append(res, Token{SYMBOL, string(runes[i : i+uint64(n)]), i, uint16(n)})
			i += uint64(n)
		}
	}
	return res, TokenizationError{}
}
