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
