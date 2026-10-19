// Package cli parses the reamd command line and owns the user-facing help text.
package cli

import (
	"fmt"
	"strings"
)

// Version is the development version reported by --version. Release builds may
// replace it with -ldflags "-X github.com/rumirifai/reamd/internal/cli.Version=...".
var Version = "0.1.0-dev"

// Theme is one of the color themes accepted by --theme. Palettes are applied in
// a later ticket; this package only validates the value.
type Theme string

const (
	ThemeAuto        Theme = "auto"
	ThemeDark        Theme = "dark"
	ThemeLight       Theme = "light"
	ThemeTransparent Theme = "transparent"
)

// Options is the result of a successful command-line parse.
type Options struct {
	Path    string
	Theme   Theme
	Help    bool
	Version bool
}

// UsageError reports invalid command-line usage. The caller maps it to exit
// code 2 and writes the message to stderr.
type UsageError struct {
	Message string
}

func (e *UsageError) Error() string {
	return e.Message
}

func usageError(format string, args ...any) *UsageError {
	return &UsageError{Message: fmt.Sprintf(format, args...)}
}

// Parse validates args, excluding the program name. No arguments means help.
func Parse(args []string) (Options, error) {
	opts := Options{Theme: ThemeAuto}
	if len(args) == 0 {
		opts.Help = true
		return opts, nil
	}

	// Help and version are checked before other arguments so they keep working
	// even when another argument is invalid. Arguments after "--" are paths.
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "--help" || arg == "-h" {
			opts.Help = true
			return opts, nil
		}
	}
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "--version" || arg == "-v" {
			opts.Version = true
			return opts, nil
		}
	}

	var paths []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--":
			paths = append(paths, args[i+1:]...)
			i = len(args)
		case arg == "--theme":
			if i+1 >= len(args) {
				return opts, usageError("option --theme requires a value")
			}
			i++
			theme, err := parseTheme(args[i])
			if err != nil {
				return opts, err
			}
			opts.Theme = theme
		case strings.HasPrefix(arg, "--theme="):
			theme, err := parseTheme(strings.TrimPrefix(arg, "--theme="))
			if err != nil {
				return opts, err
			}
			opts.Theme = theme
		case arg == "-":
			return opts, usageError("reading from stdin is not supported; pass a file path")
		case strings.HasPrefix(arg, "-"):
			return opts, usageError("unknown option %q", arg)
		default:
			paths = append(paths, arg)
		}
	}

	if len(paths) == 0 {
		return opts, usageError("missing file path")
	}
	if len(paths) > 1 {
		return opts, usageError("expected one file path, got %d", len(paths))
	}
	opts.Path = paths[0]
	return opts, nil
}

func parseTheme(value string) (Theme, error) {
	switch Theme(value) {
	case ThemeAuto, ThemeDark, ThemeLight, ThemeTransparent:
		return Theme(value), nil
	default:
		return "", usageError("invalid theme %q: expected auto, dark, light, or transparent", value)
	}
}

// HelpText returns the English help shown by --help and by a run without
// arguments.
func HelpText() string {
	return `reamd - read a Markdown file in your terminal

Usage:
  reamd [--theme <auto|dark|light|transparent>] <path>
  reamd --help
  reamd --version

Options:
  --theme <value>   Color theme: auto (default), dark, light, transparent
  -h, --help        Show this help and exit
  -v, --version     Show the version and exit

Keys:
  Up, Down, j, k     Scroll one line
  PageUp, PageDown   Scroll one screen
  Home, End          Jump to the start or end
  ?                  Show or hide the key help
  q, Esc             Quit and return to the shell

Examples:
  reamd README.md
  reamd --theme dark docs/guide.md
  reamd README.md > README.txt

With no path, reamd prints this help. Reading from stdin is not supported.
`
}

// VersionText returns the line printed by --version.
func VersionText() string {
	return "reamd " + Version + "\n"
}
