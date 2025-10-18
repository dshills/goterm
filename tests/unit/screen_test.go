package unit

import (
	"testing"

	"github.com/dshills/goterm"
)

// T020: Unit tests for Screen buffer allocation
func TestScreenBufferAllocation(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"small_screen", 10, 5},
		{"medium_screen", 80, 24},
		{"large_screen", 120, 40},
		{"single_cell", 1, 1},
		{"wide_screen", 200, 10},
		{"tall_screen", 10, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen := goterm.NewScreen(tt.width, tt.height)

			w, h := screen.Size()
			if w != tt.width {
				t.Errorf("NewScreen(%d, %d).Size() width = %d, want %d", tt.width, tt.height, w, tt.width)
			}
			if h != tt.height {
				t.Errorf("NewScreen(%d, %d).Size() height = %d, want %d", tt.width, tt.height, h, tt.height)
			}

			// Verify all cells are initialized to defaults
			for y := 0; y < h; y++ {
				for x := 0; x < w; x++ {
					cell := screen.GetCell(x, y)
					if cell.Ch != ' ' {
						t.Errorf("Initial cell at (%d, %d) has Ch = %v, want ' '", x, y, cell.Ch)
					}
					if cell.Fg != goterm.ColorDefault() {
						t.Errorf("Initial cell at (%d, %d) has non-default Fg", x, y)
					}
					if cell.Bg != goterm.ColorDefault() {
						t.Errorf("Initial cell at (%d, %d) has non-default Bg", x, y)
					}
					if cell.Style != goterm.StyleNone {
						t.Errorf("Initial cell at (%d, %d) has Style = %v, want StyleNone", x, y, cell.Style)
					}
				}
			}
		})
	}
}

func TestScreenBufferAllocationEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		width        int
		height       int
		shouldCreate bool
	}{
		{"zero_width", 0, 10, false},
		{"zero_height", 10, 0, false},
		{"negative_width", -5, 10, false},
		{"negative_height", 10, -5, false},
		{"both_zero", 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.shouldCreate && r != nil {
					t.Errorf("NewScreen(%d, %d) panicked unexpectedly: %v", tt.width, tt.height, r)
				}
				if !tt.shouldCreate && r == nil {
					t.Errorf("NewScreen(%d, %d) should panic but didn't", tt.width, tt.height)
				}
			}()

			_ = goterm.NewScreen(tt.width, tt.height)
		})
	}
}

// T021: Unit tests for SetCell/GetCell operations
func TestScreenSetGetCell(t *testing.T) {
	screen := goterm.NewScreen(80, 24)

	tests := []struct {
		name string
		x, y int
		cell goterm.Cell
	}{
		{
			name: "set_basic_char",
			x:    0, y: 0,
			cell: goterm.NewCell('A', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold),
		},
		{
			name: "set_unicode",
			x:    10, y: 5,
			cell: goterm.NewCell('🎮', goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleNone),
		},
		{
			name: "set_cjk",
			x:    20, y: 10,
			cell: goterm.NewCell('日', goterm.ColorYellow, goterm.ColorBlack, goterm.StyleUnderline),
		},
		{
			name: "set_corner_top_left",
			x:    0, y: 0,
			cell: goterm.NewCell('┌', goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone),
		},
		{
			name: "set_corner_bottom_right",
			x:    79, y: 23,
			cell: goterm.NewCell('┘', goterm.ColorMagenta, goterm.ColorDefault(), goterm.StyleNone),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen.SetCell(tt.x, tt.y, tt.cell)
			got := screen.GetCell(tt.x, tt.y)

			if !got.Equal(tt.cell) {
				t.Errorf("SetCell(%d, %d) then GetCell() = %+v, want %+v", tt.x, tt.y, got, tt.cell)
			}
		})
	}
}

func TestScreenSetCellOutOfBounds(t *testing.T) {
	screen := goterm.NewScreen(80, 24)
	cell := goterm.NewCell('X', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold)

	tests := []struct {
		name string
		x, y int
	}{
		{"negative_x", -1, 0},
		{"negative_y", 0, -1},
		{"x_too_large", 80, 0},
		{"y_too_large", 0, 24},
		{"both_too_large", 100, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SetCell should not panic, just ignore out-of-bounds
			screen.SetCell(tt.x, tt.y, cell)
		})
	}
}

func TestScreenGetCellOutOfBounds(t *testing.T) {
	screen := goterm.NewScreen(80, 24)

	tests := []struct {
		name string
		x, y int
	}{
		{"negative_x", -1, 0},
		{"negative_y", 0, -1},
		{"x_too_large", 80, 0},
		{"y_too_large", 0, 24},
		{"both_too_large", 100, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// GetCell should return a default empty cell for out-of-bounds
			cell := screen.GetCell(tt.x, tt.y)
			expected := goterm.NewCell(' ', goterm.ColorDefault(), goterm.ColorDefault(), goterm.StyleNone)
			if !cell.Equal(expected) {
				t.Errorf("GetCell(%d, %d) out of bounds = %+v, want default cell %+v", tt.x, tt.y, cell, expected)
			}
		})
	}
}

// T022: Unit tests for Clear() functionality
func TestScreenClear(t *testing.T) {
	screen := goterm.NewScreen(80, 24)

	// Fill screen with non-default content
	testCell := goterm.NewCell('X', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold)
	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			screen.SetCell(x, y, testCell)
		}
	}

	// Verify screen is filled
	cell := screen.GetCell(40, 12)
	if !cell.Equal(testCell) {
		t.Fatal("Screen was not filled before Clear() test")
	}

	// Clear the screen
	screen.Clear()

	// Verify all cells are reset to defaults
	defaultCell := goterm.NewCell(' ', goterm.ColorDefault(), goterm.ColorDefault(), goterm.StyleNone)
	for y := 0; y < 24; y++ {
		for x := 0; x < 80; x++ {
			cell := screen.GetCell(x, y)
			if !cell.Equal(defaultCell) {
				t.Errorf("After Clear(), cell at (%d, %d) = %+v, want %+v", x, y, cell, defaultCell)
			}
		}
	}
}

func TestScreenClearPreservesSize(t *testing.T) {
	screen := goterm.NewScreen(100, 30)

	w1, h1 := screen.Size()
	screen.Clear()
	w2, h2 := screen.Size()

	if w1 != w2 || h1 != h2 {
		t.Errorf("Clear() changed size from (%d, %d) to (%d, %d)", w1, h1, w2, h2)
	}
}

func TestScreenMultipleClearCalls(t *testing.T) {
	screen := goterm.NewScreen(50, 20)

	// Multiple clears should be safe
	screen.Clear()
	screen.Clear()
	screen.Clear()

	// Verify screen is still in valid state
	w, h := screen.Size()
	if w != 50 || h != 20 {
		t.Errorf("Multiple Clear() calls corrupted size: got (%d, %d), want (50, 20)", w, h)
	}

	defaultCell := goterm.NewCell(' ', goterm.ColorDefault(), goterm.ColorDefault(), goterm.StyleNone)
	cell := screen.GetCell(25, 10)
	if !cell.Equal(defaultCell) {
		t.Errorf("Multiple Clear() calls left non-default cell: %+v", cell)
	}
}

func TestScreenDrawText(t *testing.T) {
	screen := goterm.NewScreen(80, 24)

	tests := []struct {
		name   string
		x, y   int
		text   string
		fg, bg goterm.Color
		style  goterm.Style
	}{
		{
			name: "simple_text",
			x:    0, y: 0,
			text: "Hello",
			fg:   goterm.ColorRed, bg: goterm.ColorDefault(),
			style: goterm.StyleNone,
		},
		{
			name: "styled_text",
			x:    10, y: 5,
			text: "Bold Text",
			fg:   goterm.ColorGreen, bg: goterm.ColorBlue,
			style: goterm.StyleBold,
		},
		{
			name: "unicode_text",
			x:    20, y: 10,
			text: "Hello 世界 🎮",
			fg:   goterm.ColorCyan, bg: goterm.ColorDefault(),
			style: goterm.StyleNone,
		},
		{
			name: "empty_string",
			x:    30, y: 15,
			text: "",
			fg:   goterm.ColorWhite, bg: goterm.ColorBlack,
			style: goterm.StyleNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen.Clear()
			screen.DrawText(tt.x, tt.y, tt.text, tt.fg, tt.bg, tt.style)

			// Verify each character was placed correctly
			runes := []rune(tt.text)
			for i, ch := range runes {
				cell := screen.GetCell(tt.x+i, tt.y)
				expected := goterm.NewCell(ch, tt.fg, tt.bg, tt.style)
				if !cell.Equal(expected) {
					t.Errorf("DrawText() char %d at (%d, %d) = %+v, want %+v",
						i, tt.x+i, tt.y, cell, expected)
				}
			}
		})
	}
}

func TestScreenDrawTextWrapping(t *testing.T) {
	screen := goterm.NewScreen(80, 24)

	// Draw text that would exceed screen width
	longText := "This is a very long text that definitely exceeds the screen width and should be clipped"
	screen.DrawText(70, 10, longText, goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)

	// Characters beyond screen width should not cause issues
	// GetCell at x=80 should return default (out of bounds)
	cell := screen.GetCell(80, 10)
	defaultCell := goterm.NewCell(' ', goterm.ColorDefault(), goterm.ColorDefault(), goterm.StyleNone)
	if !cell.Equal(defaultCell) {
		t.Errorf("DrawText() wrote beyond screen bounds")
	}

	// Last visible character should be at x=79
	cell = screen.GetCell(79, 10)
	if cell.Ch == ' ' && cell.Fg == goterm.ColorDefault() {
		t.Error("DrawText() didn't write to last column as expected")
	}
}

func TestScreenResize(t *testing.T) {
	screen := goterm.NewScreen(80, 24)

	// Fill with content
	testCell := goterm.NewCell('A', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold)
	screen.SetCell(10, 10, testCell)

	// Resize to larger
	screen.Resize(100, 30)
	w, h := screen.Size()
	if w != 100 || h != 30 {
		t.Errorf("Resize(100, 30) resulted in size (%d, %d)", w, h)
	}

	// Verify content is preserved (where it fits)
	cell := screen.GetCell(10, 10)
	if !cell.Equal(testCell) {
		t.Error("Resize() lost existing content")
	}

	// Resize to smaller
	screen.Resize(40, 12)
	w, h = screen.Size()
	if w != 40 || h != 12 {
		t.Errorf("Resize(40, 12) resulted in size (%d, %d)", w, h)
	}

	// Content outside new bounds should not be accessible
	cell = screen.GetCell(50, 15)
	defaultCell := goterm.NewCell(' ', goterm.ColorDefault(), goterm.ColorDefault(), goterm.StyleNone)
	if !cell.Equal(defaultCell) {
		t.Error("Resize() didn't properly handle out-of-bounds access")
	}
}
