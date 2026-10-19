package cli

import (
	"strings"
	"testing"
)

func TestParseNoArgumentsShowsHelp(t *testing.T) {
	opts, err := Parse(nil)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if !opts.Help {
		t.Fatal("Parse(nil).Help = false, want true")
	}
}

func TestParseHelpAndVersion(t *testing.T) {
	tests := []struct {
		name string
		args []string
		help bool
		ver  bool
	}{
		{name: "long help", args: []string{"--help"}, help: true},
		{name: "short help", args: []string{"-h"}, help: true},
		{name: "help with path", args: []string{"README.md", "--help"}, help: true},
		{name: "long version", args: []string{"--version"}, ver: true},
		{name: "short version", args: []string{"-v"}, ver: true},
		{name: "version with path", args: []string{"--version", "README.md"}, ver: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := Parse(tt.args)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", tt.args, err)
			}
			if opts.Help != tt.help || opts.Version != tt.ver {
				t.Fatalf("Parse(%q) = help:%v version:%v, want help:%v version:%v",
					tt.args, opts.Help, opts.Version, tt.help, tt.ver)
			}
		})
	}
}

func TestParseTheme(t *testing.T) {
	tests := []struct {
		args  []string
		theme Theme
	}{
		{args: []string{"README.md"}, theme: ThemeAuto},
		{args: []string{"--theme", "dark", "README.md"}, theme: ThemeDark},
		{args: []string{"--theme", "light", "README.md"}, theme: ThemeLight},
		{args: []string{"--theme", "transparent", "README.md"}, theme: ThemeTransparent},
		{args: []string{"--theme=dark", "README.md"}, theme: ThemeDark},
		{args: []string{"README.md", "--theme", "light"}, theme: ThemeLight},
	}
	for _, tt := range tests {
		opts, err := Parse(tt.args)
		if err != nil {
			t.Fatalf("Parse(%q) error = %v", tt.args, err)
		}
		if opts.Theme != tt.theme {
			t.Fatalf("Parse(%q).Theme = %q, want %q", tt.args, opts.Theme, tt.theme)
		}
	}
}

func TestParseUsageErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "invalid theme", args: []string{"--theme", "neon", "README.md"}, want: `invalid theme "neon"`},
		{name: "empty theme", args: []string{"--theme=", "README.md"}, want: `invalid theme ""`},
		{name: "missing theme value", args: []string{"--theme"}, want: "--theme requires a value"},
		{name: "unknown option", args: []string{"--color", "README.md"}, want: `unknown option "--color"`},
		{name: "two paths", args: []string{"a.md", "b.md"}, want: "expected one file path, got 2"},
		{name: "stdin", args: []string{"-"}, want: "stdin is not supported"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.args)
			if err == nil {
				t.Fatalf("Parse(%q) error = nil, want usage error", tt.args)
			}
			var usage *UsageError
			if !asUsageError(err, &usage) {
				t.Fatalf("Parse(%q) error type = %T, want *UsageError", tt.args, err)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Parse(%q) error = %q, want it to contain %q", tt.args, err, tt.want)
			}
		})
	}
}

func asUsageError(err error, target **UsageError) bool {
	u, ok := err.(*UsageError)
	if ok {
		*target = u
	}
	return ok
}

func TestParseDoubleDashAllowsDashPath(t *testing.T) {
	opts, err := Parse([]string{"--", "-notes.md"})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if opts.Path != "-notes.md" {
		t.Fatalf("Path = %q, want -notes.md", opts.Path)
	}
}

func TestHelpTextHasRequiredContent(t *testing.T) {
	help := HelpText()
	for _, want := range []string{
		"Usage:",
		"--theme",
		"dark, light, transparent",
		"reamd README.md",
		"Up, Down, j, k",
		"q, Esc",
		"?",
	} {
		if !strings.Contains(help, want) {
			t.Errorf("HelpText() does not contain %q", want)
		}
	}
	if strings.Contains(help, "\x1b") {
		t.Error("HelpText() contains an escape character")
	}
}

func TestVersionText(t *testing.T) {
	if got := VersionText(); !strings.HasPrefix(got, "reamd ") || !strings.HasSuffix(got, "\n") {
		t.Fatalf("VersionText() = %q, want a reamd version line", got)
	}
}
