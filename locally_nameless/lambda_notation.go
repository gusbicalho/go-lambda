package locally_nameless

import (
	"fmt"
	"io"
	"strings"
)

func ToLambdaNotation(expr Expr, displayBoundVarAs DisplayBoundVarAs) string {
	builder := strings.Builder{}
	expr.writeAsNotation(
		EmptyContext().WithDisplayBoundVarAs(displayBoundVarAs),
		&builder,
	)
	return builder.String()
}
func (expr FreeVar) writeAsNotation(_ DisplayContext, writer io.StringWriter) {
	writer.WriteString(expr.name)
}

func (expr BoundVar) writeAsNotation(ctx DisplayContext, writer io.StringWriter) {
	switch ctx.displayBoundVarAs {
	case DisplayIndex:
		writer.WriteString(fmt.Sprint(expr.index))
		return
	case DisplayName:
		if name, found := ctx.bound.Nth(expr.index, ""); found {
			writer.WriteString(name)
			return
		}
	case DisplayBoth:
		if name, found := ctx.bound.Nth(expr.index, ""); found {
			writer.WriteString(fmt.Sprint(expr.index))
			writer.WriteString(":")
			writer.WriteString(name)
			return
		}
	}
	writer.WriteString(fmt.Sprint(expr.index))
	writer.WriteString(":<outofscope>")
}

func (expr Lambda) writeAsNotation(ctx DisplayContext, writer io.StringWriter) {
	writer.WriteString("\\")

	ctx, argName := ctx.bindFree(expr.argName)

	writer.WriteString(argName)
	writer.WriteString(". ")
	expr.body.writeAsNotation(ctx, writer)
}

func (expr App) writeAsNotation(ctx DisplayContext, writer io.StringWriter) {
	switch callee := expr.callee.(type) {
	case Lambda:
		writer.WriteString("(")
		callee.writeAsNotation(ctx, writer)
		writer.WriteString(")")
	default:
		callee.writeAsNotation(ctx, writer)
	}
	writer.WriteString(" ")
	switch arg := expr.arg.(type) {
	case App, Lambda:
		writer.WriteString("(")
		arg.writeAsNotation(ctx, writer)
		writer.WriteString(")")
	default:
		arg.writeAsNotation(ctx, writer)
	}
}
