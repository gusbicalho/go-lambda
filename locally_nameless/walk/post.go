package walk

import (
	"iter"

	"github.com/gusbicalho/go-lambda/locally_nameless/expr"
	"github.com/gusbicalho/go-lambda/locally_nameless/hole"
)

func Post(e expr.Expr) iter.Seq2[hole.Hole, expr.Expr] {
	return func(yield func(hole.Hole, expr.Expr) bool) {
		postwalk(e, hole.IdentityHole(), yield)
	}
}

func postwalk(e expr.Expr, h hole.Hole, yield func(hole.Hole, expr.Expr) bool) bool {
	return expr.CaseExpr(e, postVisit{h, yield})
}

type postVisit struct {
	hole  hole.Hole
	yield func(hole.Hole, expr.Expr) bool
}

func (v postVisit) CaseBound(e expr.BoundVar) bool {
	return v.yield(v.hole, e)
}

func (v postVisit) CaseFree(e expr.FreeVar) bool {
	return v.yield(v.hole, e)
}

func (v postVisit) CaseLambda(e expr.Lambda) bool {
	return postwalk(e.Body(), hole.ComposeHoles(v.hole, hole.BodyHole(e)), v.yield) &&
		v.yield(v.hole, e)
}

func (v postVisit) CaseApp(e expr.App) bool {
	return postwalk(e.Callee(), hole.ComposeHoles(v.hole, hole.CalleeHole(e)), v.yield) &&
		postwalk(e.Arg(), hole.ComposeHoles(v.hole, hole.ArgHole(e)), v.yield) &&
		v.yield(v.hole, e)
}
