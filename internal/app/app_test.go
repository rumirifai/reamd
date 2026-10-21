package app

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rumirifai/reamd/internal/cli"
)

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
	return path
}

func runCLI(t *testing.T, args []string, interactive bool) (code int, stdout, stderr string) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = run(args, &out, &errOut, interactive, func(string, cli.Theme, []string) error {
		return nil
	})
	return code, out.String(), errOut.String()
}

func TestRunWithoutArgumentsPrintsHelp(t *testing.T) {
	code, stdout, stderr := runCLI(t, nil, false)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout, "Usage:") || !strings.Contains(stdout, "reamd README.md") {
		t.Fatalf("stdout does not contain help:\n%s", stdout)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
}

func TestPipeAndTUIReceiveParsedDocument(t *testing.T) {
	path := writeTemp(t, "render.md", "# Title\n\n**bold** and `code`\n")
	code, stdout, stderr := runCLI(t, []string{path}, false)
	if code != 0 || stderr != "" {
		t.Fatalf("pipe run = %d, stderr %q", code, stderr)
	}
	for _, want := range []string{"# Title", "**bold**", "`code`"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("pipe output lacks %q: %q", want, stdout)
		}
	}
	if strings.Contains(stdout, "\x1b") {
		t.Fatalf("pipe output contains ANSI escape: %q", stdout)
	}
	var tuiLines []string
	var out, errOut bytes.Buffer
	code = run([]string{path}, &out, &errOut, true, func(_ string, _ cli.Theme, lines []string) error { tuiLines = lines; return nil })
	if code != 0 || errOut.Len() != 0 {
		t.Fatalf("TUI run = %d, stderr %q", code, errOut.String())
	}
	if !strings.Contains(strings.Join(tuiLines, "\n"), "**bold**") {
		t.Fatalf("TUI did not receive parsed lines: %#v", tuiLines)
	}
}

func TestRunHelpAndVersion(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{args: []string{"--help"}, want: "Usage:"},
		{args: []string{"-h"}, want: "Usage:"},
		{args: []string{"--version"}, want: "reamd "},
		{args: []string{"-v"}, want: "reamd "},
	}
	for _, tt := range tests {
		code, stdout, stderr := runCLI(t, tt.args, true)
		if code != 0 {
			t.Fatalf("%v exit code = %d, want 0", tt.args, code)
		}
		if !strings.Contains(stdout, tt.want) {
			t.Fatalf("%v stdout = %q, want it to contain %q", tt.args, stdout, tt.want)
		}
		if stderr != "" {
			t.Fatalf("%v stderr = %q, want empty", tt.args, stderr)
		}
	}
}

func TestRunUsageErrors(t *testing.T) {
	tests := [][]string{
		{"--theme", "neon", "README.md"},
		{"a.md", "b.md"},
		{"--bogus", "README.md"},
		{"--theme"},
		{"-"},
	}
	for _, args := range tests {
		code, stdout, stderr := runCLI(t, args, false)
		if code != 2 {
			t.Fatalf("%v exit code = %d, want 2", args, code)
		}
		if stdout != "" {
			t.Fatalf("%v stdout = %q, want empty", args, stdout)
		}
		if !strings.HasPrefix(stderr, "reamd: ") || !strings.Contains(stderr, "Usage:") {
			t.Fatalf("%v stderr does not contain a usage error and help:\n%s", args, stderr)
		}
	}
}

func TestRunMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.md")
	code, stdout, stderr := runCLI(t, []string{path}, false)
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "does not exist") || !strings.Contains(stderr, "missing.md") {
		t.Fatalf("stderr = %q, want a missing-file message naming the path", stderr)
	}
}

func TestRunInvalidFileContent(t *testing.T) {
	notUTF8 := writeTemp(t, "latin1.md", string([]byte{0x63, 0x61, 0x66, 0xe9}))
	binary := writeTemp(t, "binary.md", string([]byte{'a', 0x00, 'b'}))
	dir := t.TempDir()

	tests := []struct {
		path string
		want string
	}{
		{path: notUTF8, want: "not valid UTF-8"},
		{path: binary, want: "binary file"},
		{path: dir, want: "not a regular file"},
	}
	for _, tt := range tests {
		code, _, stderr := runCLI(t, []string{tt.path}, false)
		if code != 1 {
			t.Fatalf("%s exit code = %d, want 1", tt.path, code)
		}
		if !strings.Contains(stderr, tt.want) {
			t.Fatalf("%s stderr = %q, want it to contain %q", tt.path, stderr, tt.want)
		}
	}
}

func TestRunPipeOutputIsPlain(t *testing.T) {
	path := writeTemp(t, "plain.md", "# Title\n\nBody with **bold**.\n")
	code, stdout, stderr := runCLI(t, []string{path}, false)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
	if strings.Contains(stdout, "\x1b") {
		t.Fatalf("stdout contains ANSI escape:\n%q", stdout)
	}
	if stdout != "# Title\n\nBody with **bold**.\n" {
		t.Fatalf("stdout = %q, want the source text with plain newlines", stdout)
	}
}

func TestRunPipeSanitizesControlCharacters(t *testing.T) {
	path := writeTemp(t, "controls.md", "before\x1b[2Jafter\n")
	code, stdout, _ := runCLI(t, []string{path}, false)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if strings.Contains(stdout, "\x1b") {
		t.Fatalf("stdout contains an escape character: %q", stdout)
	}
	if !strings.Contains(stdout, "^[[2J") {
		t.Fatalf("stdout = %q, want a visible ESC notation", stdout)
	}
}

func TestRunInteractiveDoesNotWriteDocumentToStdout(t *testing.T) {
	path := writeTemp(t, "doc.md", "# Title\n")
	var out, errOut bytes.Buffer
	called := false
	code := run([]string{path}, &out, &errOut, true, func(gotPath string, theme cli.Theme, lines []string) error {
		called = true
		if gotPath != path {
			t.Errorf("TUI path = %q, want %q", gotPath, path)
		}
		if theme != cli.ThemeAuto {
			t.Errorf("TUI theme = %q, want auto", theme)
		}
		if len(lines) != 1 || lines[0] != "# Title" {
			t.Errorf("TUI lines = %#v, want the document lines", lines)
		}
		return nil
	})
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !called {
		t.Fatal("interactive run did not start the TUI")
	}
	if out.Len() != 0 || errOut.Len() != 0 {
		t.Fatalf("stdout = %q, stderr = %q, want both empty", out.String(), errOut.String())
	}
}

func TestRunInteractiveBadPathDoesNotStartTUI(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.md")
	var out, errOut bytes.Buffer
	called := false
	code := run([]string{path}, &out, &errOut, true, func(string, cli.Theme, []string) error {
		called = true
		return nil
	})
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if called {
		t.Fatal("TUI started before the file error was reported")
	}
	if !strings.Contains(errOut.String(), "does not exist") {
		t.Fatalf("stderr = %q, want a file error", errOut.String())
	}
}

func TestRunInteractivePassesTheme(t *testing.T) {
	path := writeTemp(t, "doc.md", "text\n")
	var out, errOut bytes.Buffer
	var got cli.Theme
	code := run([]string{"--theme", "transparent", path}, &out, &errOut, true, func(_ string, theme cli.Theme, _ []string) error {
		got = theme
		return nil
	})
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if got != cli.ThemeTransparent {
		t.Fatalf("theme = %q, want transparent", got)
	}
}

func TestRunTUIErrorReturnsOne(t *testing.T) {
	path := writeTemp(t, "doc.md", "text\n")
	var out, errOut bytes.Buffer
	code := run([]string{path}, &out, &errOut, true, func(string, cli.Theme, []string) error {
		return errors.New("no terminal")
	})
	if code != 1 {
		t.Fatalf("exit code = %d, want 1", code)
	}
	if !strings.Contains(errOut.String(), "cannot start the terminal UI") {
		t.Fatalf("stderr = %q, want a terminal UI error", errOut.String())
	}
}

func TestRunUnicodeAndSpacedPath(t *testing.T) {
	path := writeTemp(t, "dokumen dengan spasi 日本語.md", "isi café 🙂\n")
	code, stdout, stderr := runCLI(t, []string{path}, false)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0 (stderr %q)", code, stderr)
	}
	if stdout != "isi café 🙂\n" {
		t.Fatalf("stdout = %q, want the file content", stdout)
	}
}

func TestRunDoesNotModifyFile(t *testing.T) {
	path := writeTemp(t, "keep.md", "# Keep\n\nunchanged\n")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	infoBefore, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if code, _, _ := runCLI(t, []string{path}, false); code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if code, _, _ := runCLI(t, []string{path}, true); code != 0 {
		t.Fatalf("interactive exit code = %d, want 0", code)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	infoAfter, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("file content changed after running")
	}
	if !infoBefore.ModTime().Equal(infoAfter.ModTime()) {
		t.Fatal("file modification time changed after running")
	}
}

func TestRunEmptyFilePipe(t *testing.T) {
	path := writeTemp(t, "empty.md", "")
	code, stdout, stderr := runCLI(t, []string{path}, false)
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if stdout != "" || stderr != "" {
		t.Fatalf("stdout = %q, stderr = %q, want both empty", stdout, stderr)
	}
}
