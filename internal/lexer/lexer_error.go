package lexer

import "fmt"

type LexerErrorKind string

const (
	// TODO
	Default         LexerErrorKind = "lexer error: %s"
	InvalidOperator LexerErrorKind = "invalid operator %s"
)

type LexerError struct {
	Kind    LexerErrorKind
	Message string
}

func (l *LexerError) Error() string {
	if l.Message == "" {
		return fmt.Sprintf("lexer error: %v", l.Kind)
	}
	return l.Message
}
