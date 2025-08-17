package locally_nameless

import (
	"fmt"
	"slices"

	"github.com/gusbicalho/go-lambda/pretty"
)

type Hole interface {
	Fill(expr Expr) Expr
	pretty.Pretty[displayHoleContext]
}

func HoleToPrettyDoc(hole Hole, fill func(DisplayContext) pretty.Doc) pretty.Doc {
	return hole.ToPrettyDoc(displayHoleContext{EmptyContext(), fill})
}

func HoleToPrettyString(hole Hole, fill func(DisplayContext) pretty.Doc) string {
	return HoleToPrettyDoc(hole, fill).String()
}

type displayHoleContext struct {
	context DisplayContext
	fill    func(DisplayContext) pretty.Doc
}

// Identity

func IdentityHole() Hole {
	return identityHole{}
}

type identityHole struct{}

func (h identityHole) Fill(expr Expr) Expr {
	return expr
}

func (h identityHole) ToPrettyDoc(displayCtx displayHoleContext) pretty.Doc {
	return displayCtx.fill(displayCtx.context)
}

// Compose

func ComposeHoles(holes ...Hole) Hole {
	return composeHoles{holes: holes}
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

func (h composeHoles) ToPrettyDoc(displayHoleCtx displayHoleContext) pretty.Doc {
	return composeToPrettyDoc(h.holes, displayHoleCtx.context, displayHoleCtx.fill)
}

func composeToPrettyDoc(holes []Hole, ctx DisplayContext, fill func(DisplayContext) pretty.Doc) pretty.Doc {
	switch len(holes) {
	case 0:
		return fill(ctx)
	case 1:
		return holes[0].ToPrettyDoc(displayHoleContext{ctx, fill})
	}
	hole := holes[0]
	more := holes[1:]
	fillMore := func(ctx DisplayContext) pretty.Doc {
		return composeToPrettyDoc(more, ctx, fill)
	}
	return hole.ToPrettyDoc(displayHoleContext{ctx, fillMore})
}

// Lambda

func (expr Lambda) Hole() Hole {
	return lambdaBodyHole{argName: expr.argName}
}

type lambdaBodyHole struct {
	argName string
}

func (h lambdaBodyHole) Fill(expr Expr) Expr {
	return NewLambda(h.argName, expr)
}

func (h lambdaBodyHole) ToPrettyDoc(displayHoleCtx displayHoleContext) pretty.Doc {
	ctx := displayHoleCtx.context
	ctx, argName := ctx.bindFree(h.argName)
	nameLength := uint(len(argName))
	return pretty.Sequence(
		pretty.FromString(fmt.Sprint("λ", argName, " ─┬─")),
		pretty.Indent(nameLength+1, pretty.PrefixLines([]string{"  │ "},
			displayHoleCtx.fill(ctx),
		)),
		pretty.Indent(nameLength+1, pretty.FromString("  ╰─")),
	)
}

// App: Callee
func (expr App) CalleeHole() Hole {
	return appCalleeHole{arg: expr.arg}
}

type appCalleeHole struct {
	arg Expr
}

func (h appCalleeHole) Fill(expr Expr) Expr {
	return NewApp(expr, h.arg)
}

func (h appCalleeHole) ToPrettyDoc(displayHoleCtx displayHoleContext) pretty.Doc {
	ctx := displayHoleCtx.context
	return pretty.Sequence(
		displayHoleCtx.fill(ctx),
		pretty.PrefixLines([]string{
			"└► ",
			"   ",
		}, h.arg.ToPrettyDoc(ctx)),
	)
}

// App: Arg
func (expr App) ArgHole() Hole {
	return appArgHole{callee: expr.callee}
}

type appArgHole struct {
	callee Expr
}

func (h appArgHole) Fill(expr Expr) Expr {
	return NewApp(h.callee, expr)
}

func (h appArgHole) ToPrettyDoc(displayHoleCtx displayHoleContext) pretty.Doc {
	ctx := displayHoleCtx.context
	return pretty.Sequence(
		h.callee.ToPrettyDoc(ctx),
		pretty.PrefixLines([]string{
			"└► ",
			"   ",
		}, displayHoleCtx.fill(ctx)),
	)
}
