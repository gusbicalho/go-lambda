package runes_reader

import (
	"bufio"
	"io"

	"github.com/gusbicalho/go-lambda/position"
)

type RunesReader struct {
	reader      *bufio.Reader
	pos         position.Position
	currentRune struct {
		ready bool
		rune  rune
	}
}

func New(r io.Reader) *RunesReader {
	return &RunesReader{
		reader: bufio.NewReader(r),
		pos: position.Position{
			Line:   1,
			Column: 0,
		},
	}
}

func (t *RunesReader) Pos() position.Position {
	return t.pos
}

func (t *RunesReader) Peek() (rune, error) {
	if t.currentRune.ready {
		return t.currentRune.rune, nil
	}

	r, _, err := t.reader.ReadRune()
	if err != nil {
		return 0, err
	}

	t.currentRune.ready = true
	t.currentRune.rune = r
	return r, nil
}

func (t *RunesReader) Consume() {
	if !t.currentRune.ready {
		return
	}

	if t.currentRune.rune == '\n' {
		t.pos.Line++
		t.pos.Column = 0
	} else {
		t.pos.Column++
	}

	t.currentRune.ready = false
}
