package expr

import (
	"fmt"

	"github.com/gusbicalho/go-lambda/stack"
)

type DisplayContext struct {
	bound             stack.Stack[string]
	displayBoundVarAs DisplayBoundVarAs
}

type DisplayBoundVarAs = uint

const (
	DisplayBoth DisplayBoundVarAs = iota
	DisplayName
	DisplayIndex
)

func EmptyContext() DisplayContext {
	return DisplayContext{}
}

func (ctx DisplayContext) WithDisplayBoundVarAs(displayAs DisplayBoundVarAs) DisplayContext {
	ctx.displayBoundVarAs = displayAs
	return ctx
}

func (ctx DisplayContext) BindFree(name string) (DisplayContext, string) {
	if ctx.isBound(name) {
		for i := 0; ; i++ {
			newName := fmt.Sprint(name, "_", i)
			if !ctx.isBound(newName) {
				return ctx.addBinding(newName), newName
			}
		}
	}
	return ctx.addBinding(name), name
}

func (ctx DisplayContext) addBinding(name string) DisplayContext {
	ctx.bound = ctx.bound.Push(name)
	return ctx
}

func (ctx DisplayContext) isBound(name string) bool {
	for bound := range ctx.bound.Items() {
		if name == bound {
			return true
		}
	}
	return false
}
