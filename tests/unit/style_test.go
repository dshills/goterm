package unit

import (
	"testing"

	"github.com/dshills/goterm"
)

func TestStyleBitmask(t *testing.T) {
	// Test basic style operations
	s := goterm.StyleNone

	// Test Set
	s = s.Set(goterm.StyleBold)
	if !s.Has(goterm.StyleBold) {
		t.Error("Set(StyleBold) failed, Has(StyleBold) returned false")
	}

	// Test combining multiple styles
	s = s.Set(goterm.StyleUnderline)
	if !s.Has(goterm.StyleBold) || !s.Has(goterm.StyleUnderline) {
		t.Error("Multiple styles not working correctly")
	}

	// Test Clear
	s = s.Clear(goterm.StyleBold)
	if s.Has(goterm.StyleBold) {
		t.Error("Clear(StyleBold) failed, Has(StyleBold) still returns true")
	}
	if !s.Has(goterm.StyleUnderline) {
		t.Error("Clear affected wrong style")
	}

	// Test Toggle
	s = goterm.StyleNone
	s = s.Toggle(goterm.StyleItalic)
	if !s.Has(goterm.StyleItalic) {
		t.Error("Toggle on failed")
	}
	s = s.Toggle(goterm.StyleItalic)
	if s.Has(goterm.StyleItalic) {
		t.Error("Toggle off failed")
	}
}

func TestStyleCombinations(t *testing.T) {
	tests := []struct {
		name   string
		styles []goterm.Style
	}{
		{"bold_underline", []goterm.Style{goterm.StyleBold, goterm.StyleUnderline}},
		{"italic_strikethrough", []goterm.Style{goterm.StyleItalic, goterm.StyleStrikethrough}},
		{"all_styles", []goterm.Style{
			goterm.StyleBold,
			goterm.StyleDim,
			goterm.StyleItalic,
			goterm.StyleUnderline,
			goterm.StyleSlowBlink,
			goterm.StyleRapidBlink,
			goterm.StyleReverse,
			goterm.StyleConceal,
			goterm.StyleStrikethrough,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := goterm.StyleNone
			for _, style := range tt.styles {
				s = s.Set(style)
			}

			// Verify all styles are set
			for _, style := range tt.styles {
				if !s.Has(style) {
					t.Errorf("Style %v not set correctly", style)
				}
			}
		})
	}
}

func TestStyleBitwiseOR(t *testing.T) {
	// Test combining styles with bitwise OR
	combined := goterm.StyleBold | goterm.StyleUnderline | goterm.StyleItalic

	if !combined.Has(goterm.StyleBold) {
		t.Error("Combined style missing StyleBold")
	}
	if !combined.Has(goterm.StyleUnderline) {
		t.Error("Combined style missing StyleUnderline")
	}
	if !combined.Has(goterm.StyleItalic) {
		t.Error("Combined style missing StyleItalic")
	}
	if combined.Has(goterm.StyleReverse) {
		t.Error("Combined style has unexpected StyleReverse")
	}
}

func TestStyleNone(t *testing.T) {
	s := goterm.StyleNone

	styles := []goterm.Style{
		goterm.StyleBold,
		goterm.StyleDim,
		goterm.StyleItalic,
		goterm.StyleUnderline,
		goterm.StyleSlowBlink,
		goterm.StyleRapidBlink,
		goterm.StyleReverse,
		goterm.StyleConceal,
		goterm.StyleStrikethrough,
	}

	for _, style := range styles {
		if s.Has(style) {
			t.Errorf("StyleNone should not have %v", style)
		}
	}
}
