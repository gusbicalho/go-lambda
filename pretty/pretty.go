package pretty

import (
	"strings"
	"unicode/utf8"
)

type Pretty[context any] interface {
	ToPrettyDoc(context context) Doc
}

type Doc struct{ impl prettyDocImpl }

func (doc Doc) String() string {
	return strings.Join(doc.toLines(writeNothing, writeNothing), "\n")
}

func (doc Doc) toLines(writePrefix func(*strings.Builder), writeSuffix func(*strings.Builder)) []string {
	return doc.impl.toLines(writePrefix, writeSuffix)
}

type prettyDocImpl interface {
	toLines(writePrefix func(*strings.Builder), writeSuffix func(*strings.Builder)) []string
}

func FromString(s string) Doc {
	lines := strings.Split(s, "\n")
	switch len(lines) {
	case 0:
		return Doc{lineDoc{line: ""}}
	case 1:
		return Doc{lineDoc{line: lines[0]}}
	default:
		lineDocs := make([]prettyDocImpl, 0, len(lines))
		for i := 0; i < len(lines); i++ {
			lineDocs = append(lineDocs, lineDoc{line: lines[i]})
		}
		return Doc{sequenceDoc{items: lineDocs}}
	}
}

func Indent(indent uint, doc Doc) Doc {
	return Doc{indentDoc{indent: indent, item: doc}}
}

func PrefixLines(prefixes []string, doc Doc) Doc {
	prefixes = padPrefixes(prefixes)
	if len(prefixes) == 0 {
		return doc
	}
	return Doc{linePrefixDoc{prefixes: prefixes, item: doc}}
}

func padPrefixes(prefixes []string) []string {
	if len(prefixes) == 0 {
		return nil
	}
	maxLen := 0
	for _, prefix := range prefixes {
		if prefixLen := utf8.RuneCountInString(prefix); prefixLen > maxLen {
			maxLen = prefixLen
		}
	}
	pad := func(prefix string) string {
		prefixLen := utf8.RuneCountInString(prefix)
		if prefixLen == maxLen {
			return prefix
		}

		builder := strings.Builder{}
		builder.WriteString(prefix)

		for runesToAdd := maxLen - prefixLen; runesToAdd > 0; runesToAdd-- {
			builder.WriteString(" ")
		}
		return builder.String()
	}
	prefixesCopy := make([]string, 0, len(prefixes))
	for _, prefix := range prefixes {
		prefixesCopy = append(prefixesCopy, pad(prefix))
	}
	return prefixesCopy
}

func Sequence(doc Doc, moreDocs ...Doc) Doc {
	if len(moreDocs) == 0 {
		return doc
	}
	items := make([]prettyDocImpl, 1+len(moreDocs))
	items[0] = doc.impl
	for i, doc := range moreDocs {
		items[i+1] = doc.impl
	}
	return Doc{sequenceDoc{items: items}}
}

type lineDoc struct {
	line string
}

func (s lineDoc) toLines(writePrefix func(*strings.Builder), writeSuffix func(*strings.Builder)) []string {
	builder := strings.Builder{}
	writePrefix(&builder)
	builder.WriteString(s.line)
	writeSuffix(&builder)
	return []string{builder.String()}
}

type indentDoc struct {
	indent uint
	item   prettyDocImpl
}

func (i indentDoc) toLines(writePrefix func(*strings.Builder), writeSuffix func(*strings.Builder)) []string {
	return i.item.toLines(func(builder *strings.Builder) {
		writePrefix(builder)
		writeIndent(builder, i.indent)
	}, writeSuffix)
}

type linePrefixDoc struct {
	prefixes []string
	item     prettyDocImpl
}

func (d linePrefixDoc) toLines(writePrefix func(*strings.Builder), writeSuffix func(*strings.Builder)) []string {
	line := 0
	lastPrefix := d.prefixes[len(d.prefixes)-1]
	return d.item.toLines(
		func(builder *strings.Builder) {
			writePrefix(builder)
			if line < len(d.prefixes) {
				builder.WriteString(d.prefixes[line])
			} else {
				builder.WriteString(lastPrefix)
			}
			line++
		},
		writeSuffix,
	)
}

type sequenceDoc struct {
	items []prettyDocImpl
}

func (seq sequenceDoc) toLines(writePrefix func(*strings.Builder), writeSuffix func(*strings.Builder)) []string {
	lines := make([]string, 0, len(seq.items))
	for i := 0; i < len(seq.items); i++ {
		lines = append(lines, seq.items[i].toLines(writePrefix, writeSuffix)...)
	}
	return lines
}

func writeNothing(_ *strings.Builder) {}

func writeIndent(builder *strings.Builder, indent uint) {
	builder.Grow(int(indent))
	for i := uint(0); i < indent; i++ {
		builder.WriteString(" ")
	}
}
