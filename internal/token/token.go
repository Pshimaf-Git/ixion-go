package token

import "strings"

type TokenType int

const (
	NUMBER TokenType = iota

	PLUS
	MINUS
	DIV
	MUL

	RBRACE
	LBRACE

	LPAREN
	RPAREN

	IDENT

	SEMICOLON

	FN
	FOR
	VAR
	CONST
	PRINT
)

var tokenTypes = [...]string{
	NUMBER: "NUMBER",

	PLUS:  "PLUS",
	MINUS: "MINUS",
	DIV:   "DIV",
	MUL:   "MUL",

	RBRACE: "RBRACE",
	LBRACE: "LBRACE",

	LPAREN: "LPAREN",
	RPAREN: "RPAREN",

	IDENT: "IDENT",

	SEMICOLON: "SEMICOLON",

	FN:    "FN",
	FOR:   "FOR",
	VAR:   "VAR",
	CONST: "CONST",
	PRINT: "PRINT",
}

var keywords = map[string]TokenType{
	"const": CONST,
	"var":   VAR,
	"print": PRINT,
	"fn":    FN,
	"for":   FOR,
}

var operators = map[rune]TokenType{
	'+': PLUS,
	'-': MINUS,
	'*': MUL,
	'/': DIV,

	// TODO: is operators???
	';': SEMICOLON,
	'(': LPAREN,
	')': RPAREN,
	'{': LBRACE,
	'}': RBRACE,
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

type Token struct {
	Type TokenType
	Text string
	//  TODO: add meta data for token?
	// Meta any
}

func New(_type TokenType, text string) Token {
	return Token{
		Type: _type,
		Text: text,
	}
}

func (t Token) String() string {
	var buff strings.Builder

	buff.WriteString(t.Type.String())

	if t.Text != "" {
		buff.WriteString(" " + t.Text)
	}

	/*
	   if t.Meta != nil {
	     buff.WriteString(fmt.Sprintf("%v", t.Meta))
	   }
	*/

	return buff.String()
}

// If you add metadata for [Token] uncomment this code
/*
func WithMeta(_type TokenType, text string, meta any) Token {
	return Token{
		Type: _type,
		Text: text,
		Meta: meta,
	}
}
*/

type Tokens []Token

func NewTokens() *Tokens {
	return new(Tokens)
}

func (t *Tokens) String() string {
	var buff strings.Builder

	buff.WriteByte('[')

	for _, tok := range *t {
		buff.WriteString("\n\t")
		buff.WriteString(tok.String())
	}

	buff.WriteByte('\n')
	buff.WriteByte(']')

	return buff.String()
}

func (t *Tokens) Append(token Token) {
	*t = append(*t, token)
}

func (t *Tokens) Reset() []Token {
	if len(*t) == 0 {
		return nil
	}

	res := *t

	*t = []Token{}

	return res
}
