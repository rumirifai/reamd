package tui

import (
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/uniseg"
)

func newScreen(t *testing.T, width, height int) tcell.SimulationScreen {
	t.Helper()
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("SimulationScreen.Init() error = %v", err)
	}
	screen.SetSize(width, height)
	t.Cleanup(screen.Fini)
	return screen
}

func screenContains(screen tcell.SimulationScreen, want string) bool {
	cells, width, height := screen.GetContents()
	var b strings.Builder
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cell := cells[y*width+x]
			if len(cell.Runes) > 0 {
				b.WriteString(string(cell.Runes))
			} else {
				b.WriteByte(' ')
			}
		}
		b.WriteByte('\n')
	}
	return strings.Contains(b.String(), want)
}

func runeWidth(s string) int {
	return uniseg.StringWidth(s)
}

func initializedScreen(t *testing.T, width, height int) tcell.SimulationScreen {
	t.Helper()
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("SimulationScreen.Init() error = %v", err)
	}
	screen.SetSize(width, height)
	return screen
}
