package lexer

import "fmt"

type LexerErrorKind int

const (
	// TODO
	Default LexerErrorKind = iota
	InvalidOperator
	UnexpectedCharacter
)

var kinds = map[LexerErrorKind]string{
	Default:             "Default",
	InvalidOperator:     "Invalid Operator",
	UnexpectedCharacter: "Unexpected Character",
}

func (k LexerErrorKind) String() string {
	return kinds[k]
}

type LexerError struct {
	Kind    LexerErrorKind
	Message string
}

func newError(kind LexerErrorKind, msg string) error {
	return &LexerError{
		Kind:    kind,
		Message: msg,
	}
}

func (l *LexerError) Error() string {
	if l.Message == "" {
		return fmt.Sprintf("lexer error: %s", l.Kind.String())
	}
	return fmt.Sprintf("lexer error: %s: %s", l.Kind.String(), l.Message)
}
