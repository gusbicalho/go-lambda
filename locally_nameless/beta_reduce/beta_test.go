package beta_reduce

import (
	"fmt"
	"testing"

	"github.com/gusbicalho/go-lambda/locally_nameless/expr"
)

func assertExprRendersAs(t *testing.T, testName string, e expr.Expr, expected string) {
	actual := expr.ToLambdaNotation(e, expr.DisplayName)
	if actual != expected {
		t.Errorf("%s - Expected: %#v\nActual:   %#v", testName, expected, actual)
	}
}

func TestBetaReduce(t *testing.T) {
	cases := []struct {
		testName  string
		lambda    expr.Lambda
		arg       expr.Expr
		rendersAs string
		reducesTo string
	}{
		{
			"Identity Function",
			expr.NewLambda("x", expr.NewBound(0)),
			expr.NewFree("a"),
			"(\\x. x) a",
			"a"},
		{
			"Constant Function",
			expr.NewLambda("x", expr.NewFree("y")),
			expr.NewFree("a"),
			"(\\x. y) a",
			"y",
		},
		{
			"Simple Application",
			expr.NewLambda("f", expr.NewApp(expr.NewBound(0), expr.NewFree("a"))),
			expr.NewLambda("x", expr.NewBound(0)),
			"(\\f. f a) (\\x. x)",
			"(\\x. x) a",
		},
		{
			"Higher-Order Function",
			expr.NewLambda("f", expr.NewLambda("x", expr.NewApp(expr.NewBound(1), expr.NewBound(0)))),
			expr.NewFree("g"),
			"(\\f. \\x. f x) g",
			"\\x. g x",
		},
		{
			"Multiple Bound Variables",
			expr.NewLambda("x", expr.NewLambda("y", expr.NewLambda("z", expr.NewApp(expr.NewApp(expr.NewBound(2), expr.NewBound(0)), expr.NewBound(1))))),
			expr.NewFree("f"),
			"(\\x. \\y. \\z. x z y) f",
			"\\y. \\z. f z y",
		},
		{
			"Self-Application",
			expr.NewLambda("x", expr.NewApp(expr.NewBound(0), expr.NewBound(0))),
			expr.NewLambda("y", expr.NewBound(0)),
			"(\\x. x x) (\\y. y)",
			"(\\y. y) (\\y. y)",
		},
		{
			"Church Numerals",
			expr.NewLambda("f", expr.NewLambda("x", expr.NewApp(expr.NewBound(1), expr.NewApp(expr.NewBound(1), expr.NewBound(0))))),
			expr.NewFree("succ"),
			"(\\f. \\x. f (f x)) succ",
			"\\x. succ (succ x)",
		},
		{
			"Substitution Under Nested expr.Lambda",
			expr.NewLambda("x", expr.NewLambda("y", expr.NewBound(1))),
			expr.NewLambda("z", expr.NewFree("y")),
			"(\\x. \\y. x) (\\z. y)",
			"\\y. \\z. y",
		},
		{
			"Variable Capture Prevention",
			expr.NewLambda("x", expr.NewLambda("y", expr.NewApp(expr.NewBound(1), expr.NewBound(0)))),
			expr.NewLambda("a", expr.NewLambda("b", expr.NewApp(expr.NewFree("y"), expr.NewBound(1)))),
			"(\\x. \\y. x y) (\\a. \\b. y a)",
			"\\y. (\\a. \\b. y a) y",
		},
		{
			"Free Variables in Argument",
			expr.NewLambda("f", expr.NewLambda("x", expr.NewApp(expr.NewBound(1), expr.NewApp(expr.NewFree("g"), expr.NewBound(0))))),
			expr.NewLambda("y", expr.NewApp(expr.NewFree("g"), expr.NewBound(0))),
			"(\\f. \\x. f (g x)) (\\y. g y)",
			"\\x. (\\y. g y) (g x)",
		},
		{
			"Deeply Nested Lambdas",
			expr.NewLambda("x", expr.NewLambda("a", expr.NewLambda("b", expr.NewLambda("c", expr.NewLambda("d", expr.NewApp(expr.NewApp(expr.NewApp(expr.NewApp(expr.NewBound(4), expr.NewBound(3)), expr.NewBound(2)), expr.NewBound(1)), expr.NewBound(0))))))),
			expr.NewLambda("p", expr.NewLambda("q", expr.NewApp(expr.NewBound(1), expr.NewBound(0)))),
			"(\\x. \\a. \\b. \\c. \\d. x a b c d) (\\p. \\q. p q)",
			"\\a. \\b. \\c. \\d. (\\p. \\q. p q) a b c d",
		},
		{
			"Church Y Combinator Pattern",
			expr.NewLambda("f", expr.NewApp(expr.NewLambda("x", expr.NewApp(expr.NewBound(1), expr.NewApp(expr.NewBound(0), expr.NewBound(0)))), expr.NewLambda("x", expr.NewApp(expr.NewBound(1), expr.NewApp(expr.NewBound(0), expr.NewBound(0)))))),
			expr.NewFree("g"),
			"(\\f. (\\x. f (x x)) (\\x. f (x x))) g",
			"(\\x. g (x x)) (\\x. g (x x))",
		},
		{
			"Multiple Nested Applications",
			expr.NewLambda("f", expr.NewApp(expr.NewApp(expr.NewApp(expr.NewBound(0), expr.NewFree("a")), expr.NewFree("b")), expr.NewFree("c"))),
			expr.NewLambda("x", expr.NewLambda("y", expr.NewLambda("z", expr.NewApp(expr.NewBound(2), expr.NewApp(expr.NewBound(1), expr.NewBound(0)))))),
			"(\\f. f a b c) (\\x. \\y. \\z. x (y z))",
			"(\\x. \\y. \\z. x (y z)) a b c",
		},
		{
			"Empty Body Reference",
			expr.NewLambda("x", expr.NewFree("y")),
			expr.NewLambda("a", expr.NewLambda("b", expr.NewApp(expr.NewBound(1), expr.NewApp(expr.NewBound(0), expr.NewBound(1))))),
			"(\\x. y) (\\a. \\b. a (b a))",
			"y",
		},
		{
			"Argument with Same Structure",
			expr.NewLambda("x", expr.NewLambda("y", expr.NewApp(expr.NewBound(1), expr.NewBound(0)))),
			expr.NewLambda("y", expr.NewApp(expr.NewFree("z"), expr.NewBound(0))),
			"(\\x. \\y. x y) (\\y. z y)",
			"\\y. (\\y_0. z y_0) y",
		},
		{
			"Out-of-scope BoundVar as Argument",
			expr.NewLambda("x", expr.NewBound(0)),
			expr.NewBound(0),
			"(\\x. x) 0:<outofscope>",
			"0:<outofscope>",
		},
		{
			"Application as Argument",
			expr.NewLambda("f", expr.NewBound(0)),
			expr.NewApp(expr.NewFree("g"), expr.NewFree("h")),
			"(\\f. f) (g h)",
			"g h",
		},
		{
			"Index Arithmetic Boundary",
			expr.NewLambda("x", expr.NewLambda("y", expr.NewLambda("z", expr.NewBound(1)))),
			expr.NewFree("a"),
			"(\\x. \\y. \\z. y) a",
			"\\y. \\z. y",
		},
	}
	for i, c := range cases {
		testName := fmt.Sprint("Case ", i+1, " :", c.testName)
		assertExprRendersAs(t, testName, expr.NewApp(c.lambda, c.arg), c.rendersAs)
		assertExprRendersAs(t, testName, BetaReduce(c.lambda, c.arg), c.reducesTo)
	}
}
