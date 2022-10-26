package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func Repl() {
	stack := Stack{}
	stack.AddLayer(GenerateBuiltins())
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("hi")
	for {
		fmt.Println(stack)
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return
			}
			fmt.Println(err)
			fmt.Println("^^^^^^^^^^^")
			continue
		}
		fmt.Println(line)
		code := []rune(line)
		tokenized, token_err := Tokenize(code)
		if token_err.Exists() {
			fmt.Println(err)
			fmt.Println("^^^^^^^^^^^")
			continue
		}
		fmt.Println(TokensToString(tokenized, code))
		parsed := Parse(tokenized, code)
		res, runtime_err := Eval(parsed, stack)
		if runtime_err.Exists() {
			fmt.Println(err)
			fmt.Println("^^^^^^^^^^^")
			continue
		}
		fmt.Println(res.Data)

	}

}

func Test() {
	bytes, err := os.ReadFile("test.mylisp")
	if err != nil {
		panic(err)
	}
	code := []rune(string(bytes))
	tokenized, _ := Tokenize(code)
	fmt.Println("\n")
	//	for _, t := range tokenized {
	//		fmt.Println(string(code[t.Index : t.Index+uint64(t.Length)]))
	//	}

	parsed := Parse(tokenized, code)
	fmt.Println(parsed)

	fmt.Println(Eval(parsed, *new(Stack)))

}

func main() {
	Repl()
}
