package parse_tree_to_locally_nameless

import (
	"github.com/gusbicalho/go-lambda/locally_nameless/expr"
	"github.com/gusbicalho/go-lambda/parse_tree"
	"github.com/gusbicalho/go-lambda/stack"
)

func ToLocallyNameless(parsed parse_tree.ParseTree) expr.Expr {
	return toLocallyNameless(parsed, stack.Empty[string]())
}

func toLocallyNameless(parsed parse_tree.ParseTree, bound stack.Stack[string]) expr.Expr {
	switch item := parsed.Item.(type) {
	case parse_tree.Parens:
		return toLocallyNameless(item.Child, bound)
	case parse_tree.Var:
		for index, boundName := range bound.IndexedItems() {
			if boundName == item.Name {
				return expr.NewBound(index)
			}
		}
		return expr.NewFree(item.Name)
	case parse_tree.Lambda:
		return expr.NewLambda(
			item.ArgName,
			toLocallyNameless(item.Body, bound.Push(item.ArgName)),
		)
	case parse_tree.App:
		app := expr.NewApp(
			toLocallyNameless(item.Callee, bound),
			toLocallyNameless(item.Args.First, bound),
		)
		for _, arg := range item.Args.More {
			app = expr.NewApp(app, toLocallyNameless(arg, bound))
		}
		return app
	default:
		panic("unknown parse tree")
	}
}
