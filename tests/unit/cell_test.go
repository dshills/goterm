package unit

import (
	"testing"

	"github.com/dshills/goterm"
)

func TestNewCell(t *testing.T) {
	ch := 'A'
	fg := goterm.ColorRed
	bg := goterm.ColorBlue
	style := goterm.StyleBold

	cell := goterm.NewCell(ch, fg, bg, style)

	if cell.Ch != ch {
		t.Errorf("NewCell().Ch = %v, want %v", cell.Ch, ch)
	}
	if cell.Fg != fg {
		t.Errorf("NewCell().Fg = %v, want %v", cell.Fg, fg)
	}
	if cell.Bg != bg {
		t.Errorf("NewCell().Bg = %v, want %v", cell.Bg, bg)
	}
	if cell.Style != style {
		t.Errorf("NewCell().Style = %v, want %v", cell.Style, style)
	}
}

func TestCellClear(t *testing.T) {
	// Create a cell with non-default values
	cell := goterm.NewCell('X', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold)

	// Clear it
	cell.Clear()

	// Verify it's reset to defaults
	if cell.Ch != ' ' {
		t.Errorf("Clear() did not reset Ch, got %v, want ' '", cell.Ch)
	}
	if cell.Fg != goterm.ColorDefault() {
		t.Errorf("Clear() did not reset Fg to default")
	}
	if cell.Bg != goterm.ColorDefault() {
		t.Errorf("Clear() did not reset Bg to default")
	}
	if cell.Style != goterm.StyleNone {
		t.Errorf("Clear() did not reset Style, got %v, want %v", cell.Style, goterm.StyleNone)
	}
}

func TestCellEqual(t *testing.T) {
	tests := []struct {
		name  string
		cell1 goterm.Cell
		cell2 goterm.Cell
		want  bool
	}{
		{
			name:  "identical cells",
			cell1: goterm.NewCell('A', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold),
			cell2: goterm.NewCell('A', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold),
			want:  true,
		},
		{
			name:  "different character",
			cell1: goterm.NewCell('A', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold),
			cell2: goterm.NewCell('B', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold),
			want:  false,
		},
		{
			name:  "different foreground",
			cell1: goterm.NewCell('A', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold),
			cell2: goterm.NewCell('A', goterm.ColorGreen, goterm.ColorBlue, goterm.StyleBold),
			want:  false,
		},
		{
			name:  "different background",
			cell1: goterm.NewCell('A', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold),
			cell2: goterm.NewCell('A', goterm.ColorRed, goterm.ColorGreen, goterm.StyleBold),
			want:  false,
		},
		{
			name:  "different style",
			cell1: goterm.NewCell('A', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold),
			cell2: goterm.NewCell('A', goterm.ColorRed, goterm.ColorBlue, goterm.StyleItalic),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cell1.Equal(tt.cell2)
			if got != tt.want {
				t.Errorf("Cell.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCellUnicode(t *testing.T) {
	// Test with various Unicode characters
	tests := []rune{
		'A', // ASCII
		'é', // Latin with diacritic
		'日', // CJK (double-width)
		'🎮', // Emoji
		'→', // Arrow
	}

	for _, ch := range tests {
		cell := goterm.NewCell(ch, goterm.ColorDefault(), goterm.ColorDefault(), goterm.StyleNone)
		if cell.Ch != ch {
			t.Errorf("Cell with rune %v has Ch = %v", ch, cell.Ch)
		}
	}
}
