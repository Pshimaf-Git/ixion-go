package lexer

import "fmt"

type LexerErrorKind int

const (
	// TODO
	Default LexerErrorKind = iota
	InvalidOperator
	UnexpectedCharacter
	UnclosedStringLiteral
)

var kinds = map[LexerErrorKind]string{
	Default:               "Default",
	InvalidOperator:       "Invalid Operator",
	UnexpectedCharacter:   "Unexpected Character",
	UnclosedStringLiteral: "Unclosed String Literal",
}

func (k LexerErrorKind) String() string {
	return kinds[k]
}

type LexerError struct {
	Kind    LexerErrorKind
	Pos     string
	Message string
}

func newError(kind LexerErrorKind, pos string, msg string) error {
	return &LexerError{
		Kind:    kind,
		Pos:     pos,
		Message: msg,
	}
}

func (l *LexerError) Error() string {
	if l.Message == "" {
		return fmt.Sprintf("lexer error at %s: %s", l.Pos, l.Kind.String())
	}
	return fmt.Sprintf("lexer error at %s: %s: %s", l.Pos, l.Kind.String(), l.Message)
}
