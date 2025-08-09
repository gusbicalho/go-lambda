package parse_tree

import (
	"fmt"

	"github.com/gusbicalho/go-lambda/position"
	"github.com/gusbicalho/go-lambda/pretty"
)

type ParseTree struct {
	InputLocation position.Position
	Item          ParseItem
}

type ParseItem interface {
	pretty.Pretty[any]
	sealed()
}

type Parens struct {
	Child ParseTree
}

func (v Parens) sealed() {}

func (item Parens) ToPrettyDoc(ctx any) pretty.PrettyDoc {
	return pretty.Sequence(
		pretty.FromString("("),
		pretty.Indent(2, item.Child.ToPrettyDoc(ctx)),
		pretty.FromString(")"),
	)
}

type Var struct {
	Name string
}

func (v Var) sealed() {}

func (item Var) ToPrettyDoc(ctx any) pretty.PrettyDoc {
	return pretty.FromString(item.Name)
}

type Lambda struct {
	ArgName string
	Body    ParseTree
}

func (v Lambda) sealed() {}
func (item Lambda) ToPrettyDoc(ctx any) pretty.PrettyDoc {
	return pretty.Sequence(
		pretty.FromString(fmt.Sprint("\\", item.ArgName, ".")),
		pretty.Indent(2, item.Body.ToPrettyDoc(ctx)),
	)
}

type App struct {
	Callee ParseTree
	Args   AppArgs
}

func (item App) ToPrettyDoc(ctx any) pretty.PrettyDoc {
	firstArg := item.Args.First.ToPrettyDoc(ctx)
	moreArgs := make([]pretty.PrettyDoc, 0, len(item.Args.More))
	for _, arg := range item.Args.More {
		moreArgs = append(moreArgs, arg.ToPrettyDoc(ctx))
	}
	return pretty.Sequence(
		item.Callee.ToPrettyDoc(ctx),
		pretty.Indent(2, pretty.Sequence(firstArg, moreArgs...)),
	)
}

func (v App) sealed() {}

type AppArgs struct {
	First ParseTree
	More  []ParseTree
}

func (t ParseTree) ToPrettyDoc(ctx any) pretty.PrettyDoc {
	return t.Item.ToPrettyDoc(ctx)
}
func (t ParseTree) String() string {
	return t.ToPrettyDoc(nil).String()
}
