package locally_nameless

import (
	"fmt"
	"iter"

	"github.com/gusbicalho/go-lambda/lazy"
)

func BetaReduce(lambda Lambda, arg Expr) Expr {
	return subst(lambda.body, lazy.Wrap(arg), 0)
}

type BetaReductionLocus struct {
	Hole   Hole
	Lambda Lambda
	Arg    Expr
}

func (locus BetaReductionLocus) Reduce() Expr {
	return locus.Hole.Fill(BetaReduce(locus.Lambda, locus.Arg))
}

type Hole struct {
	fill func(expr Expr) Expr
}

func (h Hole) Fill(expr Expr) Expr {
	return h.fill(expr)
}

func identityHole() Hole {
	return Hole{fill: func(expr Expr) Expr { return expr }}
}

func composeHoles(hole Hole, holes ...Hole) Hole {
	for _, h := range holes {
		fillOuter := hole.fill
		fillInner := h.fill
		hole = Hole{fill: func(expr Expr) Expr {
			return fillOuter(fillInner(expr))
		}}
	}
	return hole
}

func BetaReductionLocii(expr Expr) iter.Seq[BetaReductionLocus] {
	return func(yield func(BetaReductionLocus) bool) {
		betaReductionLocii(yield, identityHole(), expr)
	}
}

func subst(body Expr, arg lazy.Lazy[Expr], index uint) Expr {
	switch body := body.(type) {
	case FreeVar:
		return body
	case BoundVar:
		switch {
		case body.index == index:
			return arg.Get()
		case body.index > index:
			// References to bound vars above the one being bound
			// have their indexes decreased by one
			// because one lambda binder is being removed by the application
			return NewBound(body.index - 1)
		default:
			return body
		}
	case App:
		return NewApp(
			subst(body.callee, arg, index),
			subst(body.arg, arg, index),
		)
	case Lambda:
		shiftedArg := lazy.New(func() Expr { return shift(arg.Get(), 0) })
		return NewLambda(body.argName, subst(body.body, shiftedArg, index+1))
	default:
		panic(fmt.Sprint("Unknown Expr ", body))
	}
}

func shift(expr Expr, underBinders uint) Expr {
	switch expr := expr.(type) {
	case BoundVar:
		if expr.index < underBinders {
			// References to variables bound within the term being shifted
			// They stay the same
			return expr
		}
		// References to variables outside the term being shifted
		// They get increased
		return NewBound(expr.index + 1)
	case FreeVar:
		return expr
	case App:
		return NewApp(
			shift(expr.callee, underBinders),
			shift(expr.arg, underBinders),
		)
	case Lambda:
		return NewLambda(expr.argName, shift(expr.body, underBinders+1))
	default:
		panic(fmt.Sprint("Unknown Expr ", expr))
	}
}

func betaReductionLocii(yield func(BetaReductionLocus) bool, hole Hole, expr Expr) bool {
	switch expr := expr.(type) {
	case App:
		switch callee := expr.callee.(type) {
		case Lambda:
			if !yield(BetaReductionLocus{
				Hole:   hole,
				Lambda: callee,
				Arg:    expr.arg,
			}) {
				return false
			}
		}
		return betaReductionLocii(yield, composeHoles(hole, expr.calleeHole()), expr.callee) &&
			betaReductionLocii(yield, composeHoles(hole, expr.argHole()), expr.arg)
	case Lambda:
		return betaReductionLocii(
			yield,
			expr.hole(),
			expr.body,
		)
	}
	return true
}

func (expr Lambda) hole() Hole {
	argName := expr.argName
	return Hole{fill: func(expr Expr) Expr { return NewLambda(argName, expr) }}
}

func (expr App) calleeHole() Hole {
	arg := expr.arg
	return Hole{fill: func(e Expr) Expr { return NewApp(e, arg) }}
}

func (expr App) argHole() Hole {
	callee := expr.callee
	return Hole{fill: func(e Expr) Expr { return NewApp(callee, e) }}
}
