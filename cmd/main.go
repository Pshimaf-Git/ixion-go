package main

import (
	"encoding/json"
	"fmt"

	"ixion/internal/lexer"
	"ixion/internal/parser"
	"ixion/internal/semantic" // Добавляем импорт
)

// Example code
const code string = `
	var a = 1;
	var b = a + 2;
	print(b);
	
	fn test(x int) int {
		return x + 1;
	}
	
	var result = test(5);
	print(result);
`

func main() {
	l := lexer.New([]rune(code))
	toks, err := l.Tokenize()
	if err != nil {
		fmt.Printf("Lexer error: %v\n", err)
		return
	}

	p := parser.New(toks)
	program := p.ParseProgram()

	// Проверяем ошибки парсера
	if len(p.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, err := range p.Errors() {
			fmt.Printf("  %s\n", err)
		}
		return
	}

	// Запускаем семантический анализ
	analyzer := semantic.NewAnalyzer()
	errors := analyzer.Analyze(program)

	if len(errors) > 0 {
		fmt.Println("Semantic errors:")
		for _, err := range errors {
			fmt.Printf("  %s\n", err)
		}
		return
	}

	fmt.Println("Semantic analysis passed successfully!")

	// Выводим AST (опционально)
	jsonOutput, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling to JSON: %v\n", err)
		return
	}
	fmt.Println("AST:")
	fmt.Println(string(jsonOutput))
}
