package main

import (
	"fmt"
	"os"

	"github.com/gusbicalho/go-lambda/parser"
	"github.com/gusbicalho/go-lambda/tokenizer"
)

func main() {
	tokenizer := tokenizer.New(os.Stdin)
	tree, err := parser.Parse(tokenizer)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(tree.String())
	}
}
