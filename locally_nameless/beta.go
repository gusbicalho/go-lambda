package locally_nameless

import (
	"fmt"

	"github.com/gusbicalho/go-lambda/lazy"
)

func BetaReduce(lambda Lambda, arg Expr) Expr {
	return subst(lambda.body, lazy.Wrap(arg), 0)
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
