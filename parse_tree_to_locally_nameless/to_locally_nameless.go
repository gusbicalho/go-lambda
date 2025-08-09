package parse_tree_to_locally_nameless

import (
	"github.com/gusbicalho/go-lambda/locally_nameless"
	"github.com/gusbicalho/go-lambda/parse_tree"
	"github.com/gusbicalho/go-lambda/stack"
)

func ToLocallyNameless(parsed parse_tree.ParseTree) locally_nameless.Expr {
	return toLocallyNameless(parsed, stack.Empty[string]())
}

func toLocallyNameless(parsed parse_tree.ParseTree, bound stack.Stack[string]) locally_nameless.Expr {
	switch item := parsed.Item.(type) {
	case parse_tree.Parens:
		return toLocallyNameless(item.Child, bound)
	case parse_tree.Var:
		for index, boundName := range bound.IndexedItems() {
			if boundName == item.Name {
				return locally_nameless.NewBound(index)
			}
		}
		return locally_nameless.NewFree(item.Name)
	case parse_tree.Lambda:
		return locally_nameless.NewLambda(
			item.ArgName,
			toLocallyNameless(item.Body, bound.Push(item.ArgName)),
		)
	case parse_tree.App:
		app := locally_nameless.NewApp(
			toLocallyNameless(item.Callee, bound),
			toLocallyNameless(item.Args.First, bound),
		)
		for _, arg := range item.Args.More {
			app = locally_nameless.NewApp(app, toLocallyNameless(arg, bound))
		}
		return app
	default:
		panic("unknown parse tree")
	}
}
