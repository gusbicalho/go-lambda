package locally_nameless

import (
	"fmt"

	"github.com/gusbicalho/go-lambda/pretty"
	"github.com/gusbicalho/go-lambda/stack"
)

// Represents an AST with strings for free variables
// and de Bruijn indexes for lambda params

type Expr interface {
	pretty.Pretty[stack.Stack[string]]
	sealed()
}

type FreeVar struct{ name string }

func NewFree(name string) Expr { return FreeVar{name} }
func (v FreeVar) Name() string { return v.name }

func (FreeVar) sealed() {}
func (v FreeVar) ToPrettyDoc(_ stack.Stack[string]) pretty.PrettyDoc {
	return pretty.FromString(v.name)
}

type BoundVar struct{ index uint }

func NewBound(index uint) Expr { return BoundVar{index} }
func (v BoundVar) Index() uint { return v.index }

func (BoundVar) sealed() {}
func (v BoundVar) ToPrettyDoc(bound stack.Stack[string]) pretty.PrettyDoc {
	if name, found := bound.Nth(v.index, ""); found {
		return pretty.FromString(fmt.Sprint(v.index, " (", name, ")"))
	}
	return pretty.FromString(fmt.Sprint(v.index, " <unbound>"))
}

type Lambda struct {
	argName string
	body    Expr
}

func NewLambda(argName string, body Expr) Expr { return Lambda{argName, body} }
func (v Lambda) ArgName() string               { return v.argName }
func (v Lambda) Body() Expr                    { return v.body }

func (Lambda) sealed() {}
func (item Lambda) ToPrettyDoc(ctx stack.Stack[string]) pretty.PrettyDoc {
	return pretty.Sequence(
		pretty.FromString(fmt.Sprint("\\", item.argName, ".")),
		pretty.Indent(2, item.body.ToPrettyDoc(ctx.Push(item.argName))),
	)
}

type App struct {
	callee Expr
	arg    Expr
}

func NewApp(callee, arg Expr) Expr { return App{callee, arg} }
func (v App) Callee() Expr         { return v.callee }
func (v App) Arg() Expr            { return v.arg }

func (App) sealed() {}
func (item App) ToPrettyDoc(ctx stack.Stack[string]) pretty.PrettyDoc {
	return pretty.Sequence(
		item.callee.ToPrettyDoc(ctx),
		pretty.Indent(2, item.arg.ToPrettyDoc(ctx)),
	)
}
