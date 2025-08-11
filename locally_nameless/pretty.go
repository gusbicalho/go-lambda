package locally_nameless

import (
	"fmt"
	"strings"

	"github.com/gusbicalho/go-lambda/pretty"
)

func ToPrettyString(expr Expr) string {
	return expr.ToPrettyDoc(EmptyContext()).String()
}

func (expr FreeVar) ToPrettyDoc(_ DisplayContext) pretty.Doc {
	return pretty.FromString(expr.name)
}

func (expr BoundVar) ToPrettyDoc(ctx DisplayContext) pretty.Doc {
	builder := strings.Builder{}
	if err := expr.writeLambdaNotation(ctx, &builder); err != nil {
		panic(err)
	}
	return pretty.FromString(builder.String())
}

func (expr Lambda) ToPrettyDoc(ctx DisplayContext) pretty.Doc {
	ctx, argName := ctx.bindFree(expr.argName)
	nameLength := uint(len(argName))
	return pretty.Sequence(
		pretty.FromString(fmt.Sprint("λ", argName, " ─┬─")),
		pretty.Indent(nameLength+1, pretty.PrefixLines([]string{"  │ "},
			expr.body.ToPrettyDoc(ctx),
		)),
		pretty.Indent(nameLength+1, pretty.FromString("  ╰─")),
	)
}

func (expr App) ToPrettyDoc(ctx DisplayContext) pretty.Doc {
	return pretty.Sequence(
		expr.callee.ToPrettyDoc(ctx),
		pretty.PrefixLines([]string{
			"└► ",
			"   ",
		}, expr.arg.ToPrettyDoc(ctx)),
	)
}
