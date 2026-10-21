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
	Links []Link
}

// Link marks the visible terminal cells belonging to a document link.
type Link struct {
	Start, End int
	Target     string
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
		if strings.HasPrefix(source[i], "\x1dCODE:") {
			language := strings.TrimPrefix(source[i], "\x1dCODE:")
			if language == "" {
				appendWrapped(out, "┌─ code", width)
			} else {
				appendWrapped(out, "┌─ code ("+Sanitize(language)+")", width)
			}
			i++
			for i < len(source) && source[i] != "\x1dENDCODE" {
				appendCodeWrapped(out, source[i], width)
				i++
			}
			if i < len(source) {
				i++
			}
			appendWrapped(out, "└─", width)
			i++
			continue
		}
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

func appendCodeWrapped(out *Layout, raw string, width int) {
	text := Sanitize(raw)
	const firstPrefix, continuationPrefix = "│ ", "↪ "
	remaining := text
	prefix := firstPrefix
	for {
		available := width - uniseg.StringWidth(prefix)
		if available < 1 {
			available = 1
		}
		part, rest := takeWidth(remaining, available)
		line := prefix + part
		out.Lines = append(out.Lines, Line{Text: line, Width: uniseg.StringWidth(line)})
		if rest == "" {
			break
		}
		remaining, prefix = rest, continuationPrefix
	}
}

func takeWidth(s string, width int) (string, string) {
	var b strings.Builder
	used := 0
	g := uniseg.NewGraphemes(s)
	for g.Next() {
		cluster := g.Str()
		cw := g.Width()
		if used+cw > width && used > 0 {
			return b.String(), s[len(b.String()):]
		}
		b.WriteString(cluster)
		used += cw
	}
	return b.String(), ""
}

func appendWrapped(out *Layout, raw string, width int) {
	text, targets := stripLinkMarkers(raw)
	text = Sanitize(text)
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
	base := len(out.Lines)
	for _, part := range parts {
		out.Lines = append(out.Lines, Line{Text: part, Width: uniseg.StringWidth(part)})
	}
	for _, target := range targets {
		before := text[:target.start]
		prefixWidth := uniseg.StringWidth(before)
		row, col := 0, prefixWidth
		for col >= width && row+1 < len(parts) {
			row++
			col -= width
		}
		remaining := uniseg.StringWidth(text[target.start:target.end])
		for remaining > 0 && row < len(parts) {
			available := uniseg.StringWidth(parts[row]) - col
			if available < 0 {
				available = 0
			}
			count := remaining
			if count > available {
				count = available
			}
			if count > 0 {
				out.Lines[base+row].Links = append(out.Lines[base+row].Links, Link{Start: col, End: col + count, Target: target.url})
			}
			remaining -= count
			row++
			col = 0
		}
	}
}

type linkTarget struct {
	start, end int
	url        string
}

func stripLinkMarkers(raw string) (string, []linkTarget) {
	var out strings.Builder
	var links []linkTarget
	for i := 0; i < len(raw); {
		if raw[i] != '\x1c' {
			out.WriteByte(raw[i])
			i++
			continue
		}
		endTarget := strings.IndexByte(raw[i+1:], '\x1c')
		if endTarget < 0 {
			out.WriteByte(raw[i])
			i++
			continue
		}
		endTarget += i + 1
		endLabel := strings.IndexByte(raw[endTarget+1:], '\x1c')
		if endLabel < 0 {
			out.WriteByte(raw[i])
			i++
			continue
		}
		endLabel += endTarget + 1
		start := out.Len()
		out.WriteString(raw[endTarget+1 : endLabel])
		links = append(links, linkTarget{start, out.Len(), raw[i+1 : endTarget]})
		i = endLabel + 1
	}
	return out.String(), links
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

// ScrollbarThumb maps the viewport onto a track and returns the thumb start
// and size in track cells. A non-scrollable document fills the track.
func (v Viewport) ScrollbarThumb(trackHeight int) (start, size int) {
	if trackHeight < 1 {
		return 0, 0
	}
	if v.Total <= 0 || v.Total <= v.Height || v.Height <= 0 {
		return 0, trackHeight
	}
	size = trackHeight * v.Height / v.Total
	if size < 1 {
		size = 1
	}
	if size > trackHeight {
		size = trackHeight
	}
	maxOffset := v.Total - v.Height
	maxStart := trackHeight - size
	start = v.Offset * maxStart / maxOffset
	return start, size
}

// ScrollToTrack maps a track cell to the corresponding document offset.
func (v *Viewport) ScrollToTrack(position, trackHeight int) {
	if trackHeight < 1 || v.Total <= v.Height || v.Height <= 0 {
		v.Offset = 0
		return
	}
	if position < 0 {
		position = 0
	}
	if position >= trackHeight {
		position = trackHeight - 1
	}
	_, size := v.ScrollbarThumb(trackHeight)
	maxStart := trackHeight - size
	if maxStart <= 0 {
		v.Offset = 0
		return
	}
	thumbStart := position - size/2
	if thumbStart < 0 {
		thumbStart = 0
	}
	if thumbStart > maxStart {
		thumbStart = maxStart
	}
	v.Offset = thumbStart * (v.Total - v.Height) / maxStart
	v.Clamp()
}
