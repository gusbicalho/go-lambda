package walk

import (
	"iter"

	"github.com/gusbicalho/go-lambda/locally_nameless/expr"
	"github.com/gusbicalho/go-lambda/locally_nameless/hole"
	lnpretty "github.com/gusbicalho/go-lambda/locally_nameless/pretty"
	"github.com/gusbicalho/go-lambda/pretty"
)

type Walk interface {
	Focus() Focus
	UpdateExpr(func(expr.Expr) expr.Expr) Walk
	Next() Walk
	Prev() Walk
}

type Focus struct {
	hole.Hole
	expr.Expr
}

func (f Focus) Realize() expr.Expr {
	return f.Hole.Fill(f.Expr)
}

func ToSeq(nav Walk) iter.Seq2[hole.Hole, expr.Expr] {
	return func(yield func(hole.Hole, expr.Expr) bool) {
		for n := nav; n != nil; n = n.Next() {
			focus := n.Focus()
			if !yield(focus.Hole, focus.Expr) {
				return
			}
		}
	}
}

func (focus Focus) ToPrettyDoc() pretty.Doc {
	return focus.Hole.ToPrettyDoc(
		func(ctx expr.DisplayContext) pretty.Doc {
			return pretty.TViewInvert(lnpretty.ExprToPrettyDoc(focus.Expr, ctx))
		},
	)
}

func ToPrettyDoc(nav Walk) pretty.Doc {
	return nav.Focus().ToPrettyDoc()
}
