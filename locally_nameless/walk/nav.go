package walk

import (
	"slices"

	"github.com/gusbicalho/go-lambda/locally_nameless/expr"
	"github.com/gusbicalho/go-lambda/locally_nameless/hole"
)

type Nav struct {
	parent *navParent
	expr   expr.Expr
}

type navParent struct {
	parent *navParent
	hole   hole.Hole
	index  uint
}

func ToNav(e expr.Expr) Nav {
	return Nav{
		parent: nil,
		expr:   e,
	}
}

func (nav Nav) Parent() (*Nav, uint) {
	if nav.parent == nil {
		return nil, 0
	}

	parent := nav.parent.parent
	hole := nav.parent.hole
	index := nav.parent.index
	return &Nav{parent: parent, expr: hole.Fill(nav.expr)}, index
}

func (nav Nav) UpdateExpr(update func(expr.Expr) *expr.Expr) (Nav, bool) {
	newExpr := update(nav.expr)
	if newExpr != nil {
		nav.expr = *newExpr
		return nav, true
	}
	return nav, false
}

func (nav Nav) Focus() Focus {
	holes := []hole.Hole{}
	for p := nav.parent; p != nil; p = p.parent {
		holes = append(holes, p.hole)
	}
	slices.Reverse(holes)
	hole := hole.ComposeHoles(holes...)
	return Focus{Hole: hole, Expr: nav.expr}
}

func (nav Nav) Children() uint {
	return expr.CaseExpr(nav.expr, visitNavChildren{})
}

func (nav Nav) Child(child uint) *Nav {
	return expr.CaseExpr(nav.expr, visitNavChild{nav.parent, child})
}

type visitNavChildren struct{}

func (v visitNavChildren) CaseBound(_ expr.BoundVar) uint { return 0 }
func (v visitNavChildren) CaseFree(_ expr.FreeVar) uint   { return 0 }
func (v visitNavChildren) CaseLambda(e expr.Lambda) uint  { return 1 }
func (v visitNavChildren) CaseApp(e expr.App) uint        { return 2 }

type visitNavChild struct {
	parent *navParent
	child  uint
}

func (v visitNavChild) CaseBound(_ expr.BoundVar) *Nav {
	return nil
}

func (v visitNavChild) CaseFree(_ expr.FreeVar) *Nav {
	return nil
}

func (v visitNavChild) CaseLambda(e expr.Lambda) *Nav {
	if v.child != 0 {
		return nil
	}
	return &Nav{
		parent: &navParent{
			parent: v.parent,
			hole:   hole.BodyHole(e),
			index:  v.child,
		},
		expr: e.Body(),
	}
}
func (v visitNavChild) CaseApp(e expr.App) *Nav {
	switch v.child {
	case 0:
		return &Nav{
			parent: &navParent{
				parent: v.parent,
				hole:   hole.CalleeHole(e),
				index:  v.child,
			},
			expr: e.Callee(),
		}
	case 1:
		return &Nav{
			parent: &navParent{
				parent: v.parent,
				hole:   hole.ArgHole(e),
				index:  v.child,
			},
			expr: e.Arg(),
		}
	default:
		return nil
	}
}
