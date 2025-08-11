package locally_nameless

import (
	"fmt"
	"testing"
)

func assertExprRendersAs(t *testing.T, testName string, expr Expr, expected string) {
	actual := ToLambdaNotation(expr, DisplayName)
	if actual != expected {
		t.Errorf("%s - Expected: %#v\nActual:   %#v", testName, expected, actual)
	}
}

func TestBetaReduce(t *testing.T) {
	cases := []struct {
		testName  string
		lambda    Lambda
		arg       Expr
		rendersAs string
		reducesTo string
	}{
		{
			"Identity Function",
			NewLambda("x", NewBound(0)),
			NewFree("a"),
			"(\\x. x) a",
			"a"},
		{
			"Constant Function",
			NewLambda("x", NewFree("y")),
			NewFree("a"),
			"(\\x. y) a",
			"y",
		},
		{
			"Simple Application",
			NewLambda("f", NewApp(NewBound(0), NewFree("a"))),
			NewLambda("x", NewBound(0)),
			"(\\f. f a) (\\x. x)",
			"(\\x. x) a",
		},
		{
			"Higher-Order Function",
			NewLambda("f", NewLambda("x", NewApp(NewBound(1), NewBound(0)))),
			NewFree("g"),
			"(\\f. \\x. f x) g",
			"\\x. g x",
		},
		{
			"Multiple Bound Variables",
			NewLambda("x", NewLambda("y", NewLambda("z", NewApp(NewApp(NewBound(2), NewBound(0)), NewBound(1))))),
			NewFree("f"),
			"(\\x. \\y. \\z. x z y) f",
			"\\y. \\z. f z y",
		},
		{
			"Self-Application",
			NewLambda("x", NewApp(NewBound(0), NewBound(0))),
			NewLambda("y", NewBound(0)),
			"(\\x. x x) (\\y. y)",
			"(\\y. y) (\\y. y)",
		},
		{
			"Church Numerals",
			NewLambda("f", NewLambda("x", NewApp(NewBound(1), NewApp(NewBound(1), NewBound(0))))),
			NewFree("succ"),
			"(\\f. \\x. f (f x)) succ",
			"\\x. succ (succ x)",
		},
		{
			"Substitution Under Nested Lambda",
			NewLambda("x", NewLambda("y", NewBound(1))),
			NewLambda("z", NewFree("y")),
			"(\\x. \\y. x) (\\z. y)",
			"\\y. \\z. y",
		},
		{
			"Variable Capture Prevention",
			NewLambda("x", NewLambda("y", NewApp(NewBound(1), NewBound(0)))),
			NewLambda("a", NewLambda("b", NewApp(NewFree("y"), NewBound(1)))),
			"(\\x. \\y. x y) (\\a. \\b. y a)",
			"\\y. (\\a. \\b. y a) y",
		},
		{
			"Free Variables in Argument",
			NewLambda("f", NewLambda("x", NewApp(NewBound(1), NewApp(NewFree("g"), NewBound(0))))),
			NewLambda("y", NewApp(NewFree("g"), NewBound(0))),
			"(\\f. \\x. f (g x)) (\\y. g y)",
			"\\x. (\\y. g y) (g x)",
		},
		{
			"Deeply Nested Lambdas",
			NewLambda("x", NewLambda("a", NewLambda("b", NewLambda("c", NewLambda("d", NewApp(NewApp(NewApp(NewApp(NewBound(4), NewBound(3)), NewBound(2)), NewBound(1)), NewBound(0))))))),
			NewLambda("p", NewLambda("q", NewApp(NewBound(1), NewBound(0)))),
			"(\\x. \\a. \\b. \\c. \\d. x a b c d) (\\p. \\q. p q)",
			"\\a. \\b. \\c. \\d. (\\p. \\q. p q) a b c d",
		},
		{
			"Church Y Combinator Pattern",
			NewLambda("f", NewApp(NewLambda("x", NewApp(NewBound(1), NewApp(NewBound(0), NewBound(0)))), NewLambda("x", NewApp(NewBound(1), NewApp(NewBound(0), NewBound(0)))))),
			NewFree("g"),
			"(\\f. (\\x. f (x x)) (\\x. f (x x))) g",
			"(\\x. g (x x)) (\\x. g (x x))",
		},
		{
			"Multiple Nested Applications",
			NewLambda("f", NewApp(NewApp(NewApp(NewBound(0), NewFree("a")), NewFree("b")), NewFree("c"))),
			NewLambda("x", NewLambda("y", NewLambda("z", NewApp(NewBound(2), NewApp(NewBound(1), NewBound(0)))))),
			"(\\f. f a b c) (\\x. \\y. \\z. x (y z))",
			"(\\x. \\y. \\z. x (y z)) a b c",
		},
		{
			"Empty Body Reference",
			NewLambda("x", NewFree("y")),
			NewLambda("a", NewLambda("b", NewApp(NewBound(1), NewApp(NewBound(0), NewBound(1))))),
			"(\\x. y) (\\a. \\b. a (b a))",
			"y",
		},
		{
			"Argument with Same Structure",
			NewLambda("x", NewLambda("y", NewApp(NewBound(1), NewBound(0)))),
			NewLambda("y", NewApp(NewFree("z"), NewBound(0))),
			"(\\x. \\y. x y) (\\y. z y)",
			"\\y. (\\y_0. z y_0) y",
		},
		{
			"Out-of-scope BoundVar as Argument",
			NewLambda("x", NewBound(0)),
			NewBound(0),
			"(\\x. x) 0:<outofscope>",
			"0:<outofscope>",
		},
		{
			"Application as Argument",
			NewLambda("f", NewBound(0)),
			NewApp(NewFree("g"), NewFree("h")),
			"(\\f. f) (g h)",
			"g h",
		},
		{
			"Index Arithmetic Boundary",
			NewLambda("x", NewLambda("y", NewLambda("z", NewBound(1)))),
			NewFree("a"),
			"(\\x. \\y. \\z. y) a",
			"\\y. \\z. y",
		},
	}
	for i, c := range cases {
		testName := fmt.Sprint("Case ", i+1, " :", c.testName)
		assertExprRendersAs(t, testName, NewApp(c.lambda, c.arg), c.rendersAs)
		assertExprRendersAs(t, testName, BetaReduce(c.lambda, c.arg), c.reducesTo)
	}
}
