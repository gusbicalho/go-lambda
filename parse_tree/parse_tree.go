package parse_tree

type ParseTree struct {
	InputRange InputRange
	Item       ParseItem
}

type InputRange struct {
	From InputLocation
	To   InputLocation
}

type InputLocation struct {
	Line   uint64
	Column uint64
}

type ParseItem interface {
	sealed()
}

type Parens struct {
	Child ParseTree
}

func (v Parens) sealed() {}

type Var struct {
	Name string
}

func (v Var) sealed() {}

type Lambda struct {
	ArgName string
	Body    []ParseTree
}

func (v Lambda) sealed() {}

type App struct {
	Fun  ParseTree
	Args []AppArgs
}

func (v App) sealed() {}

type AppArgs struct {
	First ParseTree
	More  []ParseTree
}

func Case[r any](item ParseItem,
	onParens func(item Parens) (r, error),
	onVar func(item Var) (r, error),
	onLambda func(item Lambda) (r, error),
	onApp func(item App) (r, error),
) (r, error) {
	switch item := item.(type) {
	case Parens:
		return onParens(item)
	case Var:
		return onVar(item)
	case Lambda:
		return onLambda(item)
	case App:
		return onApp(item)
	default:
		panic("Impossible")
	}
}
