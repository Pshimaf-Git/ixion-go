package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"ixion/internal/token"
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

LOOP:
	for l.pos < len(l.input) {

		l.skipWhiteSpace()

		currentChar = l.peek(0)

		switch {
		case unicode.IsLetter(currentChar):
			l.tokenizeWord()
		case currentChar == '\000':
			break LOOP
		case l.isOperator(currentChar):
			tokenType, _ := token.IsOperator(currentChar)
			l.makeToken(tokenType, string(currentChar))
			currentChar = l.next()
		case unicode.IsDigit(currentChar) || currentChar == '"':
			if err := l.tokenizeLiteral(); err != nil {
				return nil, err
			}
		default:
			return nil, l.createError(UnexpectedCharacter, string(currentChar))
		}
	}

	return l.tokens, nil
}

func (l *Lexer) tokenizeWord() {
	var buffer strings.Builder
	buffer.WriteRune(l.peek(0))
	currentChar := l.next()

	for unicode.IsLetter(currentChar) || unicode.IsDigit(currentChar) {
		buffer.WriteRune(currentChar)
		currentChar = l.next()
	}
	word := buffer.String()

	if tokenType, ok := token.IsKeyword(word); ok {
		l.makeToken(tokenType, tokenType.String())
	} else if tokenType, ok := token.IsLangType(word); ok {
		l.makeToken(tokenType, tokenType.String())
	} else {
		l.makeToken(token.IDENT, word)
	}
}

func (l *Lexer) tokenizeLiteral() error {
	var buff strings.Builder

	currentChar := l.peek(0)

	switch {
	case unicode.IsDigit(currentChar):
		for unicode.IsDigit(currentChar) && currentChar != '\000' {
			buff.WriteRune(currentChar)
			currentChar = l.next()
		}
		l.makeToken(token.NUMBER_LITERAL, buff.String())
	case currentChar == '"':
		currentChar = l.next()
		for currentChar != '"' && currentChar != '\000' {
			buff.WriteRune(currentChar)
			currentChar = l.next()
		}

		if currentChar == '\000' {
			return l.createError(UnclosedStringLiteral, "string literal must be closed")
		}

		l.incPos()
		l.makeToken(token.STRING_LITERAL, buff.String())
	}

	return nil
}

func (l *Lexer) isOperator(char rune) bool {
	_, ok := token.IsOperator(char)
	return ok
}

func (l *Lexer) makeToken(_type token.TokenType, text string) {
	l.tokens.Append(token.New(_type, text))
}

func (l *Lexer) skipWhiteSpace() {
	for l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
		l.incPos()
	}
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return '\000'
	}

	l.incPos()
	return l.peek(0)
}

func (l *Lexer) incPos() {
	if l.pos >= len(l.input) {
		return
	}

	if l.input[l.pos] == '\n' {
		l.row++
		l.col = 1
	} else {
		l.col++
	}
	l.pos++
}

func (l *Lexer) peek(currentPos int) rune {
	finalPos := currentPos + l.pos
	if finalPos >= len(l.input) {
		return '\000'
	}
	return l.input[finalPos]
}

func (l *Lexer) createError(kind LexerErrorKind, msg string) error {
	return newError(kind, fmt.Sprintf("%d:%d", l.row, l.col), msg)
}
