package parse_tree

import (
	"fmt"
	"strings"

	"github.com/gusbicalho/go-lambda/position"
	"github.com/gusbicalho/go-lambda/pretty"
)

type ParseTree struct {
	InputLocation position.Position
	Item          ParseItem
}

type ParseItem interface {
	pretty.Pretty
	sealed()
}

type Parens struct {
	Child ParseTree
}

func (v Parens) sealed() {}

func (item Parens) ToPrettyDoc() pretty.PrettyDoc {
	return pretty.Sequence(
		pretty.FromString("("),
		pretty.Indent(2, item.Child.ToPrettyDoc()),
		pretty.FromString(")"),
	)
}

type Var struct {
	Name string
}

func (v Var) sealed() {}

func (item Var) ToPrettyDoc() pretty.PrettyDoc {
	return pretty.FromString(item.Name)
}

type Lambda struct {
	ArgName string
	Body    ParseTree
}

func (v Lambda) sealed() {}
func (item Lambda) ToPrettyDoc() pretty.PrettyDoc {
	return pretty.Sequence(
		pretty.FromString(fmt.Sprint("\\", item.ArgName, ".")),
		pretty.Indent(2, item.Body.ToPrettyDoc()),
	)
}

type App struct {
	Callee ParseTree
	Args   AppArgs
}

func (item App) ToPrettyDoc() pretty.PrettyDoc {
	firstArg := item.Args.First.ToPrettyDoc()
	moreArgs := make([]pretty.PrettyDoc, 0, len(item.Args.More))
	for _, arg := range item.Args.More {
		moreArgs = append(moreArgs, arg.ToPrettyDoc())
	}
	return pretty.Sequence(
		item.Callee.ToPrettyDoc(),
		pretty.Indent(2, pretty.Sequence(firstArg, moreArgs...)),
	)
}

func (v App) sealed() {}

type AppArgs struct {
	First ParseTree
	More  []ParseTree
}

func (t ParseTree) ToPrettyDoc() pretty.PrettyDoc {
	return t.Item.ToPrettyDoc()
}
func (t ParseTree) String() string {
	return strings.Join(t.ToPrettyDoc().ToLines(0), "\n")
}
