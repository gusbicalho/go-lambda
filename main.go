package main

import (
	"fmt"
	"iter"
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

	ast := parse_tree_to_locally_nameless.ToLocallyNameless(*parseTree)

	fmt.Println(locally_nameless.ToLambdaNotation(ast, locally_nameless.DisplayName))
	fmt.Println(locally_nameless.ToPrettyString(ast))

	for expr := range betaReductions(ast, true) {
		fmt.Print("Step? ")
		_, err = fmt.Scanln()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(locally_nameless.ToLambdaNotation(expr, locally_nameless.DisplayName))
		fmt.Println(locally_nameless.ToPrettyString(expr))
	}
	fmt.Println("Irreducible.")
}

func betaReductions(
	expr locally_nameless.Expr,
	reduceUnderLambda bool,
) iter.Seq[locally_nameless.Expr] {
	return func(yield func(locally_nameless.Expr) bool) {
		for expr, reduced := betaReduceNext(expr, reduceUnderLambda); reduced; expr, reduced = betaReduceNext(expr, reduceUnderLambda) {
			if !yield(expr) {
				return
			}
		}
	}
}

func betaReduceNext(expr locally_nameless.Expr, reduceUnderLambda bool) (locally_nameless.Expr, bool) {
	for locus := range locally_nameless.BetaRedexes(expr) {
		return locus.Reduce(), true
	}
	return expr, false
}
