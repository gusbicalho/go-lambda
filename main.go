package main

import (
	"fmt"
	"os"

	"github.com/gusbicalho/go-lambda/parse_tree_to_locally_nameless"
	"github.com/gusbicalho/go-lambda/parser"
	"github.com/gusbicalho/go-lambda/stack"
	"github.com/gusbicalho/go-lambda/tokenizer"
)

func main() {
	tokenizer := tokenizer.New(os.Stdin)
	parseTree, err := parser.Parse(tokenizer)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Parse tree")
	fmt.Println(parseTree.String())

	ast := parse_tree_to_locally_nameless.ToLocallyNameless(*parseTree)

	fmt.Println("Locally nameless")
	fmt.Println(ast.ToPrettyDoc(stack.Empty[string]()).String())
}
