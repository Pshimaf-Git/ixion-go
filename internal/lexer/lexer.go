package lexer

import (
	"ixion/internal/token"
	"strings"
	"unicode"
)

type Lexer struct {
	input         []rune
	tokens        *token.Tokens
	pos, row, col int
}

func New(in []rune) *Lexer {
	return &Lexer{
		input:  in,
		tokens: token.NewTokens(),
		pos:    0,
		row:    1,
		col:    1,
	}
}

func Tokenize(in string) (*token.Tokens, error) {
	lexer := Lexer{
		input:  []rune(in),
		tokens: token.NewTokens(),
	}
	return lexer.Tokenize()
}

func (l *Lexer) Tokenize() (*token.Tokens, error) {
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
		} else if _, ok := token.IsOperator(currentChar); ok {
			if tokenType, ok := token.IsOperator(currentChar); ok {
				l.makeToken(tokenType, string(currentChar))
				currentChar = l.next()
			} else {
				return nil, newError(InvalidOperator, string(currentChar))
			}
		} else {
			return nil, newError(UnexpectedCharacter, string(currentChar))
		}
	}

	return l.tokens, nil
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
	if tokenType, ok := token.IsKeyword(word); ok {
		l.makeToken(tokenType, "")
	} else {
		l.makeToken(token.IDENT, word)
	}
}

func (l *Lexer) makeToken(_type token.TokenType, text string) {
	l.tokens.Append(token.New(_type, text))
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
