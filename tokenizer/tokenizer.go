package tokenizer

import (
	"io"
	"unicode"

	"github.com/gusbicalho/go-lambda/runes_reader"
	"github.com/gusbicalho/go-lambda/token"
)

type Tokenizer struct {
	runes  *runes_reader.RunesReader
	buffer []token.Token
}

func New(r io.Reader) *Tokenizer {
	return &Tokenizer{
		runes:  runes_reader.New(r),
		buffer: make([]token.Token, 0, 1),
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
	t.buffer = t.buffer[:0]
	return tok
}

func (t *Tokenizer) Peek() token.Token {
	if len(t.buffer) > 0 {
		return t.buffer[0]
	}
	tok := t.nextFromRunes()
	t.buffer = append(t.buffer, tok)
	return tok
}

func (t *Tokenizer) nextFromRunes() token.Token {
	pos := t.runes.Pos()

	if err := t.skipWhitespace(); err != nil {
		if err == io.EOF {
			return token.EOFToken(pos)
		}
		return token.InvalidToken(err.Error(), pos)
	}

	pos = t.runes.Pos()

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
