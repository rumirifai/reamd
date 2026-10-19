package source

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
	return path
}

func TestLoadUTF8WithUnicodeAndSpaces(t *testing.T) {
	const content = "# Judul 日本語\n\ncafé 🙂\n"
	path := writeFile(t, "dokumen dengan spasi 日本語.md", content)

	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got != content {
		t.Fatalf("Load() = %q, want %q", got, content)
	}
}

func TestLoadDoesNotChangeFile(t *testing.T) {
	path := writeFile(t, "keep.md", "original content\n")
	infoBefore, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := Load(path); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	infoAfter, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("file content changed after Load()")
	}
	if !infoBefore.ModTime().Equal(infoAfter.ModTime()) {
		t.Fatal("file modification time changed after Load()")
	}
}

func TestLoadErrors(t *testing.T) {
	dir := t.TempDir()
	notUTF8 := writeFile(t, "latin1.md", string([]byte{0x63, 0x61, 0x66, 0xe9}))
	binary := writeFile(t, "image.md", string([]byte{0x89, 'P', 'N', 'G', 0x00, 0x1a}))

	tests := []struct {
		name string
		path string
		kind Kind
	}{
		{name: "missing", path: filepath.Join(dir, "missing.md"), kind: KindNotExist},
		{name: "directory", path: dir, kind: KindNotRegular},
		{name: "not utf8", path: notUTF8, kind: KindNotUTF8},
		{name: "binary", path: binary, kind: KindBinary},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(tt.path)
			if err == nil {
				t.Fatalf("Load(%q) error = nil, want failure", tt.path)
			}
			var loadErr *Error
			if !errors.As(err, &loadErr) {
				t.Fatalf("Load(%q) error type = %T, want *Error", tt.path, err)
			}
			if loadErr.Kind != tt.kind {
				t.Fatalf("Load(%q) kind = %v, want %v", tt.path, loadErr.Kind, tt.kind)
			}
			if loadErr.Error() == "" {
				t.Fatal("error message is empty")
			}
		})
	}
}

func TestLoadPermissionDenied(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows read-only permission does not prevent reading")
	}
	path := writeFile(t, "secret.md", "secret\n")
	if err := os.Chmod(path, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(path, 0o600) })

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want permission error")
	}
	var loadErr *Error
	if !errors.As(err, &loadErr) {
		t.Fatalf("error type = %T, want *Error", err)
	}
	if loadErr.Kind != KindPermission && loadErr.Kind != KindOther {
		t.Fatalf("kind = %v, want permission or other", loadErr.Kind)
	}
}
