package token

import (
	"fmt"

	"github.com/gusbicalho/go-lambda/position"
)

type Type int

const (
	Invalid Type = iota
	EOF
	LeftParen
	RightParen
	Lambda
	Dot
	Identifier
)

func (t Type) String() string {
	switch t {
	case Invalid:
		return "INVALID"
	case EOF:
		return "EOF"
	case LeftParen:
		return "LPAREN"
	case RightParen:
		return "RPAREN"
	case Lambda:
		return "LAMBDA"
	case Dot:
		return "DOT"
	case Identifier:
		return "IDENT"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	tokenType Type
	Value     string
	Position  position.Position
}

func (t Token) String() string {
	return fmt.Sprint(t.Type(), " at ", t.Position.Line, ":", t.Position.Column)
}

func (t Token) Type() Type {
	return t.tokenType
}

func InvalidToken(reason string, pos position.Position) Token {
	return Token{tokenType: Invalid, Value: reason, Position: pos}
}

func EOFToken(pos position.Position) Token {
	return Token{tokenType: EOF, Value: "", Position: pos}
}

func LeftParenToken(pos position.Position) Token {
	return Token{tokenType: LeftParen, Value: "(", Position: pos}
}

func RightParenToken(pos position.Position) Token {
	return Token{tokenType: RightParen, Value: ")", Position: pos}
}

func LambdaToken(pos position.Position) Token {
	return Token{tokenType: Lambda, Value: "\\", Position: pos}
}

func DotToken(pos position.Position) Token {
	return Token{tokenType: Dot, Value: ".", Position: pos}
}

func IdentifierToken(name string, pos position.Position) Token {
	return Token{tokenType: Identifier, Value: name, Position: pos}
}
