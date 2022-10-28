package main

import (
	"strings"
	"bufio"
	"fmt"
	"io"
	"os"
)

func Run(line []rune, stack *Stack) any {
		tokenized, token_err := Tokenize(line)
		if token_err.Exists() {
			fmt.Println(token_err)
			fmt.Println("^^^^^^^^^^^")
			return nil
		}
//		fmt.Println(TokensToString(tokenized, code))
		parsed := Parse(tokenized, line)
		res, runtime_err := Eval(parsed, *stack)
		if runtime_err.Exists() {
			fmt.Println(runtime_err)
			fmt.Println("^^^^^^^^^^^")
			return nil
		}
		return res.Data
}

func EvenParens(str string) bool {
	return strings.Count(str, "(") == strings.Count(str, ")")
}

func Repl() {
	fmt.Println("Atlisp v0.01")
	stack := Stack{}
	stack.AddLayer(GenerateBuiltins())
	reader := bufio.NewReader(os.Stdin)
	for {
//		fmt.Println(stack)
		fmt.Print(" > ")
		even_parens := false
		var line string
		for !even_parens {
			tline, err := reader.ReadString('\n')
			line += tline
			if err != nil {
				if err == io.EOF {
					fmt.Println()
					return
				}
				fmt.Println(err)
				continue
			}
			even_parens = EvenParens(line)
			if !even_parens {
				fmt.Print(".. ")
			}
		}

//		fmt.Println(line)
		code := []rune(line)
		fmt.Println(Run(code, &stack))
	}

}

func SplitParens(str []rune) []string {
	res := make([]string, 0)
	counter := 0
	last := 0
	for i, r := range str {
		if r == '(' {
			counter++
		} else if r == ')' {
			counter--
			if counter == 0 {
				res = append(res, string(str[last:i+1]))
				last = i+1
			}
		}
	}
	return res
}

func RunFile(filename string) {
	stack := Stack{}
	stack.AddLayer(GenerateBuiltins())
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	arr := SplitParens([]rune(string(bytes)))
	for i := range arr {
//		fmt.Println(arr[i])
		Run([]rune(arr[i]), &stack)
	}
}

func main() {
	if len(os.Args) <= 1 {
		Repl()
	} else {
		RunFile(os.Args[1])
	}
}
