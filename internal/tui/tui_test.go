package tui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/rumirifai/reamd/internal/cli"
	"github.com/rumirifai/reamd/internal/clipboard"
	"github.com/rumirifai/reamd/internal/document"
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
	app.HandleKey(runeKey('='))
	if app.Zoom() != 110 || app.doc.Width != 36 {
		t.Fatalf("zoom =: zoom=%d width=%d", app.Zoom(), app.doc.Width)
	}
	app.HandleKey(runeKey('+'))
	if app.Zoom() != 120 || app.doc.Width != 33 {
		t.Fatalf("zoom +: zoom=%d width=%d", app.Zoom(), app.doc.Width)
	}
	app.HandleKey(runeKey('-'))
	if app.Zoom() != 110 {
		t.Fatalf("zoom -: zoom=%d, want 110", app.Zoom())
	}
	app.HandleKey(runeKey('0'))
	if app.Zoom() != 100 {
		t.Fatalf("0 zoom = %d", app.Zoom())
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

func TestLinkKeyboardAnchorProtocolGateAndOpenerFailure(t *testing.T) {
	source := []string{"## First!", "\x1chttps://example.com\x1cWeb\x1c (https://example.com)", "\x1c#second\x1cJump\x1c (#second)", "\x1cmailto:user@example.com\x1cMail\x1c (mailto:user@example.com)", "\x1c#missing\x1cMissing\x1c (#missing)", "## Second", "after 1", "after 2", "after 3", "after 4", "after 5", "after 6", "after 7", "after 8", "after 9", "after 10"}
	calls := []string{}
	app := NewWithOpener("links.md", cli.ThemeAuto, source, func(url string) error { calls = append(calls, url); return errors.New("browser unavailable") })
	app.Resize(60, 10)
	if len(calls) != 0 {
		t.Fatal("constructor activated a link")
	}
	app.HandleKey(key(tcell.KeyTab, 0))
	if app.focused != 0 || len(calls) != 0 {
		t.Fatalf("Tab focus=%d opener calls=%v", app.focused, calls)
	}
	app.HandleKey(key(tcell.KeyBacktab, 0))
	if app.focused != 3 {
		t.Fatalf("Shift+Tab focus=%d, want previous last link", app.focused)
	}
	app.HandleKey(key(tcell.KeyTab, 0))
	app.HandleKey(key(tcell.KeyEnter, 0))
	if len(calls) != 1 || calls[0] != "https://example.com" {
		t.Fatalf("web activation calls=%v", calls)
	}
	if !strings.Contains(app.message, "Cannot open browser") {
		t.Fatalf("missing opener error: %q", app.message)
	}
	errorScreen := newScreen(t, 60, 10)
	app.Draw(errorScreen)
	errorScreen.Show()
	if !screenContains(errorScreen, "Cannot open browser") {
		cells, width, height := errorScreen.GetContents()
		t.Fatalf("browser failure message was not visible on screen (cells=%d, size=%dx%d, message=%q)", len(cells), width, height, app.message)
	}
	app.HandleKey(key(tcell.KeyTab, 0))
	app.HandleKey(key(tcell.KeyEnter, 0))
	if app.view.Offset != 5 {
		t.Fatalf("anchor offset=%d, want heading row 5", app.view.Offset)
	}
	app.HandleKey(key(tcell.KeyTab, 0))
	app.HandleKey(key(tcell.KeyEnter, 0))
	if len(calls) != 1 {
		t.Fatalf("non-HTTP protocol reached opener: %v", calls)
	}
	app.HandleKey(key(tcell.KeyTab, 0))
	app.HandleKey(key(tcell.KeyEnter, 0))
	if !strings.Contains(app.message, "Anchor not found") {
		t.Fatalf("missing anchor error: %q", app.message)
	}
}

func TestSimulatedMouseActivatesOnlyClickedLink(t *testing.T) {
	var calls []string
	source := []string{"\x1chttps://example.com/page\x1cOpen\x1c (https://example.com/page)"}
	baseScreen := initializedScreen(t, 60, 10)
	screen := &mouseCapabilityScreen{Screen: baseScreen, supported: true}
	done := make(chan error, 1)
	go func() {
		done <- runInitializedWithOpener(screen, "links.md", cli.ThemeAuto, source, false, func(url string) error { calls = append(calls, url); return nil })
	}()
	baseScreen.InjectMouse(1, 0, tcell.Button1, tcell.ModNone)
	baseScreen.InjectMouse(1, 0, tcell.ButtonNone, tcell.ModNone)
	baseScreen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("TUI did not exit")
	}
	if len(calls) != 1 || calls[0] != "https://example.com/page" {
		t.Fatalf("mouse opener calls=%v", calls)
	}
	if !screen.enabled {
		t.Fatal("mouse-capable screen was not enabled for mouse events")
	}
}

func TestMouseWheelTrackAndThumbDrag(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, lineRange(100))
	app.Resize(40, 10)
	app.mouseEnabled = true
	app.HandleMouse(tcell.NewEventMouse(5, 0, tcell.WheelDown, tcell.ModNone))
	if app.view.Offset != 3 {
		t.Fatalf("offset after wheel down = %d, want 3", app.view.Offset)
	}
	app.HandleMouse(tcell.NewEventMouse(5, 0, tcell.WheelUp, tcell.ModNone))
	if app.view.Offset != 0 {
		t.Fatalf("offset after wheel up = %d, want 0", app.view.Offset)
	}
	app.HandleMouse(tcell.NewEventMouse(39, 8, tcell.Button1, tcell.ModNone))
	if app.view.Offset < 85 {
		t.Fatalf("offset after bottom track click = %d, want near end", app.view.Offset)
	}
	app.view.Home()
	track := app.height - 1
	start, _ := app.view.ScrollbarThumb(track)
	app.HandleMouse(tcell.NewEventMouse(39, start, tcell.Button1, tcell.ModNone))
	app.HandleMouse(tcell.NewEventMouse(39, start+5, tcell.Button1, tcell.ModNone))
	if app.view.Offset <= 0 {
		t.Fatalf("offset after thumb drag = %d, want positive", app.view.Offset)
	}
	app.HandleMouse(tcell.NewEventMouse(39, start+5, tcell.ButtonNone, tcell.ModNone))
	if app.draggingScrollbar {
		t.Fatal("thumb drag remained active after button release")
	}
}

func TestMouseEventsIgnoredWithoutCapabilityAndKeyboardRemainsAvailable(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, lineRange(30))
	app.Resize(40, 10)
	app.HandleMouse(tcell.NewEventMouse(0, 0, tcell.WheelDown, tcell.ModNone))
	if app.view.Offset != 0 {
		t.Fatalf("mouse changed offset when disabled: %d", app.view.Offset)
	}
	app.HandleKey(key(tcell.KeyDown, 0))
	if app.view.Offset != 1 {
		t.Fatalf("keyboard fallback offset = %d, want 1", app.view.Offset)
	}
	screen := &mouseCapabilityScreen{Screen: newScreen(t, 40, 10), supported: false}
	if enableMouse(screen) || screen.enabled {
		t.Fatal("mouse input enabled without reported terminal support")
	}
}

func TestDragSelectionMapsDocumentAreaAcrossWrappedUnicodeLine(t *testing.T) {
	app := New("wrapped.md", cli.ThemeAuto, []string{strings.Repeat("a", 38) + "中文tail"})
	app.Resize(40, 10)
	app.mouseEnabled = true
	if len(app.doc.Lines) < 2 {
		t.Fatalf("layout lines = %d, want wrapped content", len(app.doc.Lines))
	}
	app.HandleMouse(tcell.NewEventMouse(38, 0, tcell.Button1, tcell.ModNone))
	app.HandleMouse(tcell.NewEventMouse(2, 1, tcell.Button1, tcell.ModNone))
	app.HandleMouse(tcell.NewEventMouse(2, 1, tcell.ButtonNone, tcell.ModNone))
	if got, want := app.selection.plainText(app.doc.Lines), "中\n文t"; got != want {
		t.Fatalf("wrapped selection = %q, want %q", got, want)
	}
	if app.selection.anchor.row != 0 || app.selection.active.row != 1 {
		t.Fatalf("selection rows = %d..%d, want 0..1", app.selection.anchor.row, app.selection.active.row)
	}
}

func TestDraggingAcrossLinkSelectsWithoutActivating(t *testing.T) {
	var activated []string
	app := NewWithOpener("links.md", cli.ThemeAuto, []string{"\x1chttps://example.com\x1cOpen link\x1c (https://example.com)"}, func(target string) error {
		activated = append(activated, target)
		return nil
	})
	app.Resize(60, 10)
	app.mouseEnabled = true
	app.HandleMouse(tcell.NewEventMouse(0, 0, tcell.Button1, tcell.ModNone))
	app.HandleMouse(tcell.NewEventMouse(9, 0, tcell.Button1, tcell.ModNone))
	app.HandleMouse(tcell.NewEventMouse(9, 0, tcell.ButtonNone, tcell.ModNone))
	if len(activated) != 0 {
		t.Fatalf("drag activated link: %v", activated)
	}
	if !app.selection.set {
		t.Fatal("drag over link did not leave a selection")
	}
}

func TestSelectionOnlyStartsInDocumentAreaAndEmptyCopyKeepsTUIOpen(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, []string{"document", "next"})
	app.Resize(40, 10)
	app.mouseEnabled = true
	app.HandleMouse(tcell.NewEventMouse(10, 8, tcell.Button1, tcell.ModNone))
	app.HandleMouse(tcell.NewEventMouse(10, 8, tcell.ButtonNone, tcell.ModNone))
	if app.selection.set {
		t.Fatal("click outside document created a selection")
	}
	if quit := app.HandleKey(key(tcell.KeyCtrlC, 0)); quit {
		t.Fatal("Ctrl+C exited the TUI without a selection")
	}
	if !strings.Contains(app.message, "Select text") {
		t.Fatalf("empty selection message = %q", app.message)
	}
}

func TestCtrlCCopiesRenderedUnicodeWithoutControlSequences(t *testing.T) {
	doc := document.New("controls.md", "café 中文 🙂 before\x1b[2J after")
	originalLines := append([]string(nil), doc.Lines...)
	writer := &fakeClipboard{}
	app := NewWithAdapters(doc.Path, cli.ThemeAuto, doc.Lines, func(string) error { return nil }, writer)
	app.Resize(80, 10)
	app.mouseEnabled = true
	line := app.doc.Lines[0]
	width := line.Width
	app.HandleMouse(tcell.NewEventMouse(0, 0, tcell.Button1, tcell.ModNone))
	app.HandleMouse(tcell.NewEventMouse(width-1, 0, tcell.Button1, tcell.ModNone))
	app.HandleMouse(tcell.NewEventMouse(width-1, 0, tcell.ButtonNone, tcell.ModNone))
	if quit := app.HandleKey(key(tcell.KeyCtrlC, 0)); quit {
		t.Fatal("Ctrl+C exited after copying")
	}
	if writer.text == "" || !strings.Contains(writer.text, "café 中文 🙂") || !strings.Contains(writer.text, "^[[2J") {
		t.Fatalf("clipboard text lacks rendered Unicode/control notation: %q", writer.text)
	}
	if strings.ContainsRune(writer.text, '\x1b') || strings.Contains(writer.text, "\x1b[") {
		t.Fatalf("clipboard contains active terminal control: %q", writer.text)
	}
	if !strings.Contains(app.message, "Copied selection") {
		t.Fatalf("copy success message = %q", app.message)
	}
	if strings.Join(doc.Lines, "\n") != strings.Join(originalLines, "\n") {
		t.Fatal("selection or copy changed the source document")
	}
}

func TestDrawShowsSelectionWithReverseVideo(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, []string{"visible text"})
	app.Resize(40, 10)
	app.selection = textSelection{anchor: cellPosition{row: 0, col: 0}, active: cellPosition{row: 0, col: 3}, set: true}
	screen := newScreen(t, 40, 10)
	app.Draw(screen)
	screen.Show()
	cells, _, _ := screen.GetContents()
	_, _, attrs := cells[0].Style.Decompose()
	if attrs&tcell.AttrReverse == 0 {
		t.Fatal("selected text has no visible reverse-video indicator")
	}
}

func TestClipboardUnavailableOrFailurePreservesSelectionAndTUI(t *testing.T) {
	for _, tc := range []struct {
		name   string
		writer *fakeClipboard
		want   string
	}{
		{name: "unavailable", want: "Clipboard is unavailable"},
		{name: "failure", writer: &fakeClipboard{err: errors.New("clipboard denied")}, want: "Cannot copy to clipboard"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var writer clipboard.Writer
			if tc.writer != nil {
				writer = tc.writer
			}
			app := NewWithAdapters("doc.md", cli.ThemeAuto, []string{"copy me"}, func(string) error { return nil }, writer)
			app.Resize(40, 10)
			app.selection = textSelection{anchor: cellPosition{row: 0, col: 0}, active: cellPosition{row: 0, col: 3}, set: true}
			if quit := app.HandleKey(key(tcell.KeyCtrlC, 0)); quit {
				t.Fatal("Ctrl+C exited after clipboard error")
			}
			if !app.selection.set || app.selection.plainText(app.doc.Lines) != "copy" {
				t.Fatalf("selection was not preserved: %#v", app.selection)
			}
			if !strings.Contains(app.message, tc.want) {
				t.Fatalf("clipboard message = %q, want %q", app.message, tc.want)
			}
		})
	}
}

type fakeClipboard struct {
	text string
	err  error
}

func (f *fakeClipboard) WriteText(text string) error {
	f.text = text
	return f.err
}

func TestDrawScrollbarThumbAtViewportPosition(t *testing.T) {
	app := New("doc.md", cli.ThemeAuto, lineRange(100))
	app.Resize(40, 10)
	app.mouseEnabled = true
	app.view.ScrollLines(45)
	screen := newScreen(t, 40, 10)
	app.Draw(screen)
	screen.Show()
	start, size := app.view.ScrollbarThumb(app.height - 1)
	for y := 0; y < app.height-1; y++ {
		main, _, _, _ := screen.GetContent(app.width-1, y)
		want := '│'
		if y >= start && y < start+size {
			want = '█'
		}
		if main != want {
			t.Errorf("scrollbar cell (%d,%d) = %q, want %q", app.width-1, y, main, want)
		}
	}
}

type mouseCapabilityScreen struct {
	tcell.Screen
	supported bool
	enabled   bool
}

func (s *mouseCapabilityScreen) HasMouse() bool { return s.supported }
func (s *mouseCapabilityScreen) EnableMouse(flags ...tcell.MouseFlags) {
	s.enabled = true
	s.Screen.EnableMouse(flags...)
}

func TestFocusedLinkHasVisibleReverseVideoIndicator(t *testing.T) {
	app := NewWithOpener("links.md", cli.ThemeAuto, []string{"\x1chttps://example.com\x1cOpen\x1c (https://example.com)"}, func(string) error { return nil })
	app.Resize(60, 10)
	app.HandleKey(key(tcell.KeyTab, 0))
	screen := newScreen(t, 60, 10)
	app.Draw(screen)
	screen.Show()
	cells, _, _ := screen.GetContents()
	_, _, attrs := cells[0].Style.Decompose()
	if attrs&tcell.AttrReverse == 0 {
		t.Fatal("focused link has no visible focus style")
	}
}

func TestFX004AnchorsNavigateAndMissingTargetReportsMessage(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "document", "testdata", "fx-004.md"))
	if err != nil {
		t.Fatal(err)
	}
	doc := document.New("fx-004.md", string(data))
	app := NewWithOpener(doc.Path, cli.ThemeAuto, doc.Lines, func(string) error { t.Fatal("fixture anchor reached browser opener"); return nil })
	app.Resize(100, 5)
	want := map[string]string{"#cara-pakai": "## Cara Pakai!", "#foo-bar": "## Foo: Bar?", "#bagian-sama-1": "## Bagian sama"}
	for i, item := range app.links() {
		if heading, ok := want[item.link.Target]; ok {
			app.focused = i
			app.activateFocused()
			if !strings.HasPrefix(app.doc.Lines[app.view.Offset].Text, heading) {
				t.Errorf("%s landed at %q, want heading %q", item.link.Target, app.doc.Lines[app.view.Offset].Text, heading)
			}
			delete(want, item.link.Target)
		}
		if item.link.Target == "#tidak-ada" {
			app.focused = i
			app.activateFocused()
			if !strings.Contains(app.message, "Anchor not found") {
				t.Errorf("missing target message = %q", app.message)
			}
		}
	}
	if len(want) != 0 {
		t.Errorf("FX-004 anchors not tested: %v", want)
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
	if !screenContains(screen, "Mouse drag") || !screenContains(screen, "Ctrl+C") {
		t.Fatal("help panel did not show selection and copy bindings")
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
