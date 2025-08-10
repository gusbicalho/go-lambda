package main

import (
	"fmt"
	"iter"
	"os"
	"strings"

	"github.com/gusbicalho/go-lambda/locally_nameless"
	"github.com/gusbicalho/go-lambda/parse_tree_to_locally_nameless"
	"github.com/gusbicalho/go-lambda/parser"
	"github.com/gusbicalho/go-lambda/stack"
	"github.com/gusbicalho/go-lambda/tokenizer"
)

func main() {
	source := os.Args[1]
	tokenizer := tokenizer.New(strings.NewReader(source))
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

	for expr := range betaReductions(ast, true) {
		fmt.Print("Step? ")
		// i := ""
		_, err = fmt.Scanln()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(expr.ToPrettyDoc(stack.Empty[string]()).String())
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
	switch expr := expr.(type) {
	case locally_nameless.App:
		switch callee := expr.Callee().(type) {
		case locally_nameless.Lambda:
			return locally_nameless.BetaReduce(callee, expr.Arg()), true
		default:
			callee, reduced := betaReduceNext(callee, reduceUnderLambda)
			if reduced {
				return locally_nameless.NewApp(callee, expr.Arg()), true
			}
			arg, reduced := betaReduceNext(expr.Arg(), reduceUnderLambda)
			if reduced {
				return locally_nameless.NewApp(expr.Callee(), arg), true
			}
			return expr, false
		}
	case locally_nameless.Lambda:
		if reduceUnderLambda {
			reducedBody, reduced := betaReduceNext(expr.Body(), reduceUnderLambda)
			if reduced {
				return locally_nameless.NewLambda(expr.ArgName(), reducedBody), true
			}
		}
		return expr, false
	default:
		return expr, false
	}
}
