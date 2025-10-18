package goterm

import "testing"

// Basic sanity tests to ensure coverage reporting works
// More comprehensive tests are in tests/unit/ directory

func TestColorDefaultBasic(t *testing.T) {
	c := ColorDefault()
	if c.Mode() != ColorModeDefault {
		t.Errorf("ColorDefault().Mode() = %v, want %v", c.Mode(), ColorModeDefault)
	}
}

func TestNewCellBasic(t *testing.T) {
	cell := NewCell('A', ColorRed, ColorBlue, StyleBold)
	if cell.Ch != 'A' {
		t.Errorf("NewCell().Ch = %v, want 'A'", cell.Ch)
	}
}

func TestNewScreenBasic(t *testing.T) {
	screen := NewScreen(80, 24)
	w, h := screen.Size()
	if w != 80 || h != 24 {
		t.Errorf("NewScreen(80, 24).Size() = (%d, %d), want (80, 24)", w, h)
	}
}
