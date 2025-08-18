package main

import (
	"fmt"
	"os"
	"strings"

	ln_beta_reduce "github.com/gusbicalho/go-lambda/locally_nameless/beta_reduce"
	ln_expr "github.com/gusbicalho/go-lambda/locally_nameless/expr"
	ln_pretty "github.com/gusbicalho/go-lambda/locally_nameless/pretty"
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

	fmt.Println(ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName))

	for {
		redex := nextBetaRedex(expr)
		if redex == nil {
			fmt.Println(ln_pretty.ToPrettyDoc(expr).String())
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
		fmt.Println(ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName))
	}
}

func nextBetaRedex(expr ln_expr.Expr) *ln_beta_reduce.BetaRedex {
	for redex := range ln_beta_reduce.BetaRedexes(expr) {
		return &redex
	}
	return nil
}
