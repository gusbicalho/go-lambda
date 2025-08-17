package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gusbicalho/go-lambda/locally_nameless"
	"github.com/gusbicalho/go-lambda/parse_tree_to_locally_nameless"
	"github.com/gusbicalho/go-lambda/parser"
	"github.com/gusbicalho/go-lambda/tokenizer"
)

func main() {
	source := os.Args[1]
	parseTree, err := parser.Parse(tokenizer.New(strings.NewReader(source)))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	expr := parse_tree_to_locally_nameless.ToLocallyNameless(*parseTree)

	fmt.Println(locally_nameless.ToLambdaNotation(expr, locally_nameless.DisplayName))

	for {
		redex := nextBetaRedex(expr)
		if redex == nil {
			fmt.Println(locally_nameless.ToPrettyString(expr))
			fmt.Println("Irreducible.")
			break
		}
		fmt.Println(redex.ToPrettyDoc(nil).String())
		fmt.Print("Step? ")
		_, err = fmt.Scanln()
		if err != nil {
			fmt.Println(err)
			return
		}
		expr = redex.Reduce()
		fmt.Println(locally_nameless.ToLambdaNotation(expr, locally_nameless.DisplayName))
	}
}

func nextBetaRedex(expr locally_nameless.Expr) *locally_nameless.BetaRedex {
	for redex := range locally_nameless.BetaRedexes(expr) {
		return &redex
	}
	return nil
}
