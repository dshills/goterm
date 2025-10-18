package unit

import (
	"testing"

	"github.com/dshills/goterm"
)

func TestColorRGB(t *testing.T) {
	tests := []struct {
		name     string
		r, g, b  uint8
		wantMode goterm.ColorMode
	}{
		{"red", 255, 0, 0, goterm.ColorModeTrueColor},
		{"green", 0, 255, 0, goterm.ColorModeTrueColor},
		{"blue", 0, 0, 255, goterm.ColorModeTrueColor},
		{"white", 255, 255, 255, goterm.ColorModeTrueColor},
		{"black", 0, 0, 0, goterm.ColorModeTrueColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := goterm.ColorRGB(tt.r, tt.g, tt.b)
			if c.Mode() != tt.wantMode {
				t.Errorf("ColorRGB().Mode() = %v, want %v", c.Mode(), tt.wantMode)
			}
			gotR, gotG, gotB := c.RGB()
			if gotR != tt.r || gotG != tt.g || gotB != tt.b {
				t.Errorf("ColorRGB().RGB() = (%d, %d, %d), want (%d, %d, %d)",
					gotR, gotG, gotB, tt.r, tt.g, tt.b)
			}
		})
	}
}

func TestColorIndex(t *testing.T) {
	tests := []struct {
		name     string
		index    uint8
		wantMode goterm.ColorMode
	}{
		{"black", 0, goterm.ColorMode16},
		{"red", 1, goterm.ColorMode16},
		{"bright_black", 8, goterm.ColorMode16},
		{"bright_white", 15, goterm.ColorMode16},
		{"color_cube_start", 16, goterm.ColorMode256},
		{"color_cube_mid", 100, goterm.ColorMode256},
		{"grayscale_start", 232, goterm.ColorMode256},
		{"grayscale_end", 255, goterm.ColorMode256},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := goterm.ColorIndex(tt.index)
			if c.Mode() != tt.wantMode {
				t.Errorf("ColorIndex(%d).Mode() = %v, want %v", tt.index, c.Mode(), tt.wantMode)
			}
			if c.Index() != tt.index {
				t.Errorf("ColorIndex(%d).Index() = %d, want %d", tt.index, c.Index(), tt.index)
			}
		})
	}
}

func TestColorDefault(t *testing.T) {
	c := goterm.ColorDefault()
	if c.Mode() != goterm.ColorModeDefault {
		t.Errorf("ColorDefault().Mode() = %v, want %v", c.Mode(), goterm.ColorModeDefault)
	}
}

func TestColorTo256(t *testing.T) {
	// Test RGB to 256-color conversion
	tests := []struct {
		name    string
		r, g, b uint8
	}{
		{"red", 255, 0, 0},
		{"green", 0, 255, 0},
		{"blue", 0, 0, 255},
		{"gray", 128, 128, 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgb := goterm.ColorRGB(tt.r, tt.g, tt.b)
			c256 := rgb.To256()
			if c256.Mode() != goterm.ColorMode256 {
				t.Errorf("To256().Mode() = %v, want %v", c256.Mode(), goterm.ColorMode256)
			}
			// Index should be in valid range (16-255 for color cube/grayscale)
			if idx := c256.Index(); idx < 16 {
				t.Errorf("To256().Index() = %d, want >= 16", idx)
			}
		})
	}
}

func TestColorTo16(t *testing.T) {
	// Test conversion to 16-color
	tests := []struct {
		name  string
		color goterm.Color
	}{
		{"from_256", goterm.ColorIndex(45)},
		{"from_rgb", goterm.ColorRGB(200, 100, 50)},
		{"from_16", goterm.ColorIndex(5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c16 := tt.color.To16()
			if c16.Mode() != goterm.ColorMode16 {
				t.Errorf("To16().Mode() = %v, want %v", c16.Mode(), goterm.ColorMode16)
			}
			if idx := c16.Index(); idx >= 16 {
				t.Errorf("To16().Index() = %d, want < 16", idx)
			}
		})
	}
}

func TestNamedColors(t *testing.T) {
	// Test that named colors are correctly defined
	colors := map[string]goterm.Color{
		"black":   goterm.ColorBlack,
		"red":     goterm.ColorRed,
		"green":   goterm.ColorGreen,
		"yellow":  goterm.ColorYellow,
		"blue":    goterm.ColorBlue,
		"magenta": goterm.ColorMagenta,
		"cyan":    goterm.ColorCyan,
		"white":   goterm.ColorWhite,
	}

	for name, color := range colors {
		t.Run(name, func(t *testing.T) {
			if color.Mode() != goterm.ColorMode16 {
				t.Errorf("%s color mode = %v, want %v", name, color.Mode(), goterm.ColorMode16)
			}
			if color.Index() >= 8 {
				t.Errorf("%s color index = %d, want < 8", name, color.Index())
			}
		})
	}
}
