// Package app wires the CLI, file loader, document model, and TUI together.
package app

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/rumirifai/reamd/internal/cli"
	"github.com/rumirifai/reamd/internal/document"
	"github.com/rumirifai/reamd/internal/layout"
	"github.com/rumirifai/reamd/internal/source"
	"github.com/rumirifai/reamd/internal/tui"
)

type tuiRunner func(path string, theme cli.Theme, lines []string) error

// Run executes reamd with the real standard streams and terminal detection.
func Run(args []string, stdout, stderr io.Writer) int {
	interactive := term.IsTerminal(int(os.Stdout.Fd())) && term.IsTerminal(int(os.Stdin.Fd()))
	return run(args, stdout, stderr, interactive, tui.Run)
}

func run(args []string, stdout, stderr io.Writer, interactive bool, runTUI tuiRunner) int {
	opts, err := cli.Parse(args)
	if err != nil {
		fmt.Fprintf(stderr, "reamd: %v\n\n", err)
		fmt.Fprint(stderr, cli.HelpText())
		return 2
	}
	if opts.Help {
		fmt.Fprint(stdout, cli.HelpText())
		return 0
	}
	if opts.Version {
		fmt.Fprint(stdout, cli.VersionText())
		return 0
	}

	content, err := source.Load(opts.Path)
	if err != nil {
		fmt.Fprintf(stderr, "reamd: %v\n", err)
		return 1
	}
	doc := document.New(opts.Path, content)

	if !interactive {
		writePlain(stdout, doc.Lines)
		return 0
	}

	if err := runTUI(doc.Path, opts.Theme, doc.Lines); err != nil {
		fmt.Fprintf(stderr, "reamd: cannot start the terminal UI: %v\n", err)
		return 1
	}
	return 0
}

// writePlain emits the parsed, laid-out document without ANSI or terminal takeover.
func writePlain(w io.Writer, lines []string) {
	if len(lines) == 0 {
		return
	}
	var b strings.Builder
	for _, line := range lines {
		b.WriteString(layout.Sanitize(line))
		b.WriteByte('\n')
	}
	io.WriteString(w, b.String())
}
