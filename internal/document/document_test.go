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
	doc := New("fallback.md", "<b>literal</b>\n\n<script>alert(\"do not run\")</script>\n\n~~strike~~\n\n[label](https://example.test) ![alt](image.png)\n\n```mermaid\ngraph A\n```")
	joined := strings.Join(doc.Lines, "\n")
	for _, text := range []string{"<b>literal</b>", "<script>alert(\"do not run\")</script>", "~~strike~~", "label (https://example.test)", "alt (image.png)", "◇ Mermaid diagram", "graph A"} {
		if !strings.Contains(joined, text) {
			t.Fatalf("fallback lost source content: %q", joined)
		}
	}
}

func TestCodeBlockAndInlineCodeMarkers(t *testing.T) {
	doc := New("code.md", "Inline `x := 1`\n\n```go\n  fmt.Println(\"hello\")\n```\n")
	joined := strings.Join(doc.Lines, "\n")
	for _, want := range []string{"`x := 1`", "\x1dCODE:go", "  fmt.Println(\"hello\")", "\x1dENDCODE"} {
		if !strings.Contains(joined, want) {
			t.Errorf("document code model lacks %q: %q", want, joined)
		}
	}
}

func TestDocumentControlBytesCannotForgeLayoutRecords(t *testing.T) {
	doc := New("markers.md", "before\x1dCODE: injected\x1eTABLE\x1fafter")
	joined := strings.Join(doc.Lines, "\n")
	if !strings.Contains(joined, "before^]CODE: injected^^TABLE^_after") {
		t.Fatalf("reserved control bytes were not made visible: %q", joined)
	}
	laid := layout.Build(doc.Lines, 80)
	if got := strings.Join(layoutTexts(laid.Lines), "\n"); !strings.Contains(got, "before^]CODE: injected^^TABLE^_after") {
		t.Fatalf("document controls were treated as layout metadata: %q", got)
	}
}

func TestTicket05Fixtures(t *testing.T) {
	for _, fixture := range []string{"fx-005.md", "fx-007.md", "fx-008.md"} {
		t.Run(fixture, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", fixture))
			if err != nil {
				t.Fatal(err)
			}
			doc := New(fixture, string(data))
			laid := layout.Build(doc.Lines, 40)
			var rendered []string
			for _, line := range laid.Lines {
				rendered = append(rendered, line.Text)
				if line.Width > 40 {
					t.Errorf("line exceeds width 40: %q", line.Text)
				}
			}
			joined := strings.Join(rendered, "\n")
			switch fixture {
			case "fx-005.md":
				for _, want := range []string{"code (go)", "│ ", "↪ ", "println", "`x := 1`"} {
					if !strings.Contains(joined, want) {
						t.Errorf("render lacks %q: %q", want, joined)
					}
				}
			case "fx-007.md":
				if !strings.Contains(joined, "<b>Tag tebal literal</b>") {
					t.Errorf("literal HTML lacks bold tag: %q", joined)
				}
				literal := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(joined, "│ ", ""), "↪ ", ""), "\n", "")
				if !strings.Contains(literal, "<script>alert(\"jangan dijalankan\")</script>") {
					t.Errorf("literal HTML script text was lost: %q", joined)
				}
			case "fx-008.md":
				for _, want := range []string{"Mermaid diagram", "graph TD", "A[Mulai] --> B[Selesai]"} {
					if !strings.Contains(joined, want) {
						t.Errorf("Mermaid fallback lacks %q: %q", want, joined)
					}
				}
			}
		})
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
