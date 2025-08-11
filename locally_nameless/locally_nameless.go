package locally_nameless

import (
	"github.com/gusbicalho/go-lambda/pretty"
	"io"
)

// Represents an AST with strings for free variables
// and de Bruijn indexes for lambda params

type Expr interface {
	pretty.Pretty[DisplayContext]
	writeAsNotation(ctx DisplayContext, writer io.StringWriter)
	sealed()
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
