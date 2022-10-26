package main

import (
	"fmt"
	"os"
)

func main() {
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
