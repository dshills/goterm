package goterm

import "fmt"

// ColorMode represents the color capability mode
type ColorMode int

const (
	// ColorModeDefault uses terminal's default colors
	ColorModeDefault ColorMode = iota
	// ColorMode16 uses 16 ANSI colors
	ColorMode16
	// ColorMode256 uses 256-color palette
	ColorMode256
	// ColorModeTrueColor uses 24-bit RGB colors
	ColorModeTrueColor
)

// Color represents a terminal color with support for multiple modes
type Color struct {
	mode    ColorMode
	r, g, b uint8 // RGB values for truecolor
	index   uint8 // Palette index for 16/256-color modes
}

// ColorDefault returns the terminal's default color
func ColorDefault() Color {
	return Color{mode: ColorModeDefault}
}

// ColorRGB creates a true color (24-bit RGB)
// Automatically degrades to 256-color or 16-color on terminals without truecolor support
func ColorRGB(r, g, b uint8) Color {
	return Color{
		mode: ColorModeTrueColor,
		r:    r,
		g:    g,
		b:    b,
	}
}

// ColorIndex creates an indexed color (0-255)
// Ranges:
//   - 0-7: Basic ANSI colors
//   - 8-15: Bright ANSI colors
//   - 16-231: 216-color cube
//   - 232-255: Grayscale ramp
func ColorIndex(index uint8) Color {
	mode := ColorMode256
	if index < 16 {
		mode = ColorMode16
	}
	return Color{
		mode:  mode,
		index: index,
	}
}

// Named color constants (ANSI colors)
var (
	ColorBlack   = ColorIndex(0)
	ColorRed     = ColorIndex(1)
	ColorGreen   = ColorIndex(2)
	ColorYellow  = ColorIndex(3)
	ColorBlue    = ColorIndex(4)
	ColorMagenta = ColorIndex(5)
	ColorCyan    = ColorIndex(6)
	ColorWhite   = ColorIndex(7)
)

// Mode returns the color's mode
func (c Color) Mode() ColorMode {
	return c.mode
}

// RGB returns the RGB components (only valid for TrueColor mode)
func (c Color) RGB() (r, g, b uint8) {
	return c.r, c.g, c.b
}

// Index returns the palette index (only valid for 16/256-color modes)
func (c Color) Index() uint8 {
	return c.index
}

// To256 converts RGB color to nearest 256-color palette index
func (c Color) To256() Color {
	if c.mode != ColorModeTrueColor {
		return c
	}

	// Convert RGB to 216-color cube (6x6x6)
	// Formula: 16 + 36*r + 6*g + b where r,g,b are in range 0-5
	// r, g, b are in range 0-5, so maximum value is:
	// 16 + 36*5 + 6*5 + 5 = 16 + 180 + 30 + 5 = 231, which fits in uint8
	r := (int(c.r) * 6) / 256
	g := (int(c.g) * 6) / 256
	b := (int(c.b) * 6) / 256

	// Safe conversion: result is always in range [16, 231]
	index := uint8(16 + 36*r + 6*g + b) // #nosec G115
	return ColorIndex(index)
}

// To16 converts color to nearest ANSI 16-color
func (c Color) To16() Color {
	if c.mode == ColorMode16 {
		return c
	}

	// Convert to 256-color first if needed
	col := c
	if c.mode == ColorModeTrueColor {
		col = c.To256()
	}

	// Map 256-color to nearest 16-color
	idx := col.index
	if idx < 16 {
		return col
	}

	// Simplified mapping: map to basic 8 colors
	// This is a basic implementation - a more sophisticated one would
	// calculate actual color distance
	if idx >= 232 { // Grayscale
		if idx < 244 {
			return ColorBlack
		}
		return ColorWhite
	}

	// For color cube, use simple modulo mapping
	return ColorIndex(idx % 8)
}

// ansiCode returns the ANSI escape sequence for this color
// fg=true for foreground, fg=false for background
func (c Color) ansiCode(fg bool) string {
	switch c.mode {
	case ColorModeDefault:
		if fg {
			return "\x1b[39m"
		}
		return "\x1b[49m"

	case ColorMode16:
		// Standard ANSI colors
		if fg {
			if c.index < 8 {
				return fmt.Sprintf("\x1b[%dm", 30+c.index)
			}
			// Bright colors (8-15)
			return fmt.Sprintf("\x1b[%dm", 90+c.index-8)
		}
		// Background
		if c.index < 8 {
			return fmt.Sprintf("\x1b[%dm", 40+c.index)
		}
		return fmt.Sprintf("\x1b[%dm", 100+c.index-8)

	case ColorMode256:
		// 256-color mode
		if fg {
			return fmt.Sprintf("\x1b[38;5;%dm", c.index)
		}
		return fmt.Sprintf("\x1b[48;5;%dm", c.index)

	case ColorModeTrueColor:
		// 24-bit RGB
		if fg {
			return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.r, c.g, c.b)
		}
		return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.r, c.g, c.b)
	}

	return ""
}
