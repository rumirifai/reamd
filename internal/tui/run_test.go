package tui

import (
	"errors"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/rumirifai/reamd/internal/cli"
)

func TestRunInitializedHandlesKeysAndRestoresScreen(t *testing.T) {
	screen := initializedScreen(t, 40, 10)

	done := make(chan error, 1)
	go func() {
		done <- runInitialized(screen, "doc.md", cli.ThemeAuto, lineRange(20), false)
	}()

	screen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	screen.InjectKey(tcell.KeyRune, 'j', tcell.ModNone)
	screen.InjectKey(tcell.KeyRune, '?', tcell.ModNone)
	screen.InjectKey(tcell.KeyRune, '?', tcell.ModNone)
	screen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("runInitialized() error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runInitialized() did not exit after q")
	}

	if width, height := screen.Size(); width != 0 || height != 0 {
		t.Fatalf("screen size after restore = %dx%d, want 0x0", width, height)
	}
}

func TestRunInitializedRestoresOnInterrupt(t *testing.T) {
	screen := initializedScreen(t, 40, 10)

	done := make(chan error, 1)
	go func() {
		done <- runInitialized(screen, "doc.md", cli.ThemeAuto, lineRange(3), false)
	}()
	screen.InjectKey(tcell.KeyRune, '?', tcell.ModNone)
	if err := screen.PostEvent(tcell.NewEventInterrupt(nil)); err != nil {
		t.Fatalf("PostEvent() error = %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("runInitialized() error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runInitialized() did not exit after interrupt")
	}
	if width, height := screen.Size(); width != 0 || height != 0 {
		t.Fatalf("screen size after interrupt = %dx%d, want 0x0", width, height)
	}
}

func TestRunInitializedRestoresAfterResize(t *testing.T) {
	screen := initializedScreen(t, 40, 10)

	done := make(chan error, 1)
	go func() {
		done <- runInitialized(screen, "doc.md", cli.ThemeAuto, lineRange(30), false)
	}()
	screen.SetSize(80, 24)
	screen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("runInitialized() error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("runInitialized() did not exit after resize and q")
	}
}

func TestWithRestoreFinishesOnNormalExit(t *testing.T) {
	screen := initializedScreen(t, 40, 10)
	if err := withRestore(screen, func() error { return nil }); err != nil {
		t.Fatalf("withRestore() error = %v", err)
	}
	if width, height := screen.Size(); width != 0 || height != 0 {
		t.Fatalf("screen size after normal exit = %dx%d, want 0x0", width, height)
	}
}

func TestWithRestoreFinishesOnError(t *testing.T) {
	screen := initializedScreen(t, 40, 10)
	want := errors.New("session failed")
	if err := withRestore(screen, func() error { return want }); !errors.Is(err, want) {
		t.Fatalf("withRestore() error = %v, want %v", err, want)
	}
	if width, height := screen.Size(); width != 0 || height != 0 {
		t.Fatalf("screen size after error = %dx%d, want 0x0", width, height)
	}
}

func TestWithRestoreFinishesBeforeRepanic(t *testing.T) {
	screen := initializedScreen(t, 40, 10)
	var recovered any
	func() {
		defer func() { recovered = recover() }()
		_ = withRestore(screen, func() error { panic("boom") })
	}()
	if recovered != "boom" {
		t.Fatalf("recovered = %v, want boom", recovered)
	}
	if width, height := screen.Size(); width != 0 || height != 0 {
		t.Fatalf("screen size after panic = %dx%d, want 0x0", width, height)
	}
}
