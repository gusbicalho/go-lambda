package locally_nameless

import (
	"fmt"
	"io"
	"strings"
)

func ToLambdaNotation(expr Expr, displayBoundVarAs DisplayBoundVarAs) string {
	builder := strings.Builder{}
	err := expr.writeLambdaNotation(
		EmptyContext().WithDisplayBoundVarAs(displayBoundVarAs),
		&builder,
	)
	if err != nil {
		panic(err)
	}
	return builder.String()
}

type asLambdaNotation interface {
	writeLambdaNotation(ctx DisplayContext, writer io.StringWriter) error
}

func writeStrings(writer io.StringWriter, strings ...string) error {
	for _, str := range strings {
		_, err := writer.WriteString(str)
		if err != nil {
			return err
		}
	}
	return nil
}

func (expr FreeVar) writeLambdaNotation(_ DisplayContext, writer io.StringWriter) error {
	return writeStrings(writer, expr.name)
}

func (expr BoundVar) writeLambdaNotation(ctx DisplayContext, writer io.StringWriter) error {
	switch ctx.displayBoundVarAs {
	case DisplayIndex:
		return writeStrings(writer, fmt.Sprint(expr.index))
	case DisplayName:
		if name, found := ctx.bound.Nth(expr.index, ""); found {
			return writeStrings(writer, name)
		}
	case DisplayBoth:
		if name, found := ctx.bound.Nth(expr.index, ""); found {
			return writeStrings(writer, fmt.Sprint(expr.index), ":", name)
		}
	}
	return writeStrings(writer, fmt.Sprint(expr.index), ":<outofscope>")
}

func (expr Lambda) writeLambdaNotation(ctx DisplayContext, writer io.StringWriter) error {
	ctx, argName := ctx.bindFree(expr.argName)
	if err := writeStrings(writer, "\\", argName, ". "); err != nil {
		return err
	}

	return expr.body.writeLambdaNotation(ctx, writer)
}

func (expr App) writeLambdaNotation(ctx DisplayContext, writer io.StringWriter) error {
	switch callee := expr.callee.(type) {
	case Lambda:
		if err := writeStrings(writer, "("); err != nil {
			return err
		}
		if err := callee.writeLambdaNotation(ctx, writer); err != nil {
			return err
		}
		if err := writeStrings(writer, ")"); err != nil {
			return err
		}
	default:
		if err := callee.writeLambdaNotation(ctx, writer); err != nil {
			return err
		}
	}

	if err := writeStrings(writer, " "); err != nil {
		return err
	}

	switch arg := expr.arg.(type) {
	case App, Lambda:
		if err := writeStrings(writer, "("); err != nil {
			return err
		}
		if err := arg.writeLambdaNotation(ctx, writer); err != nil {
			return err
		}
		if err := writeStrings(writer, ")"); err != nil {
			return err
		}
	default:
		if err := arg.writeLambdaNotation(ctx, writer); err != nil {
			return err
		}
	}
	return nil
}
