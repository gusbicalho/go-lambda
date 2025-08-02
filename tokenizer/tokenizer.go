package tokenizer

import (
	"bufio"
	"io"
	"unicode"

	"github.com/gusbicalho/go-lambda/token"
)

type Tokenizer struct {
	reader  *bufio.Reader
	line    int
	column  int
	current struct {
		ready bool
		rune  rune
	}
}

func New(r io.Reader) *Tokenizer {
	return &Tokenizer{
		reader: bufio.NewReader(r),
		line:   1,
		column: 0,
	}
}

func (t *Tokenizer) Next() token.Token {
	pos := t.pos()

	if err := t.skipWhitespace(); err != nil {
		if err == io.EOF {
			return token.EOFToken(pos)
		}
		return token.InvalidToken(err.Error(), pos)
	}

	pos = t.pos()

	r, err := t.peekRune()
	if err != nil {
		if err == io.EOF {
			return token.EOFToken(pos)
		}
		return token.InvalidToken(err.Error(), pos)
	}

	switch r {
	case '(':
		t.consumeRune()
		return token.LeftParenToken(pos)
	case ')':
		t.consumeRune()
		return token.RightParenToken(pos)
	case '\\':
		t.consumeRune()
		return token.LambdaToken(pos)
	case '.':
		t.consumeRune()
		return token.DotToken(pos)
	default:
		if unicode.IsLetter(r) || r == '_' {
			value, err := t.readIdentifier()
			if err != nil {
				return token.InvalidToken(err.Error(), pos)
			}
			return token.IdentifierToken(value, pos)
		}

		t.consumeRune()
		return token.InvalidToken(string(r), pos)
	}
}

func (t *Tokenizer) pos() token.Position {
	return token.Position{Line: t.line, Column: t.column}
}

func (t *Tokenizer) peekRune() (rune, error) {
	if t.current.ready {
		return t.current.rune, nil
	}

	r, _, err := t.reader.ReadRune()
	if err != nil {
		return 0, err
	}

	t.current.ready = true
	t.current.rune = r
	return r, nil
}

func (t *Tokenizer) consumeRune() {
	if !t.current.ready {
		return
	}

	if t.current.rune == '\n' {
		t.line++
		t.column = 0
	} else {
		t.column++
	}

	t.current.ready = false
}

func (t *Tokenizer) skipWhitespace() error {
	for {
		r, err := t.peekRune()
		if err != nil {
			return err
		}

		if !unicode.IsSpace(r) {
			break
		}

		t.consumeRune()
	}
	return nil
}

func (t *Tokenizer) readIdentifier() (string, error) {
	var result []rune

	for {
		r, err := t.peekRune()
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
		t.consumeRune()
	}

	return string(result), nil
}
