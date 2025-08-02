package main

import (
	"fmt"
	"os"

	"github.com/gusbicalho/go-lambda/token"
	"github.com/gusbicalho/go-lambda/tokenizer"
)

func main() {
	tokenizer := tokenizer.New(os.Stdin)
	for tok := tokenizer.Next(); tok.Type() != token.EOF; tok = tokenizer.Next() {
		switch tok.Type() {
		case token.Invalid:
			fmt.Printf("ERROR at %d:%d: %s\n", tok.Position.Line, tok.Position.Column, tok.Value)
		default:
			fmt.Printf("%s: %q at %d:%d\n", tok.Type(), tok.Value, tok.Position.Line, tok.Position.Column)
		}
	}
}
