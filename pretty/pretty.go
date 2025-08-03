package pretty

import (
	"strings"
)

type Pretty interface {
	ToPrettyDoc() PrettyDoc
}

type PrettyDoc interface {
	ToLines(indent uint) []string
	sealed()
}

func FromString(s string) PrettyDoc {
	lines := strings.Split(s, "\n")
	switch len(lines) {
	case 0:
		return lineDoc{line: ""}
	case 1:
		return lineDoc{line: lines[0]}
	default:
		lineDocs := make([]PrettyDoc, 0, len(lines))
		for i := 0; i < len(lines); i++ {
			lineDocs = append(lineDocs, lineDoc{line: lines[i]})
		}
		return sequenceDoc{items: lineDocs}
	}
}

func Indent(indent uint, doc PrettyDoc) PrettyDoc {
	return indentDoc{indent: indent, item: doc}
}

func Sequence(doc PrettyDoc, moreDocs ...PrettyDoc) PrettyDoc {
	if len(moreDocs) == 0 {
		return doc
	}
	return sequenceDoc{
		items: append([]PrettyDoc{doc}, moreDocs...),
	}
}

type lineDoc struct {
	line string
}

func (s lineDoc) sealed() {}

func (s lineDoc) ToLines(indent uint) []string {
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
	item   PrettyDoc
}

func (s indentDoc) sealed() {}

func (i indentDoc) ToLines(indent uint) []string {
	return i.item.ToLines(indent + i.indent)
}

type sequenceDoc struct {
	items []PrettyDoc
}

func (s sequenceDoc) sealed() {}

func (seq sequenceDoc) ToLines(indent uint) []string {
	lines := make([]string, 0, len(seq.items))
	for i := 0; i < len(seq.items); i++ {
		lines = append(lines, seq.items[i].ToLines(indent)...)
	}
	return lines
}
