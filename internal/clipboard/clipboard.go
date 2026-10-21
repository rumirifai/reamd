// Package clipboard provides the system clipboard adapter used by the TUI.
package clipboard

import atottoclipboard "github.com/atotto/clipboard"

// Writer is the TUI-facing plain-text clipboard contract.
type Writer interface {
	WriteText(string) error
}

// System writes plain text to the operating system clipboard.
type System struct{}

func (System) WriteText(text string) error {
	return atottoclipboard.WriteAll(text)
}
