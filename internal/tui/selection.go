package tui

import (
	"strings"

	"github.com/rivo/uniseg"

	"github.com/rumirifai/reamd/internal/layout"
)

type cellPosition struct {
	row int
	col int
}

type textSelection struct {
	anchor cellPosition
	active cellPosition
	set    bool
}

func (s textSelection) bounds() (cellPosition, cellPosition) {
	if s.anchor.row > s.active.row || (s.anchor.row == s.active.row && s.anchor.col > s.active.col) {
		return s.active, s.anchor
	}
	return s.anchor, s.active
}

func (s textSelection) contains(row, col, width int) bool {
	if !s.set || width <= 0 {
		return false
	}
	start, end := s.bounds()
	if row < start.row || row > end.row {
		return false
	}
	left, right := 0, int(^uint(0)>>1)
	if row == start.row {
		left = start.col
	}
	if row == end.row {
		right = end.col + 1
	}
	return col < right && col+width > left
}

func (s textSelection) plainText(lines []layout.Line) string {
	if !s.set || len(lines) == 0 {
		return ""
	}
	start, end := s.bounds()
	if start.row < 0 {
		start.row = 0
	}
	if end.row >= len(lines) {
		end.row = len(lines) - 1
	}
	if start.row > end.row {
		return ""
	}
	var result strings.Builder
	for row := start.row; row <= end.row; row++ {
		line := lines[row].Text
		from, to := 0, int(^uint(0)>>1)
		if row == start.row {
			from = start.col
		}
		if row == end.row {
			to = end.col + 1
		}
		if from < 0 {
			from = 0
		}
		if to < from {
			to = from
		}
		result.WriteString(sliceDisplayCells(line, from, to))
		if row < end.row {
			result.WriteByte('\n')
		}
	}
	return result.String()
}

func sliceDisplayCells(text string, from, to int) string {
	if to <= from {
		return ""
	}
	var result strings.Builder
	cell := 0
	graphemes := uniseg.NewGraphemes(text)
	for graphemes.Next() {
		cluster := graphemes.Str()
		width := graphemes.Width()
		next := cell + width
		if cell < to && next > from {
			result.WriteString(cluster)
		}
		cell = next
	}
	return result.String()
}
