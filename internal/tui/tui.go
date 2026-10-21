// Package tui owns the interactive terminal screen, keyboard handling, and
// drawing. It consumes the layout package and never parses Markdown itself.
package tui

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/uniseg"

	"github.com/rumirifai/reamd/internal/cli"
	"github.com/rumirifai/reamd/internal/layout"
)

// Minimum usable terminal size in cells.
const (
	MinWidth  = 40
	MinHeight = 10
)

// App is the screen-independent reader state.
type App struct {
	path   string
	theme  cli.Theme
	source []string
	doc    *layout.Layout
	view   layout.Viewport
	width  int
	height int
	help   bool
}

// New creates a reader for one document. The theme is validated by the CLI and
// is carried here for the palette work in a later ticket.
func New(path string, theme cli.Theme, source []string) *App {
	a := &App{path: path, theme: theme, source: source}
	a.Resize(80, 24)
	return a
}

// Theme reports the validated theme carried by this reader.
func (a *App) Theme() cli.Theme {
	return a.theme
}

// Resize re-wraps the document and clamps the scroll position.
func (a *App) Resize(width, height int) {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	a.width, a.height = width, height
	contentHeight := height - 1
	if contentHeight < 1 {
		contentHeight = 1
	}
	a.doc = layout.BuildTables(a.source, width)
	a.view.Resize(contentHeight, len(a.doc.Lines))
}

// HandleKey applies one key event. It reports whether the app should exit.
func (a *App) HandleKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'q':
			return true
		case '?':
			a.help = !a.help
			return false
		case 'j':
			if !a.help {
				a.view.ScrollLines(1)
			}
			return false
		case 'k':
			if !a.help {
				a.view.ScrollLines(-1)
			}
			return false
		}
	}

	switch ev.Key() {
	case tcell.KeyEsc:
		return true
	case tcell.KeyUp:
		if !a.help {
			a.view.ScrollLines(-1)
		}
	case tcell.KeyDown:
		if !a.help {
			a.view.ScrollLines(1)
		}
	case tcell.KeyPgUp:
		if !a.help {
			a.view.ScrollPage(-1)
		}
	case tcell.KeyPgDn:
		if !a.help {
			a.view.ScrollPage(1)
		}
	case tcell.KeyHome:
		if !a.help {
			a.view.Home()
		}
	case tcell.KeyEnd:
		if !a.help {
			a.view.End()
		}
	}
	return false
}

// Draw paints the current frame. It never writes document content as terminal
// control sequences because the layout has already sanitized it.
func (a *App) Draw(screen tcell.Screen) {
	screen.Clear()
	if a.width < MinWidth || a.height < MinHeight {
		a.drawTooSmall(screen)
		a.drawFooter(screen)
		return
	}
	if a.help {
		a.drawHelp(screen)
	} else {
		a.drawDocument(screen)
	}
	a.drawFooter(screen)
}

func (a *App) drawDocument(screen tcell.Screen) {
	if len(a.doc.Lines) == 0 {
		screen.PutStrStyled(0, 0, "(empty file)", tcell.StyleDefault.Dim(true))
		return
	}
	for row, line := range a.view.Visible(a.doc.Lines) {
		screen.PutStrStyled(0, row, line.Text, tcell.StyleDefault)
	}
}

func (a *App) drawHelp(screen tcell.Screen) {
	title := tcell.StyleDefault.Bold(true)
	lines := []struct {
		key  string
		what string
	}{
		{"Up / Down, j / k", "Scroll one line"},
		{"PageUp / PageDown", "Scroll one screen"},
		{"Home / End", "Jump to start/end"},
		{"?", "Close this help"},
		{"q / Esc", "Quit"},
	}
	screen.PutStrStyled(0, 0, "Keys", title)
	for i, item := range lines {
		row := i + 2
		screen.PutStrStyled(2, row, item.key, tcell.StyleDefault.Bold(true))
		screen.PutStrStyled(23, row, item.what, tcell.StyleDefault)
	}
}

func (a *App) drawTooSmall(screen tcell.Screen) {
	screen.PutStrStyled(0, 0, "Terminal too small", tcell.StyleDefault.Bold(true))
	screen.PutStrStyled(0, 2, fmt.Sprintf("Need %dx%d", MinWidth, MinHeight), tcell.StyleDefault)
	screen.PutStrStyled(0, 3, "Resize, q quits", tcell.StyleDefault)
}

func (a *App) drawFooter(screen tcell.Screen) {
	right := fmt.Sprintf("? help  q quit  %d%%", a.view.Percent())
	if a.help {
		right = fmt.Sprintf("? close  q quit  %d%%", a.view.Percent())
	}
	screen.PutStrStyled(0, a.height-1, footerText(a.path, right, a.width), tcell.StyleDefault)
}

// footerText places the path on the left and the key hints on the right,
// shortening the path from the left when space runs out.
func footerText(path, right string, width int) string {
	rightWidth := uniseg.StringWidth(right)
	if rightWidth >= width {
		return truncate(right, width)
	}
	leftWidth := width - rightWidth - 1
	left := truncateLeft(path, leftWidth)
	gap := width - uniseg.StringWidth(left) - rightWidth
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + right
}

func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if uniseg.StringWidth(s) <= width {
		return s
	}
	var b strings.Builder
	used := 0
	g := uniseg.NewGraphemes(s)
	for g.Next() {
		w := g.Width()
		if used+w > width {
			break
		}
		b.WriteString(g.Str())
		used += w
	}
	return b.String()
}

func truncateLeft(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if uniseg.StringWidth(s) <= width {
		return s
	}
	const marker = "..."
	if width <= len(marker) {
		return truncate(s, width)
	}
	var clusters []string
	g := uniseg.NewGraphemes(s)
	for g.Next() {
		clusters = append(clusters, g.Str())
	}
	remaining := width - len(marker)
	used := 0
	start := len(clusters)
	for i := len(clusters) - 1; i >= 0; i-- {
		w := uniseg.StringWidth(clusters[i])
		if used+w > remaining {
			break
		}
		used += w
		start = i
	}
	return marker + strings.Join(clusters[start:], "")
}

// Run opens the real terminal and runs the reader until the user quits.
func Run(path string, theme cli.Theme, source []string) error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := screen.Init(); err != nil {
		// Init may have changed terminal modes before failing; Fini is
		// idempotent and restores whatever was entered.
		screen.Fini()
		return err
	}
	return runInitialized(screen, path, theme, source, true)
}

// runInitialized runs the event loop on an initialized screen and owns its
// restoration. handleSignals is false in tests.
func runInitialized(screen tcell.Screen, path string, theme cli.Theme, source []string, handleSignals bool) error {
	return withRestore(screen, func() error {
		if handleSignals {
			stop := make(chan struct{})
			defer close(stop)
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
			defer signal.Stop(sigCh)
			go func() {
				select {
				case <-sigCh:
					// PostEventWait is safe from this goroutine and guarantees
					// the event loop wakes to restore the terminal.
					screen.PostEventWait(tcell.NewEventInterrupt(nil))
				case <-stop:
				}
			}()
		}

		app := New(path, theme, source)
		app.Resize(screen.Size())
		app.Draw(screen)
		screen.Show()

		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case nil:
				return nil
			case *tcell.EventKey:
				if app.HandleKey(ev) {
					return nil
				}
			case *tcell.EventResize:
				app.Resize(ev.Size())
				screen.Sync()
			case *tcell.EventInterrupt:
				return nil
			}
			app.Draw(screen)
			screen.Show()
		}
	})
}

// withRestore always restores the terminal, including when fn panics, then
// re-panics so the failure is not hidden.
func withRestore(screen tcell.Screen, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			screen.Fini()
			panic(r)
		}
		screen.Fini()
	}()
	return fn()
}
