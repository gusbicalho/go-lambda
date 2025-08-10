package lazy

type Lazy[v any] struct {
	*lazyImpl[v]
}

func New[v any](run func() v) Lazy[v] {
	return Lazy[v]{&lazyImpl[v]{todo: run}}
}

func Wrap[v any](value v) Lazy[v] {
	return Lazy[v]{&lazyImpl[v]{value: &value}}
}

func (l Lazy[v]) Get() v {
	if l.todo == nil {
		return *l.value
	}
	value := l.todo()
	l.value = &value
	l.todo = nil
	return value
}

type lazyImpl[v any] struct {
	todo  func() v
	value *v
}
