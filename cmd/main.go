package main

import (
	"fmt"
	"ixion/internal/lexer"
	"log"
)

const code string = `
	fn myfunc(a int) int {
		var res = (a + 12) / 2;
		return res;
	}

	var foo uint8 = 0;
	var bar = "\tTab\nNewLine";

	foo = (foo + 1000) * (12 - 1) / 3;

	print foo;
	print bar;
	print "Hello, world";
	print myfunc(foo);
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
