package hole

import (
	"fmt"
	"slices"

	ln "github.com/gusbicalho/go-lambda/locally_nameless/expr"
	ln_pretty "github.com/gusbicalho/go-lambda/locally_nameless/pretty"
	"github.com/gusbicalho/go-lambda/pretty"
)

type Hole struct {
	holeImpl
}

func (h Hole) ToPrettyDoc(fill func(ln.DisplayContext) pretty.Doc) pretty.Doc {
	return h.holeImpl.toPrettyDoc(ln.EmptyContext(), fill)
}

type holeImpl interface {
	Fill(expr ln.Expr) ln.Expr
	toPrettyDoc(ctx ln.DisplayContext, fill func(ln.DisplayContext) pretty.Doc) pretty.Doc
}

// Identity

func IdentityHole() Hole {
	return Hole{identityHole{}}
}

type identityHole struct{}

func (h identityHole) Fill(expr ln.Expr) ln.Expr {
	return expr
}

func (h identityHole) toPrettyDoc(ctx ln.DisplayContext, fill func(ln.DisplayContext) pretty.Doc) pretty.Doc {
	return fill(ctx)
}

// Compose

func ComposeHoles(holes ...Hole) Hole {
	impls := make([]holeImpl, 0, len(holes))
	for _, hole := range holes {
		switch hole := hole.holeImpl.(type) {
		case identityHole:
			// Noop - identity holes can be skipped in composition
		case composeHoles:
			// We can flatten nested composeHoles
			impls = append(impls, hole.holes...)
		default:
			impls = append(impls, hole)
		}
	}
	return Hole{composeHoles{holes: impls}}
}

type composeHoles struct {
	holes []holeImpl
}

func (h composeHoles) Fill(expr ln.Expr) ln.Expr {
	for _, hole := range slices.Backward(h.holes) {
		expr = hole.Fill(expr)
	}
	return expr
}

func (h composeHoles) toPrettyDoc(ctx ln.DisplayContext, fill func(ln.DisplayContext) pretty.Doc) pretty.Doc {
	return composeToPrettyDoc(h.holes, ctx, fill)
}

func composeToPrettyDoc(holes []holeImpl, ctx ln.DisplayContext, fill func(ln.DisplayContext) pretty.Doc) pretty.Doc {
	switch len(holes) {
	case 0:
		return fill(ctx)
	case 1:
		return holes[0].toPrettyDoc(ctx, fill)
	}
	hole := holes[0]
	more := holes[1:]
	fillMore := func(ctx ln.DisplayContext) pretty.Doc {
		return composeToPrettyDoc(more, ctx, fill)
	}
	return hole.toPrettyDoc(ctx, fillMore)
}

// Lambda

func BodyHole(expr ln.Lambda) Hole {
	return Hole{lambdaBodyHole{argName: expr.ArgName()}}
}

type lambdaBodyHole struct {
	argName string
}

func (h lambdaBodyHole) Fill(expr ln.Expr) ln.Expr {
	return ln.NewLambda(h.argName, expr)
}

func (h lambdaBodyHole) toPrettyDoc(ctx ln.DisplayContext, fill func(ln.DisplayContext) pretty.Doc) pretty.Doc {
	ctx, argName := ctx.BindFree(h.argName)
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
func CalleeHole(expr ln.App) Hole {
	return Hole{appCalleeHole{arg: expr.Arg()}}
}

type appCalleeHole struct {
	arg ln.Expr
}

func (h appCalleeHole) Fill(expr ln.Expr) ln.Expr {
	return ln.NewApp(expr, h.arg)
}

func (h appCalleeHole) toPrettyDoc(ctx ln.DisplayContext, fill func(ln.DisplayContext) pretty.Doc) pretty.Doc {
	return pretty.Sequence(
		fill(ctx),
		pretty.PrefixLines([]string{
			"└► ",
			"   ",
		}, ln_pretty.ExprToPrettyDoc(h.arg, ctx)),
	)
}

// App: Arg
func ArgHole(expr ln.App) Hole {
	return Hole{appArgHole{callee: expr.Callee()}}
}

type appArgHole struct {
	callee ln.Expr
}

func (h appArgHole) Fill(expr ln.Expr) ln.Expr {
	return ln.NewApp(h.callee, expr)
}

func (h appArgHole) toPrettyDoc(ctx ln.DisplayContext, fill func(ln.DisplayContext) pretty.Doc) pretty.Doc {
	return pretty.Sequence(
		ln_pretty.ExprToPrettyDoc(h.callee, ctx),
		pretty.PrefixLines([]string{
			"└► ",
			"   ",
		}, fill(ctx)),
	)
}
