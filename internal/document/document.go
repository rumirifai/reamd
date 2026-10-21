// Package document parses Markdown into the read-only text model used by layout.
package document

import (
	"strconv"
	"strings"

	"github.com/rumirifai/reamd/internal/parser"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Document is a read-only text document.
type Document struct {
	Path  string
	Lines []string
}

// New normalizes line endings and a leading byte order mark. The source file is
// never modified.
func New(path, content string) *Document {
	content = strings.TrimPrefix(content, "\uFEFF")
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	if content == "" {
		return &Document{Path: path}
	}
	parsed := parser.Parse(content)
	var lines []string
	for node := parsed.Root.FirstChild(); node != nil; node = node.NextSibling() {
		renderBlock(node, parsed.Source, 0, &lines)
	}
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return &Document{Path: path, Lines: lines}
}

func renderBlock(n ast.Node, src []byte, depth int, out *[]string) {
	switch node := n.(type) {
	case *ast.Heading:
		var b strings.Builder
		b.WriteString(strings.Repeat("#", node.Level))
		b.WriteByte(' ')
		renderInline(node, src, &b)
		*out = append(*out, b.String(), "")
	case *ast.Paragraph, *ast.TextBlock:
		indent := strings.Repeat("  ", depth)
		var b strings.Builder
		b.WriteString(indent)
		renderInline(node, src, &b)
		*out = append(*out, b.String(), "")
	case *ast.List:
		itemNumber := node.Start
		for item := node.FirstChild(); item != nil; item = item.NextSibling() {
			marker := "- "
			if node.IsOrdered() {
				marker = strconv.Itoa(itemNumber) + ". "
			}
			itemNumber++
			start := len(*out)
			*out = append(*out, strings.Repeat("  ", depth)+marker)
			for child := item.FirstChild(); child != nil; child = child.NextSibling() {
				if _, ok := child.(*ast.Paragraph); ok {
					var b strings.Builder
					renderInline(child, src, &b)
					(*out)[start] += b.String()
				} else if _, ok := child.(*ast.TextBlock); ok {
					var b strings.Builder
					renderInline(child, src, &b)
					(*out)[start] += b.String()
				} else {
					renderBlock(child, src, depth+1, out)
				}
			}
		}
		*out = append(*out, "")
	case *ast.Blockquote:
		start := len(*out)
		for child := node.FirstChild(); child != nil; child = child.NextSibling() {
			renderBlock(child, src, depth+1, out)
		}
		for i := start; i < len(*out); i++ {
			if (*out)[i] != "" {
				(*out)[i] = strings.Repeat("  ", depth) + "> " + (*out)[i]
			}
		}
	case *ast.ThematicBreak:
		*out = append(*out, strings.Repeat("  ", depth)+"────────────────", "")
	case *ast.CodeBlock:
		renderLiteralLines(node.Lines(), src, depth, out)
	case *ast.FencedCodeBlock:
		renderLiteralLines(node.Lines(), src, depth, out)
	case *ast.HTMLBlock:
		renderLiteralLines(node.Lines(), src, depth, out)
	default:
		// Unknown block nodes retain all descendant text rather than disappearing.
		var b strings.Builder
		renderInline(node, src, &b)
		if b.Len() > 0 {
			*out = append(*out, strings.Repeat("  ", depth)+b.String(), "")
		}
	}
}

func renderLiteralLines(lines *text.Segments, src []byte, depth int, out *[]string) {
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		*out = append(*out, strings.Repeat("  ", depth)+string(segment.Value(src)))
	}
	*out = append(*out, "")
}

func renderInline(n ast.Node, src []byte, b *strings.Builder) {
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		switch x := child.(type) {
		case *ast.Text:
			b.Write(x.Segment.Value(src))
			if x.SoftLineBreak() {
				b.WriteByte(' ')
			}
			if x.HardLineBreak() {
				b.WriteByte('\n')
			}
		case *ast.CodeSpan:
			b.WriteByte('`')
			b.Write(x.Text(src))
			b.WriteByte('`')
		case *ast.Link:
			renderInline(x, src, b)
			if len(x.Destination) > 0 {
				b.WriteString(" (")
				b.Write(x.Destination)
				b.WriteByte(')')
			}
		case *ast.Image:
			renderInline(x, src, b)
			if len(x.Destination) > 0 {
				b.WriteString(" (")
				b.Write(x.Destination)
				b.WriteByte(')')
			}
		case *ast.RawHTML:
			for i := 0; i < x.Segments.Len(); i++ {
				segment := x.Segments.At(i)
				b.Write(segment.Value(src))
			}
		case *ast.Emphasis:
			marker := "*"
			if x.Level == 2 {
				marker = "**"
			}
			b.WriteString(marker)
			renderInline(x, src, b)
			b.WriteString(marker)
		default:
			renderInline(child, src, b)
		}
	}
}
