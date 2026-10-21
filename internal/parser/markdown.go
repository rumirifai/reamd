// Package parser adapts Goldmark's Markdown AST for the document model.
package parser

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"
)

// Parsed retains both the AST and its source bytes for literal AST segments.
type Parsed struct {
	Root   ast.Node
	Source []byte
}

// Parse parses CommonMark without enabling out-of-scope extensions.
func Parse(source string) Parsed {
	data := []byte(source)
	return Parsed{Root: goldmark.New(goldmark.WithExtensions(extension.GFM)).Parser().Parse(text.NewReader(data)), Source: data}
}
