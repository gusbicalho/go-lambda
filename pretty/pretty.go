package pretty

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Pretty[context any] interface {
	ToPrettyDoc(context context) PrettyDoc
}

type PrettyDoc struct{ impl prettyDocImpl }

func (d PrettyDoc) String() string {
	return strings.Join(d.toLines(writeNothing, writeNothing), "\n")
}

func (d PrettyDoc) toLines(writePrefix func(*strings.Builder), writeSuffix func(*strings.Builder)) []string {
	return d.impl.toLines(writePrefix, writeSuffix)
}

type prettyDocImpl interface {
	toLines(writePrefix func(*strings.Builder), writeSuffix func(*strings.Builder)) []string
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

func (doc PrettyDoc) SuffixLines(suffixes []string) PrettyDoc {
	if len(suffixes) == 0 {
		return doc
	}
	suffixes = slices.Clone(suffixes)
	return PrettyDoc{lineSuffixDoc{suffixes: suffixes, item: doc}}
}

func PrefixLines(prefixes []string, doc PrettyDoc) PrettyDoc {
	prefixes = padPrefixes(prefixes)
	if len(prefixes) == 0 {
		return doc
	}
	return PrettyDoc{linePrefixDoc{prefixes: prefixes, item: doc}}
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

type lineSuffixDoc struct {
	suffixes []string
	item     prettyDocImpl
}

func (d lineSuffixDoc) toLines(writePrefix func(*strings.Builder), writeSuffix func(*strings.Builder)) []string {
	line := 0
	lastSuffix := d.suffixes[len(d.suffixes)-1]
	return d.item.toLines(
		writePrefix,
		func(builder *strings.Builder) {
			writePrefix(builder)
			if line < len(d.suffixes) {
				builder.WriteString(d.suffixes[line])
			} else {
				builder.WriteString(lastSuffix)
			}
			line++
		},
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
