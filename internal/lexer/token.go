package lexer

type TokenType string

const (
	NUMBER    TokenType = "NUMBER"
	PLUS      TokenType = "PLUS"
	MINUS     TokenType = "MINUS"
	MUL       TokenType = "MUL"
	DIV       TokenType = "DIV"
	SEMICOLON TokenType = "SEMICOLON"
	LPAREN    TokenType = "LPAREN"
	RPAREN    TokenType = "RPAREN"

	IDENTIFIER TokenType = "ID"

	VAR_KEYWORD   TokenType = "VAR"
	CONST_KEYWORD TokenType = "CONST"
	PRINT_KEYWORD TokenType = "PRINT"
)

type Token struct {
	TypeOfToken TokenType
	Text        string
}
