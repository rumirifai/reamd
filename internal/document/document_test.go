package document

import (
	"reflect"
	"testing"
)

func TestNewNormalizesLineEndings(t *testing.T) {
	doc := New("a.md", "\uFEFFfirst\r\nsecond\rthird\n")
	want := []string{"first", "second", "third"}
	if !reflect.DeepEqual(doc.Lines, want) {
		t.Fatalf("Lines = %#v, want %#v", doc.Lines, want)
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
	want := []string{"  indented", "\tTabbed"}
	if !reflect.DeepEqual(doc.Lines, want) {
		t.Fatalf("Lines = %#v, want %#v", doc.Lines, want)
	}
}
