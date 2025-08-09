package stack

import "iter"

type Stack[v any] struct {
	head *node[v]
}

type node[v any] struct {
	value v
	next  *node[v]
}

func Empty[v any]() Stack[v] { return Stack[v]{head: nil} }

func (s Stack[v]) Push(value v) Stack[v] {
	return Stack[v]{
		head: &node[v]{
			value: value,
			next:  s.head,
		},
	}
}

type Popped[v any] struct {
	Value v
	Stack Stack[v]
}

func (s Stack[v]) Pop() *Popped[v] {
	if s.head == nil {
		return nil
	}
	return &Popped[v]{
		Value: s.head.value,
		Stack: Stack[v]{head: s.head.next},
	}
}

func (s Stack[v]) Nth(nth uint, defaultValue v) (v, bool) {
	for node := s.head; node != nil; node = node.next {
		if nth == 0 {
			return node.value, true
		}
		nth -= 1
	}
	return defaultValue, false
}

func (s Stack[v]) Items() iter.Seq[v] {
	return func(yield func(v) bool) {
		for node := s.head; node != nil; node = node.next {
			if !yield(node.value) {
				return
			}
		}
	}
}
func (s Stack[v]) IndexedItems() iter.Seq2[uint, v] {
	return func(yield func(uint, v) bool) {
		var i uint = 0
		for node := s.head; node != nil; node = node.next {
			if !yield(i, node.value) {
				return
			}
			i++
		}
	}
}
