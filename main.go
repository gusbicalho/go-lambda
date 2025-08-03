package main

import (
	"fmt"
	"os"

	"github.com/gusbicalho/go-lambda/token"
	"github.com/gusbicalho/go-lambda/tokenizer"
)

func main() {
	tokenizer.New(os.Stdin).Each(func(tok token.Token) error {
		switch tok.Type() {
		case token.EOF:
		case token.Invalid:
			fmt.Printf("ERROR at %d:%d: %s\n", tok.Position.Line, tok.Position.Column, tok.Value)
		default:
			fmt.Printf("%s: %q at %d:%d\n", tok.Type(), tok.Value, tok.Position.Line, tok.Position.Column)
		}
		return nil
	})
}
