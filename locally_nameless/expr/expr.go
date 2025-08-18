package expr

import (
	"fmt"
)

// Represents an AST with strings for free variables
// and de Bruijn indexes for lambda params

type Expr interface {
	asLambdaNotation
	sealed()
}

type VisitExpr[R any] interface {
	CaseFree(FreeVar) R
	CaseBound(BoundVar) R
	CaseLambda(Lambda) R
	CaseApp(App) R
}

func CaseExpr[R any](expr Expr, visit VisitExpr[R]) R {
	switch expr := expr.(type) {
	case FreeVar:
		return visit.CaseFree(expr)
	case BoundVar:
		return visit.CaseBound(expr)
	case Lambda:
		return visit.CaseLambda(expr)
	case App:
		return visit.CaseApp(expr)
	default:
		panic(fmt.Sprint("Unknown expr ", expr))
	}
}

type FreeVar struct{ name string }

func NewFree(name string) FreeVar { return FreeVar{name} }
func (expr FreeVar) Name() string { return expr.name }

func (FreeVar) sealed() {}

type BoundVar struct{ index uint }

func NewBound(index uint) BoundVar { return BoundVar{index} }
func (expr BoundVar) Index() uint  { return expr.index }

func (BoundVar) sealed() {}

type Lambda struct {
	argName string
	body    Expr
}

func NewLambda(argName string, body Expr) Lambda { return Lambda{argName, body} }
func (expr Lambda) ArgName() string              { return expr.argName }
func (expr Lambda) Body() Expr                   { return expr.body }

func (Lambda) sealed() {}

type App struct {
	callee Expr
	arg    Expr
}

func NewApp(callee, arg Expr) App { return App{callee, arg} }
func (expr App) Callee() Expr     { return expr.callee }
func (expr App) Arg() Expr        { return expr.arg }

func (App) sealed() {}
