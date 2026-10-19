package layout

import (
	"reflect"
	"strings"
	"testing"
)

func TestSanitizeControlCharacters(t *testing.T) {
	got := Sanitize("before\x1b[2Jafter\x07bell\tdone")
	want := "before^[[2Jafter^Gbell    done"
	if got != want {
		t.Fatalf("Sanitize() = %q, want %q", got, want)
	}
	for _, r := range got {
		if r < 0x20 || r == 0x7f {
			t.Fatalf("Sanitize() left control character %U in %q", r, got)
		}
	}
}

func TestBuildWrapsASCII(t *testing.T) {
	got := Build([]string{"alpha beta gamma"}, 10)
	want := []string{"alpha beta", "gamma"}
	if !reflect.DeepEqual(texts(got.Lines), want) {
		t.Fatalf("Build() = %#v, want %#v", texts(got.Lines), want)
	}
	for _, line := range got.Lines {
		if line.Width > 10 {
			t.Fatalf("line %q width = %d, want <= 10", line.Text, line.Width)
		}
	}
}

func TestBuildWrapsWideCharactersByDisplayWidth(t *testing.T) {
	got := Build([]string{"日本語の見出し"}, 6)
	for _, line := range got.Lines {
		if line.Width > 6 {
			t.Fatalf("line %q width = %d, want <= 6", line.Text, line.Width)
		}
	}
	if joined := strings.Join(texts(got.Lines), ""); joined != "日本語の見出し" {
		t.Fatalf("wrapped content = %q, want no content loss", joined)
	}
}

func TestBuildKeepsEmptyLines(t *testing.T) {
	got := Build([]string{"one", "", "three"}, 20)
	want := []string{"one", "", "three"}
	if !reflect.DeepEqual(texts(got.Lines), want) {
		t.Fatalf("Build() = %#v, want %#v", texts(got.Lines), want)
	}
}

func TestBuildBreaksLongWord(t *testing.T) {
	got := Build([]string{"abcdefghij"}, 4)
	want := []string{"abcd", "efgh", "ij"}
	if !reflect.DeepEqual(texts(got.Lines), want) {
		t.Fatalf("Build() = %#v, want %#v", texts(got.Lines), want)
	}
}

func TestViewportScrollAndClamp(t *testing.T) {
	v := Viewport{}
	v.Resize(3, 10)

	v.ScrollLines(2)
	if v.Offset != 2 {
		t.Fatalf("Offset = %d, want 2", v.Offset)
	}
	v.End()
	if v.Offset != 7 {
		t.Fatalf("Offset after End = %d, want 7", v.Offset)
	}
	v.ScrollLines(100)
	if v.Offset != 7 {
		t.Fatalf("Offset after overscroll = %d, want 7", v.Offset)
	}
	v.Home()
	if v.Offset != 0 {
		t.Fatalf("Offset after Home = %d, want 0", v.Offset)
	}
	v.ScrollLines(-5)
	if v.Offset != 0 {
		t.Fatalf("Offset after underscroll = %d, want 0", v.Offset)
	}
}

func TestViewportPageAndResize(t *testing.T) {
	v := Viewport{}
	v.Resize(10, 100)
	v.ScrollPage(1)
	if v.Offset != 9 {
		t.Fatalf("Offset after page down = %d, want 9", v.Offset)
	}
	v.ScrollPage(-1)
	if v.Offset != 0 {
		t.Fatalf("Offset after page up = %d, want 0", v.Offset)
	}

	v.End()
	v.Resize(20, 100)
	if v.Offset != 80 {
		t.Fatalf("Offset after resize = %d, want 80", v.Offset)
	}
	v.Resize(20, 5)
	if v.Offset != 0 {
		t.Fatalf("Offset after shrinking document = %d, want 0", v.Offset)
	}
}

func TestViewportPercent(t *testing.T) {
	v := Viewport{}
	v.Resize(10, 10)
	if got := v.Percent(); got != 100 {
		t.Fatalf("Percent() = %d, want 100 when everything fits", got)
	}
	v.Resize(10, 20)
	v.ScrollLines(5)
	if got := v.Percent(); got != 50 {
		t.Fatalf("Percent() = %d, want 50", got)
	}
}

func texts(lines []Line) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		out[i] = line.Text
	}
	return out
}
