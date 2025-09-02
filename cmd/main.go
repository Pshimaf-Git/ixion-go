package main

import (
	"fmt"
	"ixion/internal/lexer"
	"log"
)

func main() {
	input := []rune("var a;\n\n\n\tprint (a+b)*c")
	l := lexer.New(input)
	res, err := l.Tokenize()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(res)
}
