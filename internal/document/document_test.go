package document

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rumirifai/reamd/internal/layout"
)

func TestNewNormalizesLineEndings(t *testing.T) {
	doc := New("a.md", "\uFEFFfirst\r\nsecond\rthird\n")
	if got := strings.Join(doc.Lines, " "); got != "first second third" {
		t.Fatalf("Lines = %#v, want parsed paragraph text", doc.Lines)
	}
}

func TestCommonMarkFixtures(t *testing.T) {
	for _, fixture := range []string{"fx-001.md", "fx-002.md", "fx-006.md"} {
		t.Run(fixture, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", fixture))
			if err != nil {
				t.Fatal(err)
			}
			doc := New(fixture, string(data))
			joined := strings.Join(doc.Lines, "\n")
			switch fixture {
			case "fx-001.md":
				for _, s := range []string{"# Judul utama", "**teks tebal**", "*teks miring*", "`kode`", "Baris pertama\nBaris kedua", "────"} {
					if !strings.Contains(joined, s) {
						t.Errorf("render lacks %q: %s", s, joined)
					}
				}
			case "fx-002.md":
				for _, s := range []string{"1. Satu", "2. Dua", "  - Anak A", ">   Quote tingkat pertama", ">   >     Quote tingkat kedua"} {
					if !strings.Contains(joined, s) {
						t.Errorf("render lacks %q: %s", s, joined)
					}
				}
			case "fx-006.md":
				for _, s := range []string{"日本語", "café", "é", "中文", "🙂"} {
					if !strings.Contains(joined, s) {
						t.Errorf("render lacks %q", s)
					}
				}
			}
			for _, width := range []int{40, 80, 120} {
				laid := layout.Build(doc.Lines, width)
				for _, line := range laid.Lines {
					if line.Width > width {
						t.Errorf("width %d exceeded by %q (%d)", width, line.Text, line.Width)
					}
					if fixture == "fx-002.md" && width == 40 {
						wrapped := strings.Join(layoutTexts(laid.Lines), "\n")
						if !strings.Contains(wrapped, "    ") || !strings.Contains(wrapped, "> ") {
							t.Errorf("list/quote continuation context missing: %q", wrapped)
						}
					}
				}
			}
		})
	}
}

func layoutTexts(lines []layout.Line) []string {
	texts := make([]string, len(lines))
	for i, line := range lines {
		texts[i] = line.Text
	}
	return texts
}

func TestUnknownAndHTMLContentFallback(t *testing.T) {
	doc := New("fallback.md", "<b>literal</b>\n\n~~strike~~\n\n[label](https://example.test) ![alt](image.png)\n\n```mermaid\ngraph A\n```")
	joined := strings.Join(doc.Lines, "\n")
	for _, text := range []string{"<b>literal</b>", "~~strike~~", "label (https://example.test)", "alt (image.png)", "graph A"} {
		if !strings.Contains(joined, text) {
			t.Fatalf("fallback lost source content: %q", joined)
		}
	}
}

func TestGFMRendersTablesTasksStrikethroughAndAutolinks(t *testing.T) {
	doc := New("gfm.md", "| Name | State |\n| :--- | ---: |\n| Parser | Done |\n\n- [x] checked\n- [ ] unchecked\n\n~~removed~~ https://example.com contact@example.com")
	joined := strings.Join(doc.Lines, "\n")
	for _, want := range []string{"[x] checked", "[ ] unchecked", "~~removed~~", "https://example.com", "contact@example.com", "\x1eTABLE"} {
		if !strings.Contains(joined, want) {
			t.Errorf("render missing %q: %q", want, joined)
		}
	}
	if strings.Contains(joined, "https://example.com (https://example.com)") {
		t.Fatal("autolink URL was duplicated")
	}
	if !strings.Contains(joined, "\x1eALIGN:left\x1fright") {
		t.Fatalf("table alignment not retained: %q", joined)
	}
}

func TestResponsiveTablesKeepEveryCellAtRepresentativeWidths(t *testing.T) {
	doc := New("table.md", "| Long label | Value |\n| :--- | ---: |\n| 日本語の項目 | a complete value that must remain visible |\n| Wide | another full cell value |")
	for _, width := range []int{40, 80, 120} {
		laid := layout.BuildTables(doc.Lines, width)
		var lines []string
		for _, line := range laid.Lines {
			lines = append(lines, line.Text)
			if line.Width > width {
				t.Errorf("width %d exceeded by %q", width, line.Text)
			}
		}
		joined := strings.Join(lines, "\n")
		compact := strings.Join(strings.Fields(joined), " ")
		for _, val := range []string{"a complete value that must remain visible", "another full cell value", "日本語の項目"} {
			if !strings.Contains(compact, val) {
				t.Errorf("width %d lost cell %q in %q", width, val, joined)
			}
		}
		if width == 40 && !strings.Contains(joined, "Long label: ") {
			t.Errorf("width 40 did not stack table: %q", joined)
		}
		if width >= 80 && !strings.Contains(joined, "| Long label") {
			t.Errorf("width %d did not use table grid: %q", width, joined)
		}
	}
	aligned := New("aligned.md", "| Name | State |\n| :--- | ---: |\n| Parser | Done |")
	grid := layout.BuildTables(aligned.Lines, 80)
	if got := strings.Join(layoutTexts(grid.Lines), "\n"); !strings.Contains(got, "| Parser |  Done |") {
		t.Errorf("right-aligned table cell not padded on the left: %q", got)
	}
}

func TestNewEmpty(t *testing.T) {
	for _, content := range []string{"", "\n", "\r\n"} {
		doc := New("empty.md", content)
		if len(doc.Lines) != 0 {
			t.Fatalf("New(%q).Lines = %#v, want empty", content, doc.Lines)
		}
	}
}

func TestNewPreservesWhitespace(t *testing.T) {
	doc := New("a.md", "  indented\n\tTabbed\n")
	if got := strings.Join(doc.Lines, " "); got != "indented Tabbed" {
		t.Fatalf("Lines = %#v, want Markdown-normalized paragraph", doc.Lines)
	}
}
