package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	ln_beta_reduce "github.com/gusbicalho/go-lambda/locally_nameless/beta_reduce"
	ln_expr "github.com/gusbicalho/go-lambda/locally_nameless/expr"
	ln_pretty "github.com/gusbicalho/go-lambda/locally_nameless/pretty"
	"github.com/gusbicalho/go-lambda/locally_nameless/walk"
	"github.com/gusbicalho/go-lambda/parse_tree_to_locally_nameless"
	"github.com/gusbicalho/go-lambda/parser"
	"github.com/gusbicalho/go-lambda/tokenizer"

	"github.com/rivo/tview"
)

func main() {
	var source string
	if len(os.Args) > 1 {
		source = os.Args[1]
	} else {
		line, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			panic(err)
		}
		source = line
	}

	parseTree, err := parser.Parse(tokenizer.New(strings.NewReader(source)))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	expr := parse_tree_to_locally_nameless.ToLocallyNameless(*parseTree)

	//tui(expr)
	// tui2(expr)
	tui3(expr)
	//run(expr)
}

func tui(expr ln_expr.Expr) {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(
			func() {
				app.Draw()
			},
		)

	log := []string{ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName)}
	selectedRedexIndex := 0
	redexes := slices.Collect(ln_beta_reduce.BetaRedexes(expr))

	stop := func() {
		app.Stop()
		for _, logEntry := range log {
			fmt.Println(logEntry)
		}
	}

	getSelectedRedex := func() *ln_beta_reduce.BetaRedex {
		count := len(redexes)
		if count <= 0 {
			return nil
		}
		i := selectedRedexIndex % count
		if i < 0 {
			i += count
		}
		return &redexes[i]
	}

	step := func() {
		if redex := getSelectedRedex(); redex != nil {
			expr = redex.Reduce()
			log = append(log, ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName))
			selectedRedexIndex = 0
			redexes = slices.Collect(ln_beta_reduce.BetaRedexes(expr))
		} else {
			stop()
		}
	}

	shift := func(change int) {
		if len(redexes) <= 0 {
			selectedRedexIndex = 0
			return
		}
		selectedRedexIndex += change
	}

	redraw := func() {
		var pretty string
		if redex := getSelectedRedex(); redex != nil {
			pretty = redex.ToPrettyDoc(nil).String()
		} else {
			pretty = ln_pretty.ToPrettyDoc(expr).String() + "\nIrreducible."
		}

		textView.Clear()
		fmt.Fprintf(
			textView, "%s\n\n%s",
			ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName),
			pretty,
		)
	}
	go redraw()

	textView.SetDoneFunc(
		func(key tcell.Key) {
			switch key {
			case tcell.KeyESC:
				stop()
			case tcell.KeyEnter:
				step()
				redraw()
			case tcell.KeyTab:
				if len(redexes) > 0 {
					shift(1)
					redraw()
				}
			case tcell.KeyBacktab:
				if len(redexes) > 0 {
					shift(-1)
					redraw()
				}
			default:
			}
		},
	)

	textView.SetBorder(true)
	if err := app.SetRoot(textView, true).SetFocus(textView).Run(); err != nil {
		panic(err)
	}
}

func tui2(expr ln_expr.Expr) {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(
			func() {
				app.Draw()
			},
		)

	log := []string{ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName)}
	walking := walk.Pre(expr)

	stop := func() {
		app.Stop()
		for _, logEntry := range log {
			fmt.Println(logEntry)
		}
	}

	step := func() {
		if redex := ln_beta_reduce.AsBetaRedex(walking.Focus().Expr); redex != nil {
			walking = walking.UpdateExpr(func(_ ln_expr.Expr) ln_expr.Expr { return redex.Reduce() })
			expr = walking.Focus().Realize()
			log = append(log, ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName))
		}
	}

	shift := func(change int) {
		if change > 0 {
			if next := walking.Next(); next != nil {
				walking = next
			}
		} else {
			if prev := walking.Prev(); prev != nil {
				walking = prev
			}
		}
	}

	redraw := func() {
		var pretty = walk.ToPrettyDoc(walking).String()

		textView.Clear()
		fmt.Fprintf(
			textView, "%s\n\n%s\n\n%s\n%s",
			ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName),
			pretty,
			fmt.Sprint(expr),
			fmt.Sprint(walking),
		)
	}
	go redraw()

	textView.SetDoneFunc(
		func(key tcell.Key) {
			switch key {
			case tcell.KeyESC:
				stop()
			case tcell.KeyEnter:
				step()
				redraw()
			case tcell.KeyTab:
				shift(1)
				redraw()
			case tcell.KeyBacktab:
				shift(-1)
				redraw()
			default:
			}
		},
	)

	textView.SetBorder(true)
	if err := app.SetRoot(textView, true).SetFocus(textView).Run(); err != nil {
		panic(err)
	}
}

func tui3(expr ln_expr.Expr) {
	app := tview.NewApplication()
	textView := newExprView(app)

	log := []string{ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName)}
	nav := walk.ToNav(expr)

	stop := func() {
		app.Stop()
		for _, logEntry := range log {
			fmt.Println(logEntry)
		}
	}

	step := func() {
		var updated bool
		nav, updated = nav.UpdateExpr(func(e ln_expr.Expr) *ln_expr.Expr {
			redex := ln_beta_reduce.AsBetaRedex(e)
			if redex == nil {
				return nil
			}
			e = redex.Reduce()
			return &e
		})
		if updated {
			expr = nav.Focus().Realize()
			log = append(log, ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName))
		}
	}

	left := func() {
		if parent, _ := nav.Parent(); parent != nil {
			nav = *parent
		}
	}

	right := func() {
		if child := nav.Child(0); child != nil {
			nav = *child
		}
	}

	down := func() {
		n := &nav
		for {
			parent, index := n.Parent()
			if parent == nil {
				break
			}
			if sibling := parent.Child(index + 1); sibling != nil {
				nav = *sibling
				break
			}
			n = parent
		}
	}

	up := func() {
		parent, index := nav.Parent()
		if parent == nil {
			return
		}
		if index > 0 {
			if sibling := parent.Child(index - 1); sibling != nil {
				nav = *sibling
				if child := nav.Child(nav.Children() - 1); child != nil {
					nav = *child
				}
				return
			}
		}
		nav = *parent
	}

	redraw := func() {
		var pretty = nav.Focus().ToPrettyDoc().String()

		textView.Clear()
		fmt.Fprintf(
			textView, "%s\n\n%s\n\n%s\n%s",
			ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName),
			pretty,
			fmt.Sprint(expr),
			fmt.Sprint(nav),
		)
	}
	go redraw()

	textView.SetKeyHandler(func(key tcell.Key) bool {
		switch key {
		case tcell.KeyESC:
			stop()
		case tcell.KeyEnter:
			step()
			redraw()
		case tcell.KeyLeft:
			left()
			redraw()
		case tcell.KeyRight:
			right()
			redraw()
		case tcell.KeyUp:
			up()
			redraw()
		case tcell.KeyDown:
			down()
			redraw()
		default:
			return false
		}
		return true
	})

	textView.SetBorder(true)
	if err := app.SetRoot(textView, true).SetFocus(textView).Run(); err != nil {
		panic(err)
	}
}

type exprView struct {
	*tview.TextView
	onKey func(tcell.Key) bool
}

func newExprView(app *tview.Application) *exprView {
	return &exprView{
		tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetChangedFunc(
				func() {
					app.Draw()
				},
			),
		nil,
	}
}

func (t *exprView) SetKeyHandler(onKey func(key tcell.Key) bool) {
	t.onKey = onKey
}

func (t *exprView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	super := t.TextView.InputHandler()
	return t.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if t.onKey != nil && event != nil {
			key := event.Key()
			if t.onKey(key) {
				return
			}
		}
		super(event, setFocus)
	})
}

func run(expr ln_expr.Expr) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(ln_expr.ToLambdaNotation(expr, ln_expr.DisplayName))
		redex := nextBetaRedex(expr)
		if redex == nil {
			fmt.Println(ln_pretty.ToPrettyDoc(expr).String())
			fmt.Println("Irreducible.")
			break
		}
		fmt.Println(redex.ToPrettyDoc(nil).String())
		fmt.Print("Step? ")
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		expr = redex.Reduce()
	}
}

func nextBetaRedex(expr ln_expr.Expr) *ln_beta_reduce.BetaRedex {
	for redex := range ln_beta_reduce.BetaRedexes(expr) {
		return &redex
	}
	return nil
}
