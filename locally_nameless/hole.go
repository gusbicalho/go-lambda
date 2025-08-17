package locally_nameless

import (
	"slices"
)

type Hole interface {
	Fill(expr Expr) Expr
}

// Basic

type identityHole struct{}

func (h identityHole) Fill(expr Expr) Expr {
	return expr
}

func IdentityHole() Hole {
	return identityHole{}
}

type composeHoles struct {
	holes []Hole
}

func (h composeHoles) Fill(expr Expr) Expr {
	for _, hole := range slices.Backward(h.holes) {
		expr = hole.Fill(expr)
	}
	return expr
}

func ComposeHoles(holes ...Hole) Hole {
	return composeHoles{holes: holes}
}

// Lambda

type lambdaBodyHole struct {
	argName string
}

func (h lambdaBodyHole) Fill(expr Expr) Expr {
	return NewLambda(h.argName, expr)
}

func (expr Lambda) Hole() Hole {
	return lambdaBodyHole{argName: expr.argName}
}

// App: Callee
type appCalleeHole struct {
	arg Expr
}

func (h appCalleeHole) Fill(expr Expr) Expr {
	return NewApp(expr, h.arg)
}

func (expr App) CalleeHole() Hole {
	return appCalleeHole{arg: expr.arg}
}

// App: Arg
type appArgHole struct {
	callee Expr
}

func (h appArgHole) Fill(expr Expr) Expr {
	return NewApp(h.callee, expr)
}

func (expr App) ArgHole() Hole {
	return appArgHole{callee: expr.callee}
}
