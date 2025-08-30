package lexer

import (
	"strings"
	"unicode"
)

type Lexer struct {
	input         []rune
	pos, row, col int
}

var (
	Keywords = map[string]TokenType{
		"var":   VAR_KEYWORD,
		"const": CONST_KEYWORD,
		"print": PRINT_KEYWORD,
	}

	Operators = []rune{
		'+', '-', '*', '/', '(', ')', ';',
	}
)

var tokens []Token

func New(in []rune) *Lexer {

	return &Lexer{
		input: in,
		pos:   0,
		row:   1,
		col:   1,
	}
}

func Tokenize(in string) ([]Token, error) {
	lexer := Lexer{
		input: []rune(in),
	}
	return lexer.Tokenize()
}

func (l *Lexer) Tokenize() ([]Token, error) {
	var currentChar rune

	for l.pos < len(l.input) {

		for unicode.IsSpace(l.peek(0)) {
			l.skip()
		}

		currentChar = l.peek(0)

		if unicode.IsLetter(currentChar) {
			l.tokenizeWord()
		} else if currentChar == '\000' {
			break
		} else if l.isOperator(currentChar) {
			operatorTypes := map[rune]TokenType{
				'+': PLUS,
				'-': MINUS,
				'*': MUL,
				'/': DIV,
				'(': LPAREN,
				')': RPAREN,
				';': SEMICOLON,
			}

			if tokenType, exists := operatorTypes[currentChar]; exists {
				l.makeToken(tokenType, string(currentChar))
				currentChar = l.next()
			} else {
				return nil, &LexerError{Kind: Default, Message: "invalid operator"}
			}
		} else {
			return nil, &LexerError{Kind: Default, Message: "unexpected character"}
		}
	}
	return tokens, nil
}

func (l *Lexer) tokenizeWord() {
	var buffer strings.Builder
	buffer.WriteRune(l.peek(0))
	currentChar := l.next()

	for unicode.IsLetter(currentChar) {
		buffer.WriteRune(currentChar)
		currentChar = l.next()
	}
	word := buffer.String()
	if tokenType, exists := Keywords[word]; exists {
		l.makeToken(tokenType, "")
	} else {
		l.makeToken(IDENTIFIER, word)
	}
}

func (l *Lexer) isOperator(target rune) bool {
	for _, r := range Operators {
		if r == target {
			return true
		}
	}
	return false
}

func (l *Lexer) makeToken(typo TokenType, text string) {
	tokens = append(tokens, Token{
		TypeOfToken: typo,
		Text:        text,
	})
}

func (l *Lexer) skip() {
	if l.pos >= len(l.input) {
		return
	}
	result := l.input[l.pos]

	if result == '\n' {
		l.row++
		l.col = 1
	} else {
		l.col++
	}
	l.pos++
}

func (l *Lexer) next() rune {
	l.skip()
	return l.peek(0)
}

func (l *Lexer) peek(currentPos int) rune {
	finalPos := currentPos + l.pos
	if finalPos >= len(l.input) {
		return '\000'
	}
	return l.input[finalPos]
}
