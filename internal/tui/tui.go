// Package tui owns the interactive terminal screen, keyboard handling, and
// drawing. It consumes the layout package and never parses Markdown itself.
package tui

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/uniseg"

	"github.com/rumirifai/reamd/internal/cli"
	"github.com/rumirifai/reamd/internal/clipboard"
	"github.com/rumirifai/reamd/internal/layout"
)

// Minimum usable terminal size in cells.
const (
	MinWidth  = 40
	MinHeight = 10
)

// App is the screen-independent reader state.
type App struct {
	path              string
	theme             cli.Theme
	source            []string
	doc               *layout.Layout
	view              layout.Viewport
	width             int
	height            int
	help              bool
	zoom              int
	focused           int
	message           string
	opener            func(string) error
	clipboard         clipboard.Writer
	mouseEnabled      bool
	draggingScrollbar bool
	dragGrab          int
	mouseDown         bool
	mouseDownLink     string
	selection         textSelection
}

// New creates a reader for one document. The theme is validated by the CLI and
// is carried here for the palette work in a later ticket.
func New(path string, theme cli.Theme, source []string) *App {
	return NewWithOpener(path, theme, source, openBrowser)
}

// NewWithOpener allows browser activation to be verified without launching an OS browser.
func NewWithOpener(path string, theme cli.Theme, source []string, opener func(string) error) *App {
	return NewWithAdapters(path, theme, source, opener, clipboard.System{})
}

// NewWithAdapters injects operating-system services for deterministic tests.
func NewWithAdapters(path string, theme cli.Theme, source []string, opener func(string) error, clipboardWriter clipboard.Writer) *App {
	a := &App{path: path, theme: theme, source: source, zoom: 100, focused: -1, opener: opener, clipboard: clipboardWriter}
	a.Resize(80, 24)
	return a
}

func openBrowser(url string) error {
	var name string
	var args []string
	switch runtime.GOOS {
	case "windows":
		name, args = "rundll32", []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		name, args = "open", []string{url}
	default:
		name, args = "xdg-open", []string{url}
	}
	return exec.Command(name, args...).Start()
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
	layoutWidth := width * 100 / a.zoom
	if layoutWidth > width {
		layoutWidth = width
	}
	a.doc = layout.BuildTables(a.source, layoutWidth)
	if a.zoom < 100 {
		compact := make([]layout.Line, 0, len(a.doc.Lines))
		emptySeen := 0
		for _, line := range a.doc.Lines {
			if line.Text == "" {
				emptySeen++
				if emptySeen%((100-a.zoom)/10+1) != 0 {
					continue
				}
			} else {
				emptySeen = 0
			}
			compact = append(compact, line)
		}
		a.doc.Lines = compact
	}
	a.view.Resize(contentHeight, len(a.doc.Lines))
	if a.mouseEnabled && len(a.doc.Lines) > contentHeight {
		layoutWidth = (width - 1) * 100 / a.zoom
		if layoutWidth < 1 {
			layoutWidth = 1
		}
		a.doc = layout.BuildTables(a.source, layoutWidth)
		if a.zoom < 100 {
			compact := make([]layout.Line, 0, len(a.doc.Lines))
			emptySeen := 0
			for _, line := range a.doc.Lines {
				if line.Text == "" {
					emptySeen++
					if emptySeen%((100-a.zoom)/10+1) != 0 {
						continue
					}
				} else {
					emptySeen = 0
				}
				compact = append(compact, line)
			}
			a.doc.Lines = compact
		}
		a.view.Resize(contentHeight, len(a.doc.Lines))
	}
}

// Zoom reports the application layout density. It does not change terminal font size.
func (a *App) Zoom() int { return a.zoom }

func (a *App) setZoom(zoom int) {
	if zoom < 70 {
		zoom = 70
	}
	if zoom > 150 {
		zoom = 150
	}
	if zoom == a.zoom {
		return
	}
	oldPercent := a.view.Percent()
	a.zoom = zoom
	a.Resize(a.width, a.height)
	maxOffset := a.view.Total - a.view.Height
	if maxOffset > 0 {
		a.view.Offset = maxOffset * oldPercent / 100
	}
	a.view.Clamp()
}

// HandleKey applies one key event. It reports whether the app should exit.
func (a *App) HandleKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyCtrlC || (ev.Key() == tcell.KeyRune && ev.Modifiers()&tcell.ModCtrl != 0 && ev.Rune() == 'c') {
		a.copySelection()
		return false
	}
	if ev.Key() == tcell.KeyTab || ev.Key() == tcell.KeyBacktab {
		a.moveFocus(ev.Key() == tcell.KeyBacktab || ev.Modifiers()&tcell.ModShift != 0)
		return false
	}
	if ev.Key() == tcell.KeyEnter {
		a.activateFocused()
		return false
	}
	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case '=':
			a.setZoom(a.zoom + 10)
			return false
		case '+':
			a.setZoom(a.zoom + 10)
			return false
		case '-':
			a.setZoom(a.zoom - 10)
			return false
		case '0':
			a.setZoom(100)
			return false
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
	screen.SetStyle(a.baseStyle())
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

func (a *App) baseStyle() tcell.Style {
	if os.Getenv("NO_COLOR") != "" {
		return tcell.StyleDefault
	}
	switch a.theme {
	case cli.ThemeDark:
		return tcell.StyleDefault.Foreground(tcell.ColorLightGray).Background(tcell.ColorBlack)
	case cli.ThemeLight:
		return tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite)
	default:
		return tcell.StyleDefault
	}
}

func (a *App) drawDocument(screen tcell.Screen) {
	if len(a.doc.Lines) == 0 {
		screen.PutStrStyled(0, 0, "(empty file)", a.baseStyle().Dim(true))
		return
	}
	left := (a.width - a.doc.Width) / 2
	index := a.focusedLinkIndex()
	for row, line := range a.view.Visible(a.doc.Lines) {
		style := a.styleForLine(line.Text)
		screen.PutStrStyled(left, row, line.Text, style)
		docRow := a.view.Offset + row
		x := left
		graphemes := uniseg.NewGraphemes(line.Text)
		for graphemes.Next() {
			width := graphemes.Width()
			if a.selection.contains(docRow, x-left, width) {
				cluster := []rune(graphemes.Str())
				if len(cluster) > 0 {
					screen.SetContent(x, row, cluster[0], cluster[1:], style.Reverse(true))
				}
			}
			x += width
		}
		for _, link := range line.Links {
			linkStyle := style.Underline(true)
			if index >= 0 && a.linkAt(index, a.view.Offset+row, link) {
				linkStyle = linkStyle.Reverse(true)
			}
			x := left
			graphemes := uniseg.NewGraphemes(line.Text)
			for graphemes.Next() {
				cluster := []rune(graphemes.Str())
				if x < left+link.End && x+graphemes.Width() > left+link.Start && len(cluster) > 0 {
					cellStyle := linkStyle
					if a.selection.contains(docRow, x-left, graphemes.Width()) {
						cellStyle = cellStyle.Reverse(true)
					}
					screen.SetContent(x, row, cluster[0], cluster[1:], cellStyle)
				}
				x += graphemes.Width()
			}
		}
	}
	a.drawScrollbar(screen)
	if a.message != "" {
		screen.PutStrStyled(0, a.height-2, truncate(a.message, a.width), a.baseStyle().Bold(true))
	}
}

func (a *App) links() []struct {
	row  int
	link layout.Link
} {
	var links []struct {
		row  int
		link layout.Link
	}
	for row, line := range a.doc.Lines {
		for _, link := range line.Links {
			links = append(links, struct {
				row  int
				link layout.Link
			}{row, link})
		}
	}
	return links
}
func (a *App) focusedLinkIndex() int { return a.focused }
func (a *App) linkAt(index, row int, link layout.Link) bool {
	all := a.links()
	return index >= 0 && index < len(all) && all[index].row == row && all[index].link == link
}
func (a *App) moveFocus(reverse bool) {
	links := a.links()
	if len(links) == 0 {
		a.focused = -1
		return
	}
	if reverse {
		if a.focused <= 0 {
			a.focused = len(links) - 1
		} else {
			a.focused--
		}
	} else {
		a.focused = (a.focused + 1) % len(links)
	}
	row := links[a.focused].row
	if row < a.view.Offset {
		a.view.Offset = row
	}
	if row >= a.view.Offset+a.view.Height {
		a.view.Offset = row - a.view.Height + 1
	}
	a.view.Clamp()
	a.message = ""
}
func (a *App) activateFocused() {
	links := a.links()
	if a.focused < 0 || a.focused >= len(links) {
		return
	}
	a.activate(links[a.focused].link.Target)
}
func (a *App) activate(target string) {
	if strings.HasPrefix(target, "#") {
		a.jumpAnchor(strings.TrimPrefix(target, "#"))
		return
	}
	parsed, err := url.Parse(target)
	if err != nil || (strings.ToLower(parsed.Scheme) != "http" && strings.ToLower(parsed.Scheme) != "https") || parsed.Host == "" {
		return
	}
	if err := a.opener(target); err != nil {
		a.message = "Cannot open browser: " + err.Error()
	} else {
		a.message = ""
	}
}

func (a *App) HandleMouse(ev *tcell.EventMouse) {
	if !a.mouseEnabled {
		return
	}
	if a.help {
		a.draggingScrollbar = false
		return
	}
	x, y := ev.Position()
	buttons := ev.Buttons()
	if a.draggingScrollbar {
		if buttons&tcell.Button1 != 0 {
			trackHeight := a.height - 1
			_, thumbSize := a.view.ScrollbarThumb(trackHeight)
			a.view.ScrollToTrack(y-a.dragGrab+thumbSize/2, trackHeight)
			return
		}
		a.draggingScrollbar = false
	}
	if a.mouseDown {
		if buttons&tcell.Button1 != 0 {
			if pos, ok := a.documentPosition(x, y); ok && pos != a.selection.anchor {
				a.selection.active = pos
				a.selection.set = true
				a.message = ""
			}
			return
		}
		if !a.selection.set && a.mouseDownLink != "" && a.linkTargetAt(x, y) == a.mouseDownLink {
			a.activate(a.mouseDownLink)
		}
		a.mouseDown = false
		a.mouseDownLink = ""
		return
	}
	if buttons&tcell.WheelUp != 0 {
		a.view.ScrollLines(-3)
		return
	}
	if buttons&tcell.WheelDown != 0 {
		a.view.ScrollLines(3)
		return
	}
	if buttons&tcell.Button1 == 0 {
		return
	}
	trackHeight := a.height - 1
	if x == a.width-1 && len(a.doc.Lines) > a.view.Height && y >= 0 && y < trackHeight {
		thumbStart, thumbSize := a.view.ScrollbarThumb(trackHeight)
		if a.draggingScrollbar {
			a.view.ScrollToTrack(y-a.dragGrab+thumbSize/2, trackHeight)
			return
		}
		if y >= thumbStart && y < thumbStart+thumbSize {
			if !a.draggingScrollbar {
				a.dragGrab = y - thumbStart
			}
			a.draggingScrollbar = true
			a.view.ScrollToTrack(y-a.dragGrab+thumbSize/2, trackHeight)
		} else {
			a.draggingScrollbar = true
			a.dragGrab = thumbSize / 2
			a.view.ScrollToTrack(y, trackHeight)
		}
		return
	}
	pos, ok := a.documentPosition(x, y)
	if !ok {
		return
	}
	a.mouseDown = true
	a.selection = textSelection{anchor: pos, active: pos}
	a.mouseDownLink = a.linkTargetAt(x, y)
	a.message = ""
}

func (a *App) documentPosition(x, y int) (cellPosition, bool) {
	if y < 0 || y >= a.height-1 || a.help || a.width < MinWidth || a.height < MinHeight {
		return cellPosition{}, false
	}
	left := (a.width - a.doc.Width) / 2
	row := a.view.Offset + y
	if row < 0 || row >= len(a.doc.Lines) {
		return cellPosition{}, false
	}
	lineWidth := a.doc.Lines[row].Width
	if x < left || x >= left+lineWidth {
		return cellPosition{}, false
	}
	return cellPosition{row: row, col: x - left}, true
}

func (a *App) linkTargetAt(x, y int) string {
	left := (a.width - a.doc.Width) / 2
	row := a.view.Offset + y
	if row < 0 || row >= len(a.doc.Lines) {
		return ""
	}
	for _, link := range a.doc.Lines[row].Links {
		if x >= left+link.Start && x < left+link.End {
			return link.Target
		}
	}
	return ""
}

func (a *App) copySelection() {
	if !a.selection.set {
		a.message = "Select text with the mouse before copying."
		return
	}
	text := a.selection.plainText(a.doc.Lines)
	if text == "" {
		a.message = "Select text with the mouse before copying."
		return
	}
	if a.clipboard == nil {
		a.message = "Clipboard is unavailable; selection is preserved."
		return
	}
	if err := a.clipboard.WriteText(text); err != nil {
		a.message = "Cannot copy to clipboard: " + layout.Sanitize(err.Error())
		return
	}
	a.message = "Copied selection to clipboard."
}

func (a *App) drawScrollbar(screen tcell.Screen) {
	trackHeight := a.height - 1
	if !a.mouseEnabled || trackHeight < 1 || len(a.doc.Lines) <= a.view.Height {
		return
	}
	start, size := a.view.ScrollbarThumb(trackHeight)
	style := a.baseStyle().Dim(true)
	for y := 0; y < trackHeight; y++ {
		glyph := '│'
		if y >= start && y < start+size {
			glyph = '█'
		}
		screen.SetContent(a.width-1, y, glyph, nil, style)
	}
}
func (a *App) jumpAnchor(anchor string) {
	used := map[string]int{}
	for row, line := range a.doc.Lines {
		if !strings.HasPrefix(line.Text, "# ") && !strings.HasPrefix(line.Text, "## ") && !strings.HasPrefix(line.Text, "### ") {
			continue
		}
		title := strings.TrimLeft(line.Text, "# ")
		slug := githubSlug(title)
		n := used[slug]
		used[slug]++
		if n > 0 {
			slug += "-" + fmt.Sprint(n)
		}
		if slug == anchor {
			a.view.Offset = row
			a.view.Clamp()
			a.message = ""
			return
		}
	}
	a.message = "Anchor not found: #" + anchor
}
func githubSlug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		if unicode.IsSpace(r) {
			b.WriteByte('-')
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsMark(r) || unicode.Is(unicode.Pc, r) || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func (a *App) styleForLine(line string) tcell.Style {
	style := tcell.StyleDefault
	if strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ") {
		style = style.Bold(true).Underline(a.zoom >= 120)
	} else if strings.HasPrefix(line, "│ ") || strings.HasPrefix(line, "↪ ") {
		style = style.Underline(true).Reverse(a.zoom >= 130)
	} else if strings.HasPrefix(line, "> ") || strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "+ ") {
		style = style.Bold(true)
	}
	if _, noColor := os.LookupEnv("NO_COLOR"); noColor {
		return style
	}
	switch a.theme {
	case cli.ThemeDark:
		fg := tcell.ColorLightGray
		bg := tcell.ColorBlack
		if strings.HasPrefix(line, "#") {
			fg = tcell.ColorAqua
		}
		if strings.HasPrefix(line, "│ ") || strings.HasPrefix(line, "↪ ") {
			fg = tcell.ColorYellow
		}
		style = style.Foreground(fg).Background(bg)
	case cli.ThemeLight:
		fg := tcell.ColorBlack
		bg := tcell.ColorWhite
		if strings.HasPrefix(line, "#") {
			fg = tcell.ColorBlue
		}
		if strings.HasPrefix(line, "│ ") || strings.HasPrefix(line, "↪ ") {
			fg = tcell.ColorMaroon
		}
		style = style.Foreground(fg).Background(bg)
	}
	return style
}

func (a *App) drawHelp(screen tcell.Screen) {
	base := a.baseStyle()
	title := base.Bold(true)
	lines := []struct {
		key  string
		what string
	}{
		{"Up / Down, j / k", "Scroll one line"},
		{"PageUp / PageDown", "Scroll one screen"},
		{"Home / End", "Jump to start/end"},
		{"Tab / Shift+Tab / Enter", "Focus / activate link"},
		{"Mouse drag", "Select visible text"},
		{"Ctrl+C", "Copy selected text"},
		{"= / + / - / 0", "In / out / reset"},
		{"? / q / Esc", "Close help / quit"},
	}
	screen.PutStrStyled(0, 0, "Keys", title)
	for i, item := range lines {
		row := i + 1
		screen.PutStrStyled(2, row, item.key, base.Bold(true))
		screen.PutStrStyled(23, row, item.what, base)
	}
}

func (a *App) drawTooSmall(screen tcell.Screen) {
	base := a.baseStyle()
	screen.PutStrStyled(0, 0, "Terminal too small", base.Bold(true))
	screen.PutStrStyled(0, 2, fmt.Sprintf("Need %dx%d", MinWidth, MinHeight), base)
	screen.PutStrStyled(0, 3, "Resize, q quits", base)
}

func (a *App) drawFooter(screen tcell.Screen) {
	right := fmt.Sprintf("? =,+,-,0 %d%% q quit", a.zoom)
	if a.help {
		right = fmt.Sprintf("?close =,+,-,0 %d%% q quit", a.zoom)
	}
	screen.PutStrStyled(0, a.height-1, footerText(a.path, right, a.width), a.baseStyle())
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
	return runInitializedWithOpener(screen, path, theme, source, handleSignals, openBrowser)
}

func runInitializedWithOpener(screen tcell.Screen, path string, theme cli.Theme, source []string, handleSignals bool, opener func(string) error) error {
	return withRestore(screen, func() error {
		mouseEnabled := enableMouse(screen)
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

		app := NewWithOpener(path, theme, source, opener)
		app.mouseEnabled = mouseEnabled
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
			case *tcell.EventMouse:
				app.HandleMouse(ev)
			case *tcell.EventInterrupt:
				return nil
			}
			app.Draw(screen)
			screen.Show()
		}
	})
}

func enableMouse(screen tcell.Screen) bool {
	if !screen.HasMouse() {
		return false
	}
	screen.EnableMouse(tcell.MouseDragEvents)
	return true
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
