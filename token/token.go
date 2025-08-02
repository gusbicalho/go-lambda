package token

type TokenType int

const (
	Invalid TokenType = iota
	EOF
	LeftParen
	RightParen
	Lambda
	Dot
	Identifier
)

func (t TokenType) String() string {
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

type Position struct {
	Line   int
	Column int
}

type Token struct {
	tokenType TokenType
	Value     string
	Position  Position
}

func (t Token) Type() TokenType {
	return t.tokenType
}

func InvalidToken(reason string, pos Position) Token {
	return Token{tokenType: Invalid, Value: reason, Position: pos}
}

func EOFToken(pos Position) Token {
	return Token{tokenType: EOF, Value: "", Position: pos}
}

func LeftParenToken(pos Position) Token {
	return Token{tokenType: LeftParen, Value: "(", Position: pos}
}

func RightParenToken(pos Position) Token {
	return Token{tokenType: RightParen, Value: ")", Position: pos}
}

func LambdaToken(pos Position) Token {
	return Token{tokenType: Lambda, Value: "\\", Position: pos}
}

func DotToken(pos Position) Token {
	return Token{tokenType: Dot, Value: ".", Position: pos}
}

func IdentifierToken(name string, pos Position) Token {
	return Token{tokenType: Identifier, Value: name, Position: pos}
}
