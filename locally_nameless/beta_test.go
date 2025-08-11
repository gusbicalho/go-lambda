package locally_nameless

import (
	"reflect"
	"testing"
)

// assertExprEqual checks if two Expr values are deeply equal
func assertExprEqual(t *testing.T, expected, actual Expr, testName string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s failed:\nExpected: %#v\nActual:   %#v", testName, expected, actual)
	}
}

// Test Case 1: Identity Function
// Lambda: \x. x, Arg: a (FreeVar), Expected: a
func TestBetaReduce_IdentityFunction(t *testing.T) {
	// \x. x (De Bruijn: \. 0)
	lambda := NewLambda("x", NewBound(0))
	// a (FreeVar)
	arg := NewFree("a")

	actual := BetaReduce(lambda, arg)

	expected := NewFree("a")
	assertExprEqual(t, expected, actual, "Identity Function")
}

// Test Case 2: Constant Function
// Lambda: \x. y (where y is free), Arg: a (FreeVar), Expected: y
func TestBetaReduce_ConstantFunction(t *testing.T) {
	// \x. y (De Bruijn: \. y)
	lambda := NewLambda("x", NewFree("y"))
	// a (FreeVar)
	arg := NewFree("a")

	actual := BetaReduce(lambda, arg)

	expected := NewFree("y")
	assertExprEqual(t, expected, actual, "Constant Function")
}

// Test Case 3: Simple Application
// Lambda: \f. f a, Arg: \x. x (identity), Expected: (\x. x) a
func TestBetaReduce_SimpleApplication(t *testing.T) {
	// \f. f a (De Bruijn: \. 0 a)
	lambda := NewLambda("f", NewApp(NewBound(0), NewFree("a")))
	// \x. x (De Bruijn: \. 0)
	arg := NewLambda("x", NewBound(0))

	actual := BetaReduce(lambda, arg)

	expected := NewApp(NewLambda("x", NewBound(0)), NewFree("a"))
	assertExprEqual(t, expected, actual, "Simple Application")
}

// Test Case 4: Higher-Order Function
// Lambda: \f. \x. f x, Arg: g (FreeVar), Expected: \x. g x
func TestBetaReduce_HigherOrderFunction(t *testing.T) {
	// \f. \x. f x (De Bruijn: \. \. 1 0)
	lambda := NewLambda("f", NewLambda("x", NewApp(NewBound(1), NewBound(0))))
	// g (FreeVar)
	arg := NewFree("g")

	actual := BetaReduce(lambda, arg)

	expected := NewLambda("x", NewApp(NewFree("g"), NewBound(0)))
	assertExprEqual(t, expected, actual, "Higher-Order Function")
}

// Test Case 5: Multiple Bound Variables
// Lambda: \x. \y. \z. x z y, Arg: f (FreeVar), Expected: \y. \z. f z y
func TestBetaReduce_MultipleBoundVariables(t *testing.T) {
	// \x. \y. \z. x z y (De Bruijn: \. \. \. 2 0 1)
	lambda := NewLambda("x",
		NewLambda("y",
			NewLambda("z",
				NewApp(NewApp(NewBound(2), NewBound(0)), NewBound(1)))))
	// f (FreeVar)
	arg := NewFree("f")

	actual := BetaReduce(lambda, arg)

	expected := NewLambda("y", NewLambda("z", NewApp(NewApp(NewFree("f"), NewBound(0)), NewBound(1))))
	assertExprEqual(t, expected, actual, "Multiple Bound Variables")
}

// Test Case 6: Self-Application
// Lambda: \x. x x, Arg: \y. y (identity), Expected: (\y. y) (\y. y)
func TestBetaReduce_SelfApplication(t *testing.T) {
	// \x. x x (De Bruijn: \. 0 0)
	lambda := NewLambda("x", NewApp(NewBound(0), NewBound(0)))
	// \y. y (De Bruijn: \. 0)
	arg := NewLambda("y", NewBound(0))

	actual := BetaReduce(lambda, arg)

	expected := NewApp(NewLambda("y", NewBound(0)), NewLambda("y", NewBound(0)))
	assertExprEqual(t, expected, actual, "Self-Application")
}

// Test Case 7: Church Numerals
// Lambda: \f. \x. f (f x) (Church 2), Arg: succ (FreeVar), Expected: \x. succ (succ x)
func TestBetaReduce_ChurchNumerals(t *testing.T) {
	// \f. \x. f (f x) (De Bruijn: \. \. 1 (1 0))
	lambda := NewLambda("f",
		NewLambda("x",
			NewApp(NewBound(1), NewApp(NewBound(1), NewBound(0)))))
	// succ (FreeVar)
	arg := NewFree("succ")

	actual := BetaReduce(lambda, arg)

	expected := NewLambda("x", NewApp(NewFree("succ"), NewApp(NewFree("succ"), NewBound(0))))
	assertExprEqual(t, expected, actual, "Church Numerals")
}

// Test Case 8: Substitution Under Nested Lambda
// Lambda: \x. \y. x, Arg: \z. y (contains free var y), Expected: \y'. (\z. y)
func TestBetaReduce_SubstitutionUnderNestedLambda(t *testing.T) {
	// \x. \y. x (De Bruijn: \. \. 1)
	lambda := NewLambda("x", NewLambda("y", NewBound(1)))
	// \z. y (De Bruijn: \. y)
	arg := NewLambda("z", NewFree("y"))

	actual := BetaReduce(lambda, arg)

	// Expected: \y'. (\z. y) - the lambda parameter name can be anything, body remains (\z. y)
	expected := NewLambda("y", NewLambda("z", NewFree("y")))
	assertExprEqual(t, expected, actual, "Substitution Under Nested Lambda")
}

// Test Case 9: Variable Capture Prevention
// Lambda: \x. \y. x y, Arg: \a. \b. y a (free y, bound a), Expected: \y'. (\a. \b. y a) y'
func TestBetaReduce_VariableCapturePrevention(t *testing.T) {
	// \x. \y. x y (De Bruijn: \. \. 1 0)
	lambda := NewLambda("x", NewLambda("y", NewApp(NewBound(1), NewBound(0)))) // \a. \b. y a (De Bruijn: \. \. y 1)
	arg := NewLambda("a", NewLambda("b", NewApp(NewFree("y"), NewBound(1))))

	actual := BetaReduce(lambda, arg)

	// Expected: \y'. (\a. \b. y a) y' where the argument is applied to the bound variable
	expected := NewLambda("y", NewApp(NewLambda("a", NewLambda("b", NewApp(NewFree("y"), NewBound(1)))), NewBound(0)))
	assertExprEqual(t, expected, actual, "Variable Capture Prevention")
}

// Test Case 10: Free Variables in Argument
// Lambda: \f. \x. f (g x) (g is free), Arg: \y. g y (same g), Expected: \x. (\y. g y) (g x)
func TestBetaReduce_FreeVariablesInArgument(t *testing.T) {
	// \f. \x. f (g x) (De Bruijn: \. \. 1 (g 0))
	lambda := NewLambda("f",
		NewLambda("x",
			NewApp(NewBound(1), NewApp(NewFree("g"), NewBound(0))))) // \y. g y (De Bruijn: \. g 0)
	arg := NewLambda("y", NewApp(NewFree("g"), NewBound(0)))

	actual := BetaReduce(lambda, arg)

	// Expected: \x. (\y. g y) (g x)
	expected := NewLambda("x", NewApp(NewLambda("y", NewApp(NewFree("g"), NewBound(0))), NewApp(NewFree("g"), NewBound(0))))
	assertExprEqual(t, expected, actual, "Free Variables in Argument")
}

// Test Case 11: Deeply Nested Lambdas
// Lambda: \x. \a. \b. \c. \d. x a b c d, Arg: \p. \q. p q, Expected: \a. \b. \c. \d. (\p. \q. p q) a b c d
func TestBetaReduce_DeeplyNestedLambdas(t *testing.T) {
	// \x. \a. \b. \c. \d. x a b c d (De Bruijn: \. \. \. \. \. 4 3 2 1 0)
	lambda := NewLambda("x",
		NewLambda("a",
			NewLambda("b",
				NewLambda("c",
					NewLambda("d",
						NewApp(NewApp(NewApp(NewApp(NewBound(4), NewBound(3)), NewBound(2)), NewBound(1)), NewBound(0))))))) // \p. \q. p q (De Bruijn: \. \. 1 0)
	arg := NewLambda("p", NewLambda("q", NewApp(NewBound(1), NewBound(0))))

	actual := BetaReduce(lambda, arg)

	// Expected: \a. \b. \c. \d. (\p. \q. p q) a b c d
	expected := NewLambda("a", NewLambda("b", NewLambda("c", NewLambda("d",
		NewApp(NewApp(NewApp(NewApp(
			NewLambda("p", NewLambda("q", NewApp(NewBound(1), NewBound(0)))),
			NewBound(3)), NewBound(2)), NewBound(1)), NewBound(0))))))
	assertExprEqual(t, expected, actual, "Deeply Nested Lambdas")
}

// Test Case 12: Church Y Combinator Pattern
// Lambda: \f. (\x. f (x x)) (\x. f (x x)), Arg: g (FreeVar), Expected: (\x. g (x x)) (\x. g (x x))
func TestBetaReduce_ChurchYCombinatorPattern(t *testing.T) {
	// \f. (\x. f (x x)) (\x. f (x x)) (De Bruijn: \. (\. 1 (0 0)) (\. 1 (0 0)))
	xTerm := NewLambda("x", NewApp(NewBound(1), NewApp(NewBound(0), NewBound(0))))
	lambda := NewLambda("f", NewApp(xTerm, xTerm)) // g (FreeVar)
	arg := NewFree("g")

	actual := BetaReduce(lambda, arg)

	// Expected: (\x. g (x x)) (\x. g (x x))
	xTermResult := NewLambda("x", NewApp(NewFree("g"), NewApp(NewBound(0), NewBound(0))))
	expected := NewApp(xTermResult, xTermResult)
	assertExprEqual(t, expected, actual, "Church Y Combinator Pattern")
}

// Test Case 13: Multiple Nested Applications
// Lambda: \f. ((f a) b) c, Arg: \x. \y. \z. x (y z), Expected: (((\x. \y. \z. x (y z)) a) b) c
func TestBetaReduce_MultipleNestedApplications(t *testing.T) {
	// \f. ((f a) b) c (De Bruijn: \. ((0 a) b) c)
	lambda := NewLambda("f",
		NewApp(NewApp(NewApp(NewBound(0), NewFree("a")), NewFree("b")), NewFree("c"))) // \x. \y. \z. x (y z) (De Bruijn: \. \. \. 2 (1 0))
	arg := NewLambda("x",
		NewLambda("y",
			NewLambda("z",
				NewApp(NewBound(2), NewApp(NewBound(1), NewBound(0))))))

	actual := BetaReduce(lambda, arg)

	// Expected: (((\x. \y. \z. x (y z)) a) b) c
	expected := NewApp(NewApp(NewApp(
		NewLambda("x", NewLambda("y", NewLambda("z", NewApp(NewBound(2), NewApp(NewBound(1), NewBound(0)))))),
		NewFree("a")), NewFree("b")), NewFree("c"))
	assertExprEqual(t, expected, actual, "Multiple Nested Applications")
}

// Test Case 14: Empty Body Reference
// Lambda: \x. y (x never used), Arg: complex expression, Expected: y (argument ignored)
func TestBetaReduce_EmptyBodyReference(t *testing.T) {
	// \x. y (De Bruijn: \. y)
	lambda := NewLambda("x", NewFree("y")) // Complex argument: \a. \b. a (b a)
	arg := NewLambda("a", NewLambda("b", NewApp(NewBound(1), NewApp(NewBound(0), NewBound(1)))))

	actual := BetaReduce(lambda, arg)

	// Expected: y (argument completely ignored)
	expected := NewFree("y")
	assertExprEqual(t, expected, actual, "Empty Body Reference")
}

// Test Case 15: Argument with Same Structure
// Lambda: \x. \y. x y, Arg: \y. z y (shadow y in arg), Expected: \y'. (\y. z y) y'
func TestBetaReduce_ArgumentWithSameStructure(t *testing.T) {
	// \x. \y. x y (De Bruijn: \. \. 1 0)
	lambda := NewLambda("x", NewLambda("y", NewApp(NewBound(1), NewBound(0)))) // \y. z y (De Bruijn: \. z 0)
	arg := NewLambda("y", NewApp(NewFree("z"), NewBound(0)))

	actual := BetaReduce(lambda, arg)

	// Expected: \y'. (\y. z y) y' - argument applied to bound variable with alpha conversion
	expected := NewLambda("y", NewApp(NewLambda("y", NewApp(NewFree("z"), NewBound(0))), NewBound(0)))
	assertExprEqual(t, expected, actual, "Argument with Same Structure")
}

// Test Case 16: BoundVar as Argument
// DisplayContext: \y. (\x. x) y, Lambda: \x. x, Arg: y (BoundVar), Expected: y (same BoundVar)
func TestBetaReduce_BoundVarAsArgument(t *testing.T) {
	// \x. x (De Bruijn: \. 0)
	lambda := NewLambda("x", NewBound(0)) // y as BoundVar (index 0 - referring to outer lambda in context)
	arg := NewBound(0)

	actual := BetaReduce(lambda, arg)

	// Expected: y (BoundVar with same index 0)
	expected := NewBound(0)
	assertExprEqual(t, expected, actual, "BoundVar as Argument")
}

// Test Case 17: Application as Argument
// Lambda: \f. f, Arg: g h (application), Expected: g h
func TestBetaReduce_ApplicationAsArgument(t *testing.T) {
	// \f. f (De Bruijn: \. 0)
	lambda := NewLambda("f", NewBound(0)) // g h (application)
	arg := NewApp(NewFree("g"), NewFree("h"))

	actual := BetaReduce(lambda, arg)

	// Expected: g h (same application)
	expected := NewApp(NewFree("g"), NewFree("h"))
	assertExprEqual(t, expected, actual, "Application as Argument")
}

// Test Case 18: Index Arithmetic Boundary
// Lambda: \x. \y. \z. y, Arg: a (FreeVar), Expected: \y. \z. y (index 1 unchanged, verify boundary)
func TestBetaReduce_IndexArithmeticBoundary(t *testing.T) {
	// \x. \y. \z. y (De Bruijn: \. \. \. 1)
	lambda := NewLambda("x",
		NewLambda("y",
			NewLambda("z", NewBound(1)))) // a (FreeVar)
	arg := NewFree("a")

	actual := BetaReduce(lambda, arg)

	// Expected: \y. \z. y (bound var index 1 unchanged after substitution)
	expected := NewLambda("y", NewLambda("z", NewBound(1)))
	assertExprEqual(t, expected, actual, "Index Arithmetic Boundary")
}
