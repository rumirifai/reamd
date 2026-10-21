package tui

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"

	"github.com/rumirifai/reamd/internal/cli"
	"github.com/rumirifai/reamd/internal/layout"
)

func TestNewWrapsAndShowsFirstLines(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, []string{"first", "second", "third"})
	if got := len(app.doc.Lines); got != 3 {
		t.Fatalf("doc lines = %d, want 3", got)
	}
	if app.view.Offset != 0 {
		t.Fatalf("initial offset = %d, want 0", app.view.Offset)
	}
}

func TestDrawSanitizesActualTerminalControlBytes(t *testing.T) {
	source := []string{"before\x1b[2Jafter \x07 \x1b]52;c;DATA\x07 \u0085\u009b31m"}
	app := New("controls.md", cli.ThemeAuto, source)
	app.Resize(80, 10)
	screen := newScreen(t, 80, 10)
	app.Draw(screen)
	screen.Show()
	cells, width, height := screen.GetContents()
	var rendered strings.Builder
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			for _, r := range cells[y*width+x].Runes {
				if (r < 0x20 && r != '\n') || (r >= 0x7f && r <= 0x9f) {
					t.Fatalf("TUI cell retained control U+%04X", r)
				}
				rendered.WriteRune(r)
			}
		}
	}
	for _, want := range []string{"^[[2J", "^G", "^[]52;c;DATA^G", `\u0085`, `\u009b`} {
		if !strings.Contains(rendered.String(), want) {
			t.Errorf("TUI output lacks visible notation %q: %q", want, rendered.String())
		}
	}
}

func TestHandleKeyScrollsAndQuits(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, lineRange(50))
	app.Resize(40, 10)

	app.HandleKey(key(tcell.KeyDown, 0))
	if app.view.Offset != 1 {
		t.Fatalf("offset after Down = %d, want 1", app.view.Offset)
	}
	app.HandleKey(runeKey('j'))
	if app.view.Offset != 2 {
		t.Fatalf("offset after j = %d, want 2", app.view.Offset)
	}
	app.HandleKey(runeKey('k'))
	if app.view.Offset != 1 {
		t.Fatalf("offset after k = %d, want 1", app.view.Offset)
	}
	app.HandleKey(key(tcell.KeyPgDn, 0))
	if app.view.Offset != 9 {
		t.Fatalf("offset after PageDown = %d, want 9", app.view.Offset)
	}
	app.HandleKey(key(tcell.KeyEnd, 0))
	if app.view.Offset != 41 {
		t.Fatalf("offset after End = %d, want 41", app.view.Offset)
	}
	app.HandleKey(key(tcell.KeyHome, 0))
	if app.view.Offset != 0 {
		t.Fatalf("offset after Home = %d, want 0", app.view.Offset)
	}
	if quit := app.HandleKey(key(tcell.KeyEsc, 0)); !quit {
		t.Fatal("Esc did not request quit")
	}
	if quit := app.HandleKey(runeKey('q')); !quit {
		t.Fatal("q did not request quit")
	}
}

func TestHelpToggles(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, lineRange(50))
	app.Resize(40, 10)
	app.HandleKey(runeKey('?'))
	if !app.help {
		t.Fatal("? did not open help")
	}
	app.HandleKey(runeKey('j'))
	if app.view.Offset != 0 {
		t.Fatalf("offset changed while help was open: %d", app.view.Offset)
	}
	app.HandleKey(runeKey('?'))
	if app.help {
		t.Fatal("? did not close help")
	}
}

func TestResizeReclampsOffset(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, lineRange(20))
	app.Resize(40, 10)
	app.HandleKey(key(tcell.KeyEnd, 0))
	app.Resize(40, 30)
	if app.view.Offset != 0 {
		t.Fatalf("offset after growing terminal = %d, want 0", app.view.Offset)
	}
}

func TestThemeStyles(t *testing.T) {
	for _, theme := range []cli.Theme{cli.ThemeAuto, cli.ThemeDark, cli.ThemeLight, cli.ThemeTransparent} {
		t.Run(string(theme), func(t *testing.T) {
			t.Setenv("NO_COLOR", "")
			app := New("doc.md", theme, []string{"# Heading"})
			style := app.styleForLine("# Heading")
			fg, bg, attrs := style.Decompose()
			if attrs&tcell.AttrBold == 0 {
				t.Fatal("heading lacks structural bold emphasis")
			}
			if theme == cli.ThemeAuto && (fg != tcell.ColorDefault || bg != tcell.ColorDefault) {
				t.Fatalf("auto theme overrides terminal colors: fg=%v bg=%v", fg, bg)
			}
			baseFG, baseBG, _ := app.baseStyle().Decompose()
			if theme == cli.ThemeDark && (baseFG != tcell.ColorLightGray || baseBG != tcell.ColorBlack) {
				t.Fatalf("dark palette = %v on %v", baseFG, baseBG)
			}
			if theme == cli.ThemeLight && (baseFG != tcell.ColorBlack || baseBG != tcell.ColorWhite) {
				t.Fatalf("light palette = %v on %v", baseFG, baseBG)
			}
			if theme == cli.ThemeTransparent && bg != tcell.ColorDefault {
				t.Fatalf("transparent background = %v, want terminal default", bg)
			}
		})
	}
}

func TestNoColorRetainsStructuralHierarchy(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	app := New("doc.md", cli.ThemeDark, []string{"# Heading", "- item", "│ code"})
	for _, line := range []string{"# Heading", "- item", "│ code"} {
		fg, bg, attrs := app.styleForLine(line).Decompose()
		if fg != tcell.ColorDefault || bg != tcell.ColorDefault {
			t.Fatalf("NO_COLOR style for %q contains colors: %v %v", line, fg, bg)
		}
		if attrs == 0 {
			t.Fatalf("NO_COLOR style for %q has no structural marker", line)
		}
	}
}

func TestZoomBoundsStepResetAndReflow(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, []string{"a long line that should wrap differently as zoom increases"})
	app.Resize(40, 10)
	app.HandleKey(runeKey('+'))
	if app.Zoom() != 110 || app.doc.Width != 36 {
		t.Fatalf("zoom +: zoom=%d width=%d", app.Zoom(), app.doc.Width)
	}
	app.HandleKey(tcell.NewEventKey(tcell.KeyRune, '0', tcell.ModCtrl))
	if app.Zoom() != 100 {
		t.Fatalf("Ctrl+0 zoom = %d", app.Zoom())
	}
	for i := 0; i < 5; i++ {
		app.HandleKey(runeKey('-'))
	}
	if app.Zoom() != 70 {
		t.Fatalf("lower bound zoom = %d", app.Zoom())
	}
	app.HandleKey(runeKey('-'))
	if app.Zoom() != 70 {
		t.Fatalf("zoom passed lower bound: %d", app.Zoom())
	}
	for i := 0; i < 8; i++ {
		app.HandleKey(runeKey('+'))
	}
	if app.Zoom() != 150 {
		t.Fatalf("upper bound zoom = %d", app.Zoom())
	}
	app.HandleKey(runeKey('+'))
	if app.Zoom() != 150 {
		t.Fatalf("zoom passed upper bound: %d", app.Zoom())
	}
}

func TestResizeReflowsResponsiveTable(t *testing.T) {
	source := []string{"\x1eTABLE", "\x1eALIGN:left\x1fleft", "\x1eROW:Name\x1fStatus", "\x1eROW:Renderer\x1fReady", "\x1eENDTABLE"}
	app := New("table.md", cli.ThemeAuto, source)
	app.Resize(40, 10)
	if app.doc.Width != 40 || !strings.Contains(joinLines(app.doc.Lines), "| Name") {
		t.Fatalf("wide table did not use grid layout: %#v", app.doc.Lines)
	}
	app.Resize(20, 10)
	if app.doc.Width != 20 || !strings.Contains(joinLines(app.doc.Lines), "Name: Renderer") {
		t.Fatalf("narrow table did not switch to label/value layout: %#v", app.doc.Lines)
	}
}

func joinLines(lines []layout.Line) string {
	var text []string
	for _, line := range lines {
		text = append(text, line.Text)
	}
	return strings.Join(text, "\n")
}

func TestDrawMinimumSizeMessage(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, []string{"content"})
	app.Resize(39, 9)
	screen := newScreen(t, 39, 9)
	app.Draw(screen)
	screen.Show()
	if !screenContains(screen, "Terminal too small") {
		t.Fatal("small terminal message not drawn")
	}
	if !screenContains(screen, "Resize, q quits") {
		t.Fatal("quit hint not drawn on small terminal")
	}
}

func TestDrawDocumentAndFooter(t *testing.T) {
	app := New("path with space.md", cli.ThemeAuto, []string{"hello world", "second line"})
	app.Resize(40, 10)
	screen := newScreen(t, 40, 10)
	app.Draw(screen)
	screen.Show()
	if !screenContains(screen, "hello world") {
		t.Fatal("document content not drawn")
	}
	if !screenContains(screen, "path with space.md") {
		t.Fatal("footer path not drawn")
	}
	if !screenContains(screen, "q quit") {
		t.Fatal("footer key hints not drawn")
	}
}

func TestDrawEmptyDocument(t *testing.T) {
	app := New("empty.md", cli.ThemeAuto, nil)
	app.Resize(40, 10)
	screen := newScreen(t, 40, 10)
	app.Draw(screen)
	screen.Show()
	if !screenContains(screen, "(empty file)") {
		t.Fatal("empty state not drawn")
	}
}

func TestDrawHelp(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, lineRange(5))
	app.Resize(40, 10)
	app.HandleKey(runeKey('?'))
	screen := newScreen(t, 40, 10)
	app.Draw(screen)
	screen.Show()
	if !screenContains(screen, "PageUp / PageDown") {
		t.Fatal("help panel not drawn")
	}
}

func TestFooterTextFitsWidth(t *testing.T) {
	for _, width := range []int{40, 80, 120} {
		got := footerText("very/long/path/that/does/not/fit/at/all/document.md", "? help  q quit  100%", width)
		if runeWidth(got) > width {
			t.Fatalf("footerText width = %d, want <= %d (%q)", runeWidth(got), width, got)
		}
	}
}

func key(k tcell.Key, r rune) *tcell.EventKey {
	return tcell.NewEventKey(k, r, tcell.ModNone)
}

func runeKey(r rune) *tcell.EventKey {
	return tcell.NewEventKey(tcell.KeyRune, r, tcell.ModNone)
}

func lineRange(n int) []string {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = "line"
	}
	return lines
}
