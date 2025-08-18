package locally_nameless

import (
	"fmt"
	"slices"

	"github.com/gusbicalho/go-lambda/pretty"
)

type Hole struct {
	holeImpl
}

func (h Hole) ToPrettyDoc(fill func(DisplayContext) pretty.Doc) pretty.Doc {
	return h.holeImpl.toPrettyDoc(EmptyContext(), fill)
}

type holeImpl interface {
	Fill(expr Expr) Expr
	toPrettyDoc(ctx DisplayContext, fill func(DisplayContext) pretty.Doc) pretty.Doc
}

// Identity

func IdentityHole() Hole {
	return Hole{identityHole{}}
}

type identityHole struct{}

func (h identityHole) Fill(expr Expr) Expr {
	return expr
}

func (h identityHole) toPrettyDoc(ctx DisplayContext, fill func(DisplayContext) pretty.Doc) pretty.Doc {
	return fill(ctx)
}

// Compose

func ComposeHoles(holes ...Hole) Hole {
	impls := make([]holeImpl, len(holes))
	for i, hole := range holes {
		impls[i] = hole.holeImpl
	}
	return Hole{composeHoles{holes: impls}}
}

type composeHoles struct {
	holes []holeImpl
}

func (h composeHoles) Fill(expr Expr) Expr {
	for _, hole := range slices.Backward(h.holes) {
		expr = hole.Fill(expr)
	}
	return expr
}

func (h composeHoles) toPrettyDoc(ctx DisplayContext, fill func(DisplayContext) pretty.Doc) pretty.Doc {
	return composeToPrettyDoc(h.holes, ctx, fill)
}

func composeToPrettyDoc(holes []holeImpl, ctx DisplayContext, fill func(DisplayContext) pretty.Doc) pretty.Doc {
	switch len(holes) {
	case 0:
		return fill(ctx)
	case 1:
		return holes[0].toPrettyDoc(ctx, fill)
	}
	hole := holes[0]
	more := holes[1:]
	fillMore := func(ctx DisplayContext) pretty.Doc {
		return composeToPrettyDoc(more, ctx, fill)
	}
	return hole.toPrettyDoc(ctx, fillMore)
}

// Lambda

func (expr Lambda) Hole() Hole {
	return Hole{lambdaBodyHole{argName: expr.argName}}
}

type lambdaBodyHole struct {
	argName string
}

func (h lambdaBodyHole) Fill(expr Expr) Expr {
	return NewLambda(h.argName, expr)
}

func (h lambdaBodyHole) toPrettyDoc(ctx DisplayContext, fill func(DisplayContext) pretty.Doc) pretty.Doc {
	ctx, argName := ctx.bindFree(h.argName)
	nameLength := uint(len(argName))
	return pretty.Sequence(
		pretty.FromString(fmt.Sprint("λ", argName, " ─┬─")),
		pretty.Indent(nameLength+1, pretty.PrefixLines([]string{"  │ "},
			fill(ctx),
		)),
		pretty.Indent(nameLength+1, pretty.FromString("  ╰─")),
	)
}

// App: Callee
func (expr App) CalleeHole() Hole {
	return Hole{appCalleeHole{arg: expr.arg}}
}

type appCalleeHole struct {
	arg Expr
}

func (h appCalleeHole) Fill(expr Expr) Expr {
	return NewApp(expr, h.arg)
}

func (h appCalleeHole) toPrettyDoc(ctx DisplayContext, fill func(DisplayContext) pretty.Doc) pretty.Doc {
	return pretty.Sequence(
		fill(ctx),
		pretty.PrefixLines([]string{
			"└► ",
			"   ",
		}, h.arg.ToPrettyDoc(ctx)),
	)
}

// App: Arg
func (expr App) ArgHole() Hole {
	return Hole{appArgHole{callee: expr.callee}}
}

type appArgHole struct {
	callee Expr
}

func (h appArgHole) Fill(expr Expr) Expr {
	return NewApp(h.callee, expr)
}

func (h appArgHole) toPrettyDoc(ctx DisplayContext, fill func(DisplayContext) pretty.Doc) pretty.Doc {
	return pretty.Sequence(
		h.callee.ToPrettyDoc(ctx),
		pretty.PrefixLines([]string{
			"└► ",
			"   ",
		}, fill(ctx)),
	)
}
