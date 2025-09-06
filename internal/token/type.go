package token

type TokenType int

const (
	// Types:

	// Signed
	INT TokenType = iota
	INT8
	INT16
	INT32
	INT64

	// Unsigned
	UINT
	UINT8
	UINT16
	UINT32
	UINT64

	NUMBER_LITERAL // var a = 1234;

	// String
	STRING

	STRING_LITERAL // var a = "STRING_LITERAL";

	PLUS
	MINUS
	DIV
	MUL

	ASSIGN

	RBRACE
	LBRACE

	LPAREN
	RPAREN

	IDENT

	SEMICOLON
	COMMA

	FN
	FOR
	VAR
	CONST
	PRINT
	RETURN

	EOF
)

var tokenTypes = [...]string{
	INT:   "INT",
	INT8:  "INT8",
	INT16: "INT16",
	INT32: "INT32",
	INT64: "INT64",

	UINT:   "UINT",
	UINT8:  "UINT8",
	UINT16: "UINT16",
	UINT32: "UINT32",
	UINT64: "UINT64",

	NUMBER_LITERAL: "NUMBER_LITERAL",

	STRING: "STRING",

	STRING_LITERAL: "STRING_LITERAL",

	PLUS:  "PLUS",
	MINUS: "MINUS",
	DIV:   "DIV",
	MUL:   "MUL",

	ASSIGN: "ASSIGN",

	RBRACE: "RBRACE",
	LBRACE: "LBRACE",

	LPAREN: "LPAREN",
	RPAREN: "RPAREN",

	IDENT: "IDENT",

	SEMICOLON: "SEMICOLON",
	COMMA:     "COMMA",

	FN:     "FN",
	FOR:    "FOR",
	VAR:    "VAR",
	CONST:  "CONST",
	PRINT:  "PRINT",
	RETURN: "RETURN",

	EOF: "EOF",
}

var keywords = map[string]TokenType{
	"const":  CONST,
	"var":    VAR,
	"print":  PRINT,
	"fn":     FN,
	"for":    FOR,
	"return": RETURN,
}

var operators = map[rune]TokenType{
	'+': PLUS,
	'-': MINUS,
	'*': MUL,
	'/': DIV,

	'=': ASSIGN,

	// TODO: is operators???
	';': SEMICOLON,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
	',': COMMA,
}

var types = map[string]TokenType{
	"int":   INT,
	"int8":  INT8,
	"int16": INT16,
	"int32": INT32,
	"int64": INT64,

	"uint":   UINT,
	"uint8":  UINT8,
	"uint16": UINT16,
	"uint32": UINT32,
	"uint64": UINT64,

	"string": STRING,
}

func (tt TokenType) String() string {
	var s string
	if int(tt) <= len(tokenTypes)-1 {
		s = tokenTypes[tt]
	}

	return s
}

func (tt TokenType) Is(other TokenType) bool {
	return tt == other
}

func (tt *TokenType) All(others ...TokenType) bool {
	for i := range others {
		if !tt.Is(others[i]) {
			return false
		}
	}

	return true
}

func (tt TokenType) IsKeyword() bool {
	_, ok := keywords[tt.String()]
	return ok
}

func (tt TokenType) Valid() bool {
	return int(tt) <= len(tokenTypes)-1
}

func IsOperator(char rune) (TokenType, bool) {
	tt, ok := operators[char]
	return tt, ok
}

func IsKeyword(s string) (TokenType, bool) {
	tt, ok := keywords[s]
	return tt, ok
}

func IsLangType(s string) (TokenType, bool) {
	tt, ok := types[s]
	return tt, ok
}

func GetTokenTypeFromString(s string) (TokenType, bool) {
	tt, ok := types[s]
	return tt, ok
}

func (t Token) IsType() bool {
	switch t.Type {
	case INT, INT8, INT16, INT32, INT64,
		UINT, UINT8, UINT16, UINT32, UINT64,
		STRING:
		return true
	default:
		return false
	}
}
