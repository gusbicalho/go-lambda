package pretty

import (
	"strings"
)

type Pretty[context any] interface {
	ToPrettyDoc(context context) PrettyDoc
}

type PrettyDoc struct{ impl prettyDocImpl }

func (d PrettyDoc) String() string {
	return strings.Join(d.toLines(0), "\n")
}

func (d PrettyDoc) ToLines(indent uint) []string {
	return d.impl.toLines(indent)
}

func (d PrettyDoc) toLines(indent uint) []string {
	return d.impl.toLines(indent)
}

type prettyDocImpl interface {
	toLines(indent uint) []string
}

func FromString(s string) PrettyDoc {
	lines := strings.Split(s, "\n")
	switch len(lines) {
	case 0:
		return PrettyDoc{lineDoc{line: ""}}
	case 1:
		return PrettyDoc{lineDoc{line: lines[0]}}
	default:
		lineDocs := make([]prettyDocImpl, 0, len(lines))
		for i := 0; i < len(lines); i++ {
			lineDocs = append(lineDocs, lineDoc{line: lines[i]})
		}
		return PrettyDoc{sequenceDoc{items: lineDocs}}
	}
}

func Indent(indent uint, doc PrettyDoc) PrettyDoc {
	return PrettyDoc{indentDoc{indent: indent, item: doc}}
}

func Sequence(doc PrettyDoc, moreDocs ...PrettyDoc) PrettyDoc {
	if len(moreDocs) == 0 {
		return doc
	}
	items := make([]prettyDocImpl, 1+len(moreDocs))
	items[0] = doc.impl
	for i, doc := range moreDocs {
		items[i+1] = doc.impl
	}
	return PrettyDoc{sequenceDoc{items: items}}
}

type lineDoc struct {
	line string
}

func (s lineDoc) toLines(indent uint) []string {
	builder := strings.Builder{}
	builder.Grow(int(indent) + len([]byte(s.line)))
	for i := uint(0); i < indent; i++ {
		builder.WriteString(" ")
	}
	builder.WriteString(s.line)
	return []string{builder.String()}
}

type indentDoc struct {
	indent uint
	item   prettyDocImpl
}

func (i indentDoc) toLines(indent uint) []string {
	return i.item.toLines(indent + i.indent)
}

type sequenceDoc struct {
	items []prettyDocImpl
}

func (seq sequenceDoc) toLines(indent uint) []string {
	lines := make([]string, 0, len(seq.items))
	for i := 0; i < len(seq.items); i++ {
		lines = append(lines, seq.items[i].toLines(indent)...)
	}
	return lines
}
