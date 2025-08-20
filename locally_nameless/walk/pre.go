package walk

import (
	"github.com/gusbicalho/go-lambda/locally_nameless/expr"
	"github.com/gusbicalho/go-lambda/locally_nameless/hole"
)

func Pre(e expr.Expr) Walk {
	return preWalk(hole.IdentityHole(), e)
}

func preWalk(h hole.Hole, e expr.Expr) Walk {
	return expr.CaseExpr(e, preWalkVisit{h})
}

type preWalkVisit struct{ hole hole.Hole }

func (v preWalkVisit) CaseBound(e expr.BoundVar) Walk {
	return preWalkFirst[expr.BoundVar]{v.hole, e, nil}
}
func (v preWalkVisit) CaseFree(e expr.FreeVar) Walk {
	return preWalkFirst[expr.FreeVar]{v.hole, e, nil}
}
func (v preWalkVisit) CaseLambda(e expr.Lambda) Walk {
	return preWalkFirst[expr.Lambda]{
		v.hole, e,
		func(e expr.Lambda) Walk {
			return lambdaBodyPreWalk{
				parent:  v.hole,
				argName: e.ArgName(),
				body:    preWalk(hole.IdentityHole(), e.Body()),
			}
		},
	}
}

func (v preWalkVisit) CaseApp(e expr.App) Walk {
	return preWalkFirst[expr.App]{
		v.hole, e,
		func(e expr.App) Walk {
			return appCalleePreWalk{
				parent: v.hole,
				arg:    e.Arg(),
				callee: preWalk(hole.IdentityHole(), e.Callee()),
			}
		},
	}
}

type preWalkFirst[E expr.Expr] struct {
	hole hole.Hole
	expr E
	next func(e E) Walk
}

func (w preWalkFirst[E]) Focus() Focus { return Focus{w.hole, w.expr} }
func (w preWalkFirst[E]) UpdateExpr(update func(expr.Expr) expr.Expr) Walk {
	return preWalk(w.hole, update(w.expr))
}
func (w preWalkFirst[E]) Prev() Walk { return nil }
func (w preWalkFirst[E]) Next() Walk {
	if w.next == nil {
		return nil
	}
	return w.next(w.expr)
}

type lambdaBodyPreWalk struct {
	parent  hole.Hole
	argName string
	body    Walk
}

func (w lambdaBodyPreWalk) Focus() Focus {
	f := w.body.Focus()
	bodyHole := hole.BodyHole(expr.NewLambda(w.argName, expr.NewFree("irrelevant")))
	return Focus{hole.ComposeHoles(w.parent, bodyHole, f.Hole), f.Expr}
}

func (w lambdaBodyPreWalk) UpdateExpr(update func(expr.Expr) expr.Expr) Walk {
	w.body = w.body.UpdateExpr(update)
	return w
}
func (w lambdaBodyPreWalk) Prev() Walk {
	if prev := w.body.Prev(); prev != nil {
		w.body = prev
		return w
	}
	return preWalk(w.parent, expr.NewLambda(w.argName, w.body.Focus().Realize()))
}
func (w lambdaBodyPreWalk) Next() Walk {
	if next := w.body.Next(); next != nil {
		w.body = next
		return w
	}
	return nil
}

type appCalleePreWalk struct {
	parent hole.Hole
	arg    expr.Expr
	callee Walk
}

func (w appCalleePreWalk) Focus() Focus {
	focus := w.callee.Focus()
	calleeHole := hole.CalleeHole(expr.NewApp(expr.NewFree("irrelevant"), w.arg))
	return Focus{hole.ComposeHoles(w.parent, calleeHole, focus.Hole), focus.Expr}
}
func (w appCalleePreWalk) UpdateExpr(update func(expr.Expr) expr.Expr) Walk {
	w.callee = w.callee.UpdateExpr(update)
	return w
}
func (w appCalleePreWalk) Prev() Walk {
	if prev := w.callee.Prev(); prev != nil {
		w.callee = prev
		return w
	}
	return preWalk(w.parent, expr.NewApp(w.callee.Focus().Realize(), w.arg))
}
func (w appCalleePreWalk) Next() Walk {
	if next := w.callee.Next(); next != nil {
		w.callee = next
		return w
	}
	return appArgPreWalk{
		parent:  w.parent,
		callee:  w.callee,
		argWalk: preWalk(hole.IdentityHole(), w.arg),
	}
}

type appArgPreWalk struct {
	parent  hole.Hole
	callee  Walk
	argWalk Walk
}

func (w appArgPreWalk) Focus() Focus {
	callee := w.callee.Focus().Realize()
	argHole := hole.ArgHole(expr.NewApp(callee, expr.NewFree("irrelevant")))
	focus := w.argWalk.Focus()
	return Focus{hole.ComposeHoles(w.parent, argHole, focus.Hole), focus.Expr}
}
func (w appArgPreWalk) UpdateExpr(update func(expr.Expr) expr.Expr) Walk {
	w.argWalk = w.argWalk.UpdateExpr(update)
	return w
}
func (w appArgPreWalk) Prev() Walk {
	if prev := w.argWalk.Prev(); prev != nil {
		w.argWalk = prev
		return w
	}
	arg := w.argWalk.Focus().Realize()
	return appCalleePreWalk{
		parent: w.parent,
		arg:    arg,
		callee: w.callee,
	}
}
func (w appArgPreWalk) Next() Walk {
	if next := w.argWalk.Next(); next != nil {
		w.argWalk = next
		return w
	}
	return nil
}
