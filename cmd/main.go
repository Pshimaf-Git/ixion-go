package main

import (
	"fmt"
	"ixion/internal/lexer"
	"log"
)

const code string = `
	var foo uint8 = 0;
	var bar = "\tTab\nNewLine";

	foo = (foo + 1000) * (12 - 1) / 3;

	print foo;
	print bar;
	print "Hello, world";
`

func main() {
	input := []rune(code)
	l := lexer.New(input)
	res, err := l.Tokenize()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(res)
}
