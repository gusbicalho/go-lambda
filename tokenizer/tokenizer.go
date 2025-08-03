package tokenizer

import (
	"io"
	"unicode"

	"github.com/gusbicalho/go-lambda/runes_reader"
	"github.com/gusbicalho/go-lambda/token"
)

type Tokenizer struct {
	runes   *runes_reader.RunesReader
	current struct {
		ready bool
		token token.Token
	}
}

type currentToken struct {
	ready bool
	token token.Token
}

func New(r io.Reader) *Tokenizer {
	return &Tokenizer{
		runes:   runes_reader.New(r),
		current: currentToken{ready: false},
	}
}

func (t *Tokenizer) Each(action func(token token.Token) error) error {
	for tok := t.Next(); ; tok = t.Next() {
		if err := action(tok); err != nil {
			return err
		}
		if tok.Type() == token.EOF {
			return nil
		}
	}
}

func (t *Tokenizer) Next() token.Token {
	tok := t.Peek()
	t.current = currentToken{ready: false}
	return tok
}

func (t *Tokenizer) Peek() token.Token {
	if t.current.ready {
		return t.current.token
	}
	tok := t.nextFromRunes()
	t.current = currentToken{
		ready: true,
		token: tok,
	}
	return tok
}

func (t *Tokenizer) nextFromRunes() token.Token {
	pos := t.pos()

	if err := t.skipWhitespace(); err != nil {
		if err == io.EOF {
			return token.EOFToken(pos)
		}
		return token.InvalidToken(err.Error(), pos)
	}

	pos = t.pos()

	r, err := t.runes.Peek()
	if err != nil {
		if err == io.EOF {
			return token.EOFToken(pos)
		}
		return token.InvalidToken(err.Error(), pos)
	}

	switch r {
	case '(':
		t.runes.Consume()
		return token.LeftParenToken(pos)
	case ')':
		t.runes.Consume()
		return token.RightParenToken(pos)
	case '\\':
		t.runes.Consume()
		return token.LambdaToken(pos)
	case '.':
		t.runes.Consume()
		return token.DotToken(pos)
	default:
		if unicode.IsLetter(r) || r == '_' {
			value, err := t.readIdentifier()
			if err != nil {
				return token.InvalidToken(err.Error(), pos)
			}
			return token.IdentifierToken(value, pos)
		}

		t.runes.Consume()
		return token.InvalidToken(string(r), pos)
	}
}

func (t *Tokenizer) pos() token.Position {
	pos := t.runes.Pos()
	return token.Position{Line: pos.Line, Column: pos.Column}
}

func (t *Tokenizer) skipWhitespace() error {
	for {
		r, err := t.runes.Peek()
		if err != nil {
			return err
		}

		if !unicode.IsSpace(r) {
			break
		}

		t.runes.Consume()
	}
	return nil
}

func (t *Tokenizer) readIdentifier() (string, error) {
	var result []rune

	for {
		r, err := t.runes.Peek()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if !unicode.IsLetter(r) && r != '_' {
			break
		}

		result = append(result, r)
		t.runes.Consume()
	}

	return string(result), nil
}
