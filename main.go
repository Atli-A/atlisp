package main

import (
	"fmt"
	"bufio"
	"os"
)

func Repl() {
	stack := *new(Stack)
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("hi")
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(line)
		code := []rune(line)
		tokenized, token_err := Tokenize(code)
		if token_err.Exists() {
			fmt.Println(err)
			continue
		}
		fmt.Print(">>")
		fmt.Println(tokenized)
		parsed := Parse(tokenized, code)
		res, runtime_err := Eval(parsed, stack)
		if runtime_err.Exists() {
			fmt.Println(err)
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
