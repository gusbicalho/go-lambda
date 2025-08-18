package beta_reduce

import (
	"iter"

	"github.com/gusbicalho/go-lambda/lazy"
	"github.com/gusbicalho/go-lambda/locally_nameless/expr"
	"github.com/gusbicalho/go-lambda/locally_nameless/hole"
	ln_pretty "github.com/gusbicalho/go-lambda/locally_nameless/pretty"
	"github.com/gusbicalho/go-lambda/locally_nameless/walk"
	"github.com/gusbicalho/go-lambda/pretty"
)

func BetaReduce(lambda expr.Lambda, arg expr.Expr) expr.Expr {
	return subst(lambda.Body(), lazy.Wrap(arg), 0)
}

type BetaRedex struct {
	Hole   hole.Hole
	Lambda expr.Lambda
	Arg    expr.Expr
}

func (redex BetaRedex) ToPrettyDoc(_ any) pretty.Doc {
	return redex.Hole.ToPrettyDoc(func(ctx expr.DisplayContext) pretty.Doc {
		return pretty.ForegroundColor(pretty.ColorYellow,
			ln_pretty.ExprToPrettyDoc(expr.NewApp(redex.Lambda, redex.Arg), ctx),
		)
	})
}

func AsBetaRedex(e expr.Expr) *BetaRedex {
	if app, ok := e.(expr.App); ok {
		if callee, ok := app.Callee().(expr.Lambda); ok {
			return &BetaRedex{
				Hole:   hole.IdentityHole(),
				Lambda: callee,
				Arg:    app.Arg(),
			}
		}
	}
	return nil
}

func (locus BetaRedex) Reduce() expr.Expr {
	return locus.Hole.Fill(BetaReduce(locus.Lambda, locus.Arg))
}

func BetaRedexes(e expr.Expr) iter.Seq[BetaRedex] {
	return func(yield func(BetaRedex) bool) {
		for h, e := range walk.Pre(e) {
			if redex := AsBetaRedex(e); redex != nil {
				redex.Hole = hole.ComposeHoles(h, redex.Hole)
				if !yield(*redex) {
					return
				}
			}
		}
	}
}

func subst(body expr.Expr, arg lazy.Lazy[expr.Expr], index uint) expr.Expr {
	return expr.CaseExpr(body, substVisit{index, arg})
}

type substVisit struct {
	index uint
	arg   lazy.Lazy[expr.Expr]
}

func (v substVisit) CaseFree(body expr.FreeVar) expr.Expr {
	return body
}
func (v substVisit) CaseBound(body expr.BoundVar) expr.Expr {
	switch {
	case body.Index() == v.index:
		return v.arg.Get()
	case body.Index() > v.index:
		// References to bound vars above the one being bound
		// have their indexes decreased by one
		// because one lambda binder is being removed by the application
		return expr.NewBound(body.Index() - 1)
	default:
		return body
	}
}
func (v substVisit) CaseApp(body expr.App) expr.Expr {
	return expr.NewApp(
		subst(body.Callee(), v.arg, v.index),
		subst(body.Arg(), v.arg, v.index),
	)
}
func (v substVisit) CaseLambda(body expr.Lambda) expr.Expr {
	shiftedArg := lazy.New(func() expr.Expr { return shift(v.arg.Get(), 0) })
	return expr.NewLambda(body.ArgName(), subst(body.Body(), shiftedArg, v.index+1))
}

func shift(e expr.Expr, underBinders uint) expr.Expr {
	return expr.CaseExpr(e, shiftVisit{underBinders})
}

type shiftVisit struct {
	underBinders uint
}

func (v shiftVisit) CaseBound(e expr.BoundVar) expr.Expr {
	if e.Index() < v.underBinders {
		// References to variables bound within the term being shifted
		// They stay the same
		return e
	}
	// References to variables outside the term being shifted
	// They get increased
	return expr.NewBound(e.Index() + 1)
}
func (v shiftVisit) CaseFree(e expr.FreeVar) expr.Expr {
	return e
}
func (v shiftVisit) CaseApp(e expr.App) expr.Expr {
	return expr.NewApp(
		shift(e.Callee(), v.underBinders),
		shift(e.Arg(), v.underBinders),
	)
}
func (v shiftVisit) CaseLambda(e expr.Lambda) expr.Expr {
	return expr.NewLambda(e.ArgName(), shift(e.Body(), v.underBinders+1))
}
