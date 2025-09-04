package token

import "strings"

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
