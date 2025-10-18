package integration

import (
	"os"
	"testing"

	"github.com/dshills/goterm"
)

// T023: Integration test for basic rendering
// This test verifies the complete flow: Init -> SetCell -> Show -> Close

func TestBasicRendering(t *testing.T) {
	// Skip if not running in a terminal
	if !isTerminal() {
		t.Skip("Not running in a terminal, skipping integration test")
	}

	// Initialize terminal
	screen, err := goterm.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer func() {
		if err := screen.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	}()

	// Get terminal size
	w, h := screen.Size()
	if w <= 0 || h <= 0 {
		t.Errorf("Size() returned invalid dimensions: (%d, %d)", w, h)
	}

	// Clear screen
	screen.Clear()

	// Draw some content
	screen.DrawText(0, 0, "Hello, Terminal!", goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleBold)
	screen.DrawText(0, 1, "Testing basic rendering...", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)

	// Show the content
	if err := screen.Show(); err != nil {
		t.Errorf("Show() failed: %v", err)
	}

	// Sync to ensure all output is flushed
	if err := screen.Sync(); err != nil {
		t.Errorf("Sync() failed: %v", err)
	}
}

func TestRenderingWithColors(t *testing.T) {
	if !isTerminal() {
		t.Skip("Not running in a terminal, skipping integration test")
	}

	screen, err := goterm.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer func() {
		if err := screen.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	}()

	screen.Clear()

	// Test basic 16 colors
	colors := []struct {
		color goterm.Color
		name  string
	}{
		{goterm.ColorBlack, "Black"},
		{goterm.ColorRed, "Red"},
		{goterm.ColorGreen, "Green"},
		{goterm.ColorYellow, "Yellow"},
		{goterm.ColorBlue, "Blue"},
		{goterm.ColorMagenta, "Magenta"},
		{goterm.ColorCyan, "Cyan"},
		{goterm.ColorWhite, "White"},
	}

	for i, c := range colors {
		screen.DrawText(0, i, c.name, c.color, goterm.ColorDefault(), goterm.StyleNone)
	}

	if err := screen.Show(); err != nil {
		t.Errorf("Show() failed: %v", err)
	}
}

func TestRenderingWithStyles(t *testing.T) {
	if !isTerminal() {
		t.Skip("Not running in a terminal, skipping integration test")
	}

	screen, err := goterm.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer func() {
		if err := screen.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	}()

	screen.Clear()

	// Test various styles
	styles := []struct {
		style goterm.Style
		name  string
	}{
		{goterm.StyleNone, "Normal"},
		{goterm.StyleBold, "Bold"},
		{goterm.StyleDim, "Dim"},
		{goterm.StyleItalic, "Italic"},
		{goterm.StyleUnderline, "Underline"},
		{goterm.StyleReverse, "Reverse"},
		{goterm.StyleBold | goterm.StyleUnderline, "Bold+Underline"},
	}

	for i, s := range styles {
		screen.DrawText(0, i, s.name, goterm.ColorWhite, goterm.ColorDefault(), s.style)
	}

	if err := screen.Show(); err != nil {
		t.Errorf("Show() failed: %v", err)
	}
}

func TestRenderingUnicode(t *testing.T) {
	if !isTerminal() {
		t.Skip("Not running in a terminal, skipping integration test")
	}

	screen, err := goterm.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer func() {
		if err := screen.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	}()

	screen.Clear()

	// Test various Unicode content
	unicodeTexts := []string{
		"ASCII: Hello World",
		"Latin: Héllo Wörld",
		"CJK: 你好世界 こんにちは",
		"Emoji: 🎮 🚀 🌟 ⚡",
		"Box Drawing: ┌─┬─┐",
		"              │ │ │",
		"              ├─┼─┤",
		"              │ │ │",
		"              └─┴─┘",
	}

	for i, text := range unicodeTexts {
		screen.DrawText(0, i, text, goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleNone)
	}

	if err := screen.Show(); err != nil {
		t.Errorf("Show() failed: %v", err)
	}
}

func TestRenderingPerformance(t *testing.T) {
	if !isTerminal() {
		t.Skip("Not running in a terminal, skipping integration test")
	}

	screen, err := goterm.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer func() {
		if err := screen.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	}()

	w, h := screen.Size()

	// Fill entire screen
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			ch := rune('A' + (x+y)%26)
			// Safe conversion: (x+y)%7 is in range [0, 6], so 1 + (x+y)%7 is in range [1, 7]
			color := goterm.ColorIndex(uint8(1 + (x+y)%7)) // #nosec G115
			screen.SetCell(x, y, goterm.NewCell(ch, color, goterm.ColorDefault(), goterm.StyleNone))
		}
	}

	// This should complete without timing out
	if err := screen.Show(); err != nil {
		t.Errorf("Show() failed: %v", err)
	}
}

func TestMultipleShowCalls(t *testing.T) {
	if !isTerminal() {
		t.Skip("Not running in a terminal, skipping integration test")
	}

	screen, err := goterm.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer func() {
		if err := screen.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	}()

	// Multiple Show() calls should be safe
	screen.Clear()
	screen.DrawText(0, 0, "Frame 1", goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleNone)
	if err := screen.Show(); err != nil {
		t.Errorf("First Show() failed: %v", err)
	}

	screen.Clear()
	screen.DrawText(0, 0, "Frame 2", goterm.ColorBlue, goterm.ColorDefault(), goterm.StyleNone)
	if err := screen.Show(); err != nil {
		t.Errorf("Second Show() failed: %v", err)
	}

	screen.Clear()
	screen.DrawText(0, 0, "Frame 3", goterm.ColorRed, goterm.ColorDefault(), goterm.StyleNone)
	if err := screen.Show(); err != nil {
		t.Errorf("Third Show() failed: %v", err)
	}
}

func TestResizeHandling(t *testing.T) {
	if !isTerminal() {
		t.Skip("Not running in a terminal, skipping integration test")
	}

	screen, err := goterm.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer func() {
		if err := screen.Close(); err != nil {
			t.Errorf("Close() failed: %v", err)
		}
	}()

	// Get initial size
	w1, h1 := screen.Size()

	// Draw content
	screen.Clear()
	screen.DrawText(0, 0, "Initial content", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	if err := screen.Show(); err != nil {
		t.Errorf("Show() failed: %v", err)
	}

	// Note: Actual resize would require manual terminal window resize
	// This test just verifies Size() returns consistent values
	w2, h2 := screen.Size()
	if w1 != w2 || h1 != h2 {
		t.Logf("Size changed from (%d, %d) to (%d, %d) - this may be normal if terminal was resized", w1, h1, w2, h2)
	}
}

func TestInitMultipleTimes(t *testing.T) {
	if !isTerminal() {
		t.Skip("Not running in a terminal, skipping integration test")
	}

	// First init
	screen1, err := goterm.Init()
	if err != nil {
		t.Fatalf("First Init() failed: %v", err)
	}
	if err := screen1.Close(); err != nil {
		t.Errorf("First Close() failed: %v", err)
	}

	// Second init should work
	screen2, err := goterm.Init()
	if err != nil {
		t.Fatalf("Second Init() failed: %v", err)
	}
	defer func() {
		if err := screen2.Close(); err != nil {
			t.Errorf("Second Close() failed: %v", err)
		}
	}()

	// Should be usable
	screen2.Clear()
	screen2.DrawText(0, 0, "Second init test", goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleNone)
	if err := screen2.Show(); err != nil {
		t.Errorf("Show() after second init failed: %v", err)
	}
}

func TestCloseWithoutShow(t *testing.T) {
	if !isTerminal() {
		t.Skip("Not running in a terminal, skipping integration test")
	}

	screen, err := goterm.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Close without ever calling Show() - should not panic or error
	if err := screen.Close(); err != nil {
		t.Errorf("Close() without Show() failed: %v", err)
	}
}

// Helper function to check if stdout is a terminal
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
