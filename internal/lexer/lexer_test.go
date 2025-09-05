package lexer_test

import (
	"ixion/internal/lexer"
	"ixion/internal/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLexer_Tokinaze(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		want    []token.Token
		wantErr bool
		errText string
	}{
		{
			name:  "simplest code",
			input: "var a = 1;",
			want: []token.Token{
				token.New(token.VAR, token.VAR.String()),
				token.New(token.IDENT, "a"),
				token.New(token.ASSIGN, string('=')),
				token.New(token.NUMBER_LITERAL, "1"),
				token.New(token.SEMICOLON, string(';')),
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New([]rune(tt.input))

			got, err := l.Tokenize()

			if tt.wantErr {
				assert.Error(t, err, "want a nil error, has: ", err)

				assert.EqualError(t, err, tt.errText)
			} else {
				require.NoError(t, err, "want a non nil error")
			}

			assert.Equal(t, tt.want, got.Reset())
		})
	}
}
