// Package source reads the single file passed to reamd and rejects content that
// cannot be shown safely as text.
package source

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"unicode/utf8"
)

// Kind classifies a read failure. Every kind maps to exit code 1.
type Kind int

const (
	KindOther Kind = iota
	KindNotExist
	KindNotRegular
	KindPermission
	KindNotUTF8
	KindBinary
)

// Error is a file loading failure with a user-facing English message.
type Error struct {
	Path string
	Kind Kind
	Err  error
}

func (e *Error) Error() string {
	switch e.Kind {
	case KindNotExist:
		return fmt.Sprintf("%q does not exist", e.Path)
	case KindNotRegular:
		return fmt.Sprintf("%q is not a regular file", e.Path)
	case KindPermission:
		return fmt.Sprintf("cannot read %q: permission denied", e.Path)
	case KindNotUTF8:
		return fmt.Sprintf("%q is not valid UTF-8 text", e.Path)
	case KindBinary:
		return fmt.Sprintf("%q looks like a binary file", e.Path)
	default:
		return fmt.Sprintf("cannot read %q: %v", e.Path, e.Err)
	}
}

func (e *Error) Unwrap() error {
	return e.Err
}

// Load reads path as UTF-8 text. It never writes to the file.
func Load(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", statError(path, err)
	}
	if !info.Mode().IsRegular() {
		return "", &Error{Path: path, Kind: KindNotRegular}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			return "", &Error{Path: path, Kind: KindPermission, Err: err}
		}
		return "", &Error{Path: path, Kind: KindOther, Err: err}
	}

	// A NUL byte is the practical signal that the file is not meant to be read
	// as text. Other control characters are allowed here and neutralized later
	// by the layout layer.
	if bytes.IndexByte(data, 0) >= 0 {
		return "", &Error{Path: path, Kind: KindBinary}
	}
	if !utf8.Valid(data) {
		return "", &Error{Path: path, Kind: KindNotUTF8}
	}
	return string(data), nil
}

func statError(path string, err error) error {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return &Error{Path: path, Kind: KindNotExist, Err: err}
	case errors.Is(err, fs.ErrPermission):
		return &Error{Path: path, Kind: KindPermission, Err: err}
	default:
		return &Error{Path: path, Kind: KindOther, Err: err}
	}
}
