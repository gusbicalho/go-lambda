package walk

import (
	"iter"

	"github.com/gusbicalho/go-lambda/locally_nameless/expr"
	"github.com/gusbicalho/go-lambda/locally_nameless/hole"
)

func Pre(e expr.Expr) iter.Seq2[hole.Hole, expr.Expr] {
	return func(yield func(hole.Hole, expr.Expr) bool) {
		prewalk(e, hole.IdentityHole(), yield)
	}
}

func prewalk(e expr.Expr, h hole.Hole, yield func(hole.Hole, expr.Expr) bool) bool {
	return expr.CaseExpr(e, preVisit{h, yield})
}

type preVisit struct {
	hole  hole.Hole
	yield func(hole.Hole, expr.Expr) bool
}

func (v preVisit) CaseBound(e expr.BoundVar) bool {
	return v.yield(v.hole, e)
}

func (v preVisit) CaseFree(e expr.FreeVar) bool {
	return v.yield(v.hole, e)
}

func (v preVisit) CaseLambda(e expr.Lambda) bool {
	return v.yield(v.hole, e) &&
		prewalk(
			e.Body(),
			hole.ComposeHoles(v.hole, hole.BodyHole(e)),
			v.yield,
		)
}

func (v preVisit) CaseApp(e expr.App) bool {
	return v.yield(v.hole, e) &&
		prewalk(e.Callee(), hole.ComposeHoles(v.hole, hole.CalleeHole(e)), v.yield) &&
		prewalk(e.Arg(), hole.ComposeHoles(v.hole, hole.ArgHole(e)), v.yield)
}
