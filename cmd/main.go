package main

import (
	"encoding/json"
	"fmt"

	"ixion/internal/lexer"
	"ixion/internal/parser"
)

// Example code
const code string = `
	fn f(a int) int {
		var res = (a + 12) / 2;
		return res;
	}

	var funcLit = fn(a int) int {return a + 1}

	var foo uint8 = 0;
	var bar string = "\tTab\nNewLine";

	foo = (foo + 1000) * (12 - 1) / 3;

	print(foo);
	print(bar);
	print("Hello, world");
	print(myfunc(foo));
`

func main() {
	l := lexer.New([]rune(code))
	toks, _ := l.Tokenize()

	p := parser.New(toks)
	program := p.ParseProgram()

	jsonOutput, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling to JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonOutput))
}
