package pretty

import (
	"fmt"
	"strings"

	ln "github.com/gusbicalho/go-lambda/locally_nameless/expr"
	"github.com/gusbicalho/go-lambda/pretty"
)

func ToPrettyDoc(expr ln.Expr) pretty.Doc {
	return ExprToPrettyDoc(expr, ln.EmptyContext())
}

func ExprToPrettyDoc(expr ln.Expr, ctx ln.DisplayContext) pretty.Doc {
	return ln.CaseExpr(expr, visitPretty{ctx})
}

type visitPretty struct{ ln.DisplayContext }

func (v visitPretty) CaseFree(expr ln.FreeVar) pretty.Doc {
	return pretty.FromString(expr.Name())
}
func (v visitPretty) CaseBound(expr ln.BoundVar) pretty.Doc {
	builder := strings.Builder{}
	if err := expr.WriteLambdaNotation(v.DisplayContext, &builder); err != nil {
		panic(err)
	}
	return pretty.FromString(builder.String())
}
func (v visitPretty) CaseLambda(expr ln.Lambda) pretty.Doc {
	ctx, argName := v.DisplayContext.BindFree(expr.ArgName())
	nameLength := uint(len(argName))
	return pretty.Sequence(
		pretty.FromString(fmt.Sprint("λ", argName, " ─┬─")),
		pretty.Indent(nameLength+1, pretty.PrefixLines([]string{"  │ "},
			ExprToPrettyDoc(expr.Body(), ctx),
		)),
		pretty.Indent(nameLength+1, pretty.FromString("  ╰─")),
	)
}
func (v visitPretty) CaseApp(expr ln.App) pretty.Doc {
	return pretty.Sequence(
		ExprToPrettyDoc(expr.Callee(), v.DisplayContext),
		pretty.PrefixLines([]string{
			"└► ",
			"   ",
		}, ExprToPrettyDoc(expr.Arg(), v.DisplayContext)),
	)
}
