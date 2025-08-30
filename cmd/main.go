package main

import (
	"fmt"
	"ixion/internal/lexer"
	"log"
)

func main() {
	input := []rune("print(a + b);")
	l := lexer.New(input)
	res, err := l.Tokenize()
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(res)
}
