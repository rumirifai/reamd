// Package document holds the text model shown by the TUI. Markdown parsing
// arrives in a later ticket; this vertical slice keeps the source lines intact.
package document

import "strings"

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
	content = strings.TrimSuffix(content, "\n")
	if content == "" {
		return &Document{Path: path}
	}
	return &Document{Path: path, Lines: strings.Split(content, "\n")}
}
