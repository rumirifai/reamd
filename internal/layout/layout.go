// Package layout turns source lines into terminal lines and tracks the visible
// window. Markdown-aware layout arrives in a later ticket; this package handles
// the plain-text vertical slice.
package layout

import (
	"fmt"
	"strings"

	"github.com/rivo/uniseg"
)

// Line is one terminal row of the document.
type Line struct {
	Text  string
	Width int
}

// Layout is the wrapped document at one terminal width.
type Layout struct {
	Width int
	Lines []Line
}

// Build lays out document lines, including responsive table records.
func Build(source []string, width int) *Layout {
	return BuildTables(source, width)
}

// BuildTables lays out GFM table rows as a grid when all columns fit, otherwise
// as label/value pairs. The renderer preserves every cell and wraps by display width.
func BuildTables(source []string, width int) *Layout {
	if width < 1 {
		width = 1
	}
	out := &Layout{Width: width}
	for i := 0; i < len(source); {
		if source[i] != "\x1eTABLE" {
			appendWrapped(out, source[i], width)
			i++
			continue
		}
		i++
		var rows [][]string
		var alignments []string
		for i < len(source) && source[i] != "\x1eENDTABLE" {
			if strings.HasPrefix(source[i], "\x1eROW:") {
				rows = append(rows, strings.Split(strings.TrimPrefix(source[i], "\x1eROW:"), "\x1f"))
			}
			if strings.HasPrefix(source[i], "\x1eALIGN:") {
				alignments = strings.Split(strings.TrimPrefix(source[i], "\x1eALIGN:"), "\x1f")
			}
			i++
		}
		if i < len(source) {
			i++
		}
		if len(rows) == 0 {
			continue
		}
		cols := len(rows[0])
		widths := make([]int, cols)
		for _, row := range rows {
			for c, cell := range row {
				if c < cols && uniseg.StringWidth(Sanitize(cell)) > widths[c] {
					widths[c] = uniseg.StringWidth(Sanitize(cell))
				}
			}
		}
		need := 3*cols + 1
		for _, cw := range widths {
			need += cw
		}
		if need <= width {
			for r, row := range rows {
				var b strings.Builder
				b.WriteByte('|')
				for c := 0; c < cols; c++ {
					cell := ""
					if c < len(row) {
						cell = Sanitize(row[c])
					}
					pad := widths[c] - uniseg.StringWidth(cell)
					left, right := 0, pad
					if c < len(alignments) && alignments[c] == "right" {
						left, right = pad, 0
					}
					if c < len(alignments) && alignments[c] == "center" {
						left, right = pad/2, pad-pad/2
					}
					b.WriteByte(' ')
					b.WriteString(strings.Repeat(" ", left))
					b.WriteString(cell)
					b.WriteString(strings.Repeat(" ", right))
					b.WriteString(" |")
				}
				appendWrapped(out, b.String(), width)
				if r == 0 {
					sep := "|"
					for _, cw := range widths {
						sep += strings.Repeat("-", cw+2) + "|"
					}
					appendWrapped(out, sep, width)
				}
			}
		} else {
			for _, row := range rows[1:] {
				for c, label := range rows[0] {
					val := ""
					if c < len(row) {
						val = row[c]
					}
					appendWrapped(out, Sanitize(label)+": "+Sanitize(val), width)
				}
				appendWrapped(out, "", width)
			}
		}
	}
	return out
}

func appendWrapped(out *Layout, raw string, width int) {
	text := Sanitize(raw)
	prefix := continuationPrefix(text)
	parts := wrap(text, width)
	if prefix != "" {
		content := strings.TrimLeft(text[len(prefix):], " ")
		continuation := strings.Repeat(" ", uniseg.StringWidth(prefix))
		if strings.HasPrefix(prefix, ">") {
			continuation = prefix
		}
		parts = wrap(prefix+content, width)
		for i := 1; i < len(parts); i++ {
			parts[i] = continuation + strings.TrimLeft(parts[i], " ")
		}
	}
	for _, part := range parts {
		out.Lines = append(out.Lines, Line{Text: part, Width: uniseg.StringWidth(part)})
	}
}

func continuationPrefix(s string) string {
	i := 0
	for i < len(s) && s[i] == ' ' {
		i++
	}
	if i < len(s) && s[i] == '>' {
		j := i + 1
		if j < len(s) && s[j] == ' ' {
			j++
		}
		return s[:j]
	}
	if i < len(s) && (s[i] == '-' || s[i] == '*' || s[i] == '+') && i+1 < len(s) && s[i+1] == ' ' {
		return s[:i+2]
	}
	for j := i; j < len(s) && s[j] >= '0' && s[j] <= '9'; j++ {
		if j+1 < len(s) && (s[j+1] == '.' || s[j+1] == ')') && s[j+1:j+2] != "" {
			k := j + 2
			if k < len(s) && s[k] == ' ' {
				return s[:k+1]
			}
		}
	}
	return ""
}

// Sanitize replaces control characters with printable notation and expands
// tabs. Document content must never reach the terminal as an active escape
// sequence.
func Sanitize(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == '\t':
			b.WriteString("    ")
		case r == '\n' || r == '\r':
			b.WriteByte(' ')
		case r < 0x20:
			b.WriteByte('^')
			b.WriteRune(r + 0x40)
		case r == 0x7f:
			b.WriteString("^?")
		case r >= 0x80 && r <= 0x9f:
			fmt.Fprintf(&b, "\\u%04x", r)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// wrap breaks text at grapheme cluster boundaries, preferring the last space
// that fits. A single grapheme wider than the available width gets its own
// line so no content is dropped.
func wrap(text string, width int) []string {
	if text == "" {
		return []string{""}
	}

	var clusters []string
	g := uniseg.NewGraphemes(text)
	for g.Next() {
		clusters = append(clusters, g.Str())
	}

	var lines []string
	start := 0
	for start < len(clusters) {
		end := start
		used := 0
		lastSpace := -1
		for end < len(clusters) {
			cw := uniseg.StringWidth(clusters[end])
			if used+cw > width && end > start {
				break
			}
			used += cw
			if clusters[end] == " " {
				lastSpace = end
			}
			end++
		}

		cut := end
		if end < len(clusters) && clusters[end] != " " && lastSpace > start {
			// Break before the word that does not fit.
			cut = lastSpace
		}
		lines = append(lines, strings.TrimRight(strings.Join(clusters[start:cut], ""), " "))
		start = cut
		for start < len(clusters) && clusters[start] == " " {
			start++
		}
	}
	return lines
}

// Viewport tracks the first visible line and the visible height.
type Viewport struct {
	Offset int
	Height int
	Total  int
}

// Resize updates the height and clamps the offset.
func (v *Viewport) Resize(height, total int) {
	if height < 1 {
		height = 1
	}
	v.Height = height
	v.Total = total
	v.Clamp()
}

// Clamp keeps Offset within the scrollable range.
func (v *Viewport) Clamp() {
	maxOffset := v.Total - v.Height
	if maxOffset < 0 {
		maxOffset = 0
	}
	if v.Offset > maxOffset {
		v.Offset = maxOffset
	}
	if v.Offset < 0 {
		v.Offset = 0
	}
}

// ScrollLines moves by n lines, positive for down.
func (v *Viewport) ScrollLines(n int) {
	v.Offset += n
	v.Clamp()
}

// ScrollPage moves by roughly one screen, keeping one line of context.
func (v *Viewport) ScrollPage(n int) {
	step := v.Height - 1
	if step < 1 {
		step = 1
	}
	v.ScrollLines(n * step)
}

// Home moves to the first line.
func (v *Viewport) Home() {
	v.Offset = 0
}

// End moves to the last screen.
func (v *Viewport) End() {
	v.Offset = v.Total
	v.Clamp()
}

// Visible returns the lines in the current window.
func (v *Viewport) Visible(lines []Line) []Line {
	if v.Offset >= len(lines) {
		return nil
	}
	end := v.Offset + v.Height
	if end > len(lines) {
		end = len(lines)
	}
	return lines[v.Offset:end]
}

// Percent reports how far through the document the viewport is, from 0 to 100.
func (v *Viewport) Percent() int {
	maxOffset := v.Total - v.Height
	if maxOffset <= 0 {
		return 100
	}
	p := v.Offset * 100 / maxOffset
	if p < 0 {
		p = 0
	}
	if p > 100 {
		p = 100
	}
	return p
}
