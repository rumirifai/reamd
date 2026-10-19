package tui

import (
	"testing"

	"github.com/gdamore/tcell/v2"

	"github.com/rumirifai/reamd/internal/cli"
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
