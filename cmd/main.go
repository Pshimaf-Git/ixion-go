package main

import (
	"errors"
	"fmt"
	"log"

	"ixion/internal/lexer"
	"ixion/internal/parser"
)

const code string = `
	fn myfunc(a int) int {
		var res = (a + 12) / 2;
		return res;
	}

	var foo uint8 = 0;
	var bar = "\tTab\nNewLine";

	foo = (foo + 1000) * (12 - 1) / 3;

	print(foo);
	print(bar);
	print("Hello, world");
	print(myfunc(foo););
`

func main() {
	input := []rune(code)
	l := lexer.New(input)
	tokens, err := l.Tokenize()
	if err != nil {
		log.Fatal(err.Error())
	}

	p := parser.New(tokens)

	prog := p.ParseProgram()

	if errs := p.Errors(); len(errs) != 0 {
		log.Fatal(errors.Join(errs...))
	}

	fmt.Println(prog.String())
}
