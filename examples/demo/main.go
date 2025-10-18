package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dshills/goterm"
)

func main() {
	// Initialize the terminal
	screen, err := goterm.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize terminal: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := screen.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close screen: %v\n", err)
		}
	}()

	// Run the demo sequence
	demos := []struct {
		name string
		fn   func(*goterm.Screen)
	}{
		{"Welcome Screen", demoWelcome},
		{"16 ANSI Colors", demo16Colors},
		{"256 Color Palette", demo256Colors},
		{"True Color (RGB)", demoTrueColor},
		{"Color Conversions", demoColorConversions},
		{"Text Styles", demoStyles},
		{"Style Combinations", demoStyleCombinations},
		{"Unicode Support", demoUnicode},
		{"Box Drawing", demoBoxDrawing},
		{"Screen Buffer", demoScreenBuffer},
		{"Cell Manipulation", demoCellManipulation},
		{"Dynamic Resize", demoResize},
		{"Color Gradients", demoGradients},
		{"Animation", demoAnimation},
		{"Final Showcase", demoFinalShowcase},
	}

	for i, demo := range demos {
		screen.Clear()

		// Draw demo header
		drawHeader(screen, fmt.Sprintf("Demo %d/%d: %s", i+1, len(demos), demo.name))

		// Run the demo
		demo.fn(screen)

		// Draw footer
		drawFooter(screen, "Press Ctrl+C to exit, waiting 3 seconds...")

		if err := screen.Show(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to show screen: %v\n", err)
			return
		}

		time.Sleep(3 * time.Second)
	}

	// Final goodbye
	screen.Clear()
	w, h := screen.Size()
	msg := "Thanks for watching the goterm demo!"
	screen.DrawText((w-len(msg))/2, h/2, msg, goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleBold)
	if err := screen.Show(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to show screen: %v\n", err)
	}
	time.Sleep(2 * time.Second)
}

func drawHeader(screen *goterm.Screen, title string) {
	w, _ := screen.Size()

	// Top border
	border := "═"
	for i := 0; i < w; i++ {
		screen.SetCell(i, 0, goterm.NewCell(rune(border[0]), goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold))
	}

	// Title (centered)
	titleText := fmt.Sprintf(" %s ", title)
	startX := (w - len(titleText)) / 2
	screen.DrawText(startX, 1, titleText, goterm.ColorWhite, goterm.ColorBlue, goterm.StyleBold)

	// Bottom border of header
	for i := 0; i < w; i++ {
		screen.SetCell(i, 2, goterm.NewCell(rune(border[0]), goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold))
	}
}

func drawFooter(screen *goterm.Screen, message string) {
	w, h := screen.Size()

	// Bottom border
	for i := 0; i < w; i++ {
		screen.SetCell(i, h-2, goterm.NewCell('─', goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone))
	}

	// Footer message
	screen.DrawText(2, h-1, message, goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleItalic)
}

func demoWelcome(screen *goterm.Screen) {
	w, h := screen.Size()

	lines := []string{
		"Welcome to the goterm Library Demo!",
		"",
		"This demo showcases:",
		"  • Multi-mode color support (16, 256, RGB)",
		"  • Text styling (bold, italic, underline, etc.)",
		"  • Unicode and CJK character support",
		"  • Screen buffer management",
		"  • Cell-based rendering",
		"",
		"Terminal Size: " + fmt.Sprintf("%dx%d", w, h),
	}

	startY := (h - len(lines)) / 2
	for i, line := range lines {
		startX := (w - len(line)) / 2
		color := goterm.ColorGreen
		style := goterm.StyleNone

		if i == 0 {
			color = goterm.ColorCyan
			style = goterm.StyleBold
		} else if i >= 2 && i <= 6 {
			color = goterm.ColorYellow
		}

		screen.DrawText(startX, startY+i, line, color, goterm.ColorDefault(), style)
	}
}

func demo16Colors(screen *goterm.Screen) {
	w, _ := screen.Size()

	screen.DrawText(4, 4, "16 ANSI Colors (Foreground):", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

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

	y := 6
	for i, c := range colors {
		// Standard color
		screen.DrawText(4, y+i, fmt.Sprintf("%-10s", c.name), c.color, goterm.ColorDefault(), goterm.StyleNone)
		screen.DrawText(16, y+i, "█████", c.color, goterm.ColorDefault(), goterm.StyleBold)

		// Bright variant (index + 8)
		brightColor := goterm.ColorIndex(uint8(i + 8))
		screen.DrawText(w/2, y+i, fmt.Sprintf("Bright %-10s", c.name), brightColor, goterm.ColorDefault(), goterm.StyleNone)
		screen.DrawText(w/2+16, y+i, "█████", brightColor, goterm.ColorDefault(), goterm.StyleBold)
	}

	// Background colors demo
	screen.DrawText(4, y+9, "Background Colors:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
	y += 11
	for i, c := range colors {
		screen.DrawText(4, y+i, fmt.Sprintf("  %-8s  ", c.name), goterm.ColorWhite, c.color, goterm.StyleBold)
	}
}

func demo256Colors(screen *goterm.Screen) {
	w, _ := screen.Size()

	screen.DrawText(4, 4, "256 Color Palette:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	// Standard colors (0-15)
	screen.DrawText(4, 6, "Standard (0-15):", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	for i := 0; i < 16; i++ {
		x := 4 + (i * 3)
		screen.DrawText(x, 7, "██", goterm.ColorIndex(uint8(i)), goterm.ColorDefault(), goterm.StyleNone)
	}

	// 216 color cube (16-231)
	screen.DrawText(4, 9, "Color Cube (16-231):", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	y := 10
	for row := 0; row < 6; row++ {
		for g := 0; g < 6; g++ {
			for r := 0; r < 6; r++ {
				idx := uint8(16 + 36*r + 6*g + row)
				x := 4 + (r+g*6)*2
				if x < w-2 {
					screen.DrawText(x, y+g, "█", goterm.ColorIndex(idx), goterm.ColorDefault(), goterm.StyleNone)
				}
			}
		}
		y += 7
	}

	// Grayscale (232-255)
	screen.DrawText(4, y, "Grayscale (232-255):", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	y++
	for i := 0; i < 24; i++ {
		x := 4 + (i * 2)
		if x < w-2 {
			screen.DrawText(x, y, "█", goterm.ColorIndex(uint8(232+i)), goterm.ColorDefault(), goterm.StyleNone)
		}
	}
}

func demoTrueColor(screen *goterm.Screen) {
	w, _ := screen.Size()

	screen.DrawText(4, 4, "24-bit True Color (RGB):", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	// RGB color examples
	examples := []struct {
		r, g, b uint8
		name    string
	}{
		{255, 0, 0, "Pure Red (255,0,0)"},
		{0, 255, 0, "Pure Green (0,255,0)"},
		{0, 0, 255, "Pure Blue (0,0,255)"},
		{255, 255, 0, "Yellow (255,255,0)"},
		{255, 0, 255, "Magenta (255,0,255)"},
		{0, 255, 255, "Cyan (0,255,255)"},
		{255, 128, 0, "Orange (255,128,0)"},
		{128, 0, 128, "Purple (128,0,128)"},
		{0, 128, 128, "Teal (0,128,128)"},
		{255, 192, 203, "Pink (255,192,203)"},
	}

	y := 6
	for i, ex := range examples {
		color := goterm.ColorRGB(ex.r, ex.g, ex.b)
		screen.DrawText(4, y+i, "████████", color, goterm.ColorDefault(), goterm.StyleNone)
		screen.DrawText(14, y+i, ex.name, goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	}

	// RGB gradient
	screen.DrawText(4, y+11, "RGB Gradients:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	// Red gradient
	y += 13
	for i := 0; i < 32 && (4+i*2) < w-2; i++ {
		r := uint8((i * 255) / 31)
		color := goterm.ColorRGB(r, 0, 0)
		screen.DrawText(4+i*2, y, "█", color, goterm.ColorDefault(), goterm.StyleNone)
	}

	// Green gradient
	y++
	for i := 0; i < 32 && (4+i*2) < w-2; i++ {
		g := uint8((i * 255) / 31)
		color := goterm.ColorRGB(0, g, 0)
		screen.DrawText(4+i*2, y, "█", color, goterm.ColorDefault(), goterm.StyleNone)
	}

	// Blue gradient
	y++
	for i := 0; i < 32 && (4+i*2) < w-2; i++ {
		b := uint8((i * 255) / 31)
		color := goterm.ColorRGB(0, 0, b)
		screen.DrawText(4+i*2, y, "█", color, goterm.ColorDefault(), goterm.StyleNone)
	}
}

func demoColorConversions(screen *goterm.Screen) {
	screen.DrawText(4, 4, "Color Mode Conversions:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	// Example: RGB -> 256 -> 16
	rgb := goterm.ColorRGB(128, 200, 255)
	c256 := rgb.To256()
	c16 := rgb.To16()

	y := 6
	screen.DrawText(4, y, "Original RGB (128,200,255):", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(34, y, "████████", rgb, goterm.ColorDefault(), goterm.StyleNone)

	y++
	screen.DrawText(4, y, "Converted to 256-color:", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(34, y, "████████", c256, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(44, y, fmt.Sprintf("(index: %d)", c256.Index()), goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)

	y++
	screen.DrawText(4, y, "Converted to 16-color:", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(34, y, "████████", c16, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(44, y, fmt.Sprintf("(index: %d)", c16.Index()), goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)

	// More examples
	y += 3
	screen.DrawText(4, y, "More Conversion Examples:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	examples := []struct {
		r, g, b uint8
		name    string
	}{
		{255, 0, 0, "Red"},
		{0, 255, 0, "Green"},
		{255, 128, 64, "Orange"},
		{100, 100, 255, "Light Blue"},
	}

	y += 2
	screen.DrawText(6, y, "RGB", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold)
	screen.DrawText(20, y, "256", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold)
	screen.DrawText(34, y, "16", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold)
	screen.DrawText(48, y, "Name", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold)

	y++
	for i, ex := range examples {
		rgb := goterm.ColorRGB(ex.r, ex.g, ex.b)
		screen.DrawText(6, y+i, "████", rgb, goterm.ColorDefault(), goterm.StyleNone)
		screen.DrawText(20, y+i, "████", rgb.To256(), goterm.ColorDefault(), goterm.StyleNone)
		screen.DrawText(34, y+i, "████", rgb.To16(), goterm.ColorDefault(), goterm.StyleNone)
		screen.DrawText(48, y+i, ex.name, goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	}
}

func demoStyles(screen *goterm.Screen) {
	screen.DrawText(4, 4, "Text Styles:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	styles := []struct {
		style goterm.Style
		name  string
	}{
		{goterm.StyleNone, "Normal (no style)"},
		{goterm.StyleBold, "Bold"},
		{goterm.StyleDim, "Dim (faint)"},
		{goterm.StyleItalic, "Italic"},
		{goterm.StyleUnderline, "Underline"},
		{goterm.StyleSlowBlink, "Slow Blink"},
		{goterm.StyleRapidBlink, "Rapid Blink"},
		{goterm.StyleReverse, "Reverse (inverted)"},
		{goterm.StyleConceal, "Conceal (hidden)"},
		{goterm.StyleStrikethrough, "Strikethrough"},
	}

	y := 6
	for i, s := range styles {
		// Style name
		screen.DrawText(4, y+i, fmt.Sprintf("%-25s", s.name), goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)

		// Example text with style
		screen.DrawText(30, y+i, "The quick brown fox jumps over the lazy dog", goterm.ColorWhite, goterm.ColorDefault(), s.style)
	}

	// Note about terminal support
	y += len(styles) + 1
	screen.DrawText(4, y, "Note: Some styles may not be supported by all terminals", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleItalic)
}

func demoStyleCombinations(screen *goterm.Screen) {
	screen.DrawText(4, 4, "Style Combinations (using bitwise OR):", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	combinations := []struct {
		style goterm.Style
		name  string
	}{
		{goterm.StyleBold | goterm.StyleUnderline, "Bold + Underline"},
		{goterm.StyleBold | goterm.StyleItalic, "Bold + Italic"},
		{goterm.StyleItalic | goterm.StyleUnderline, "Italic + Underline"},
		{goterm.StyleBold | goterm.StyleItalic | goterm.StyleUnderline, "Bold + Italic + Underline"},
		{goterm.StyleDim | goterm.StyleItalic, "Dim + Italic"},
		{goterm.StyleBold | goterm.StyleReverse, "Bold + Reverse"},
		{goterm.StyleUnderline | goterm.StyleSlowBlink, "Underline + Blink"},
	}

	y := 6
	for i, c := range combinations {
		screen.DrawText(4, y+i, fmt.Sprintf("%-30s", c.name), goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
		screen.DrawText(36, y+i, "The quick brown fox", goterm.ColorWhite, goterm.ColorDefault(), c.style)
	}

	// Demonstrate style operations
	y += len(combinations) + 2
	screen.DrawText(4, y, "Style Operations:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	baseStyle := goterm.StyleBold
	y += 2
	screen.DrawText(4, y, "Base: Bold", goterm.ColorWhite, goterm.ColorDefault(), baseStyle)

	y++
	withUnderline := baseStyle.Set(goterm.StyleUnderline)
	screen.DrawText(4, y, "Set(Underline): Bold + Underline", goterm.ColorWhite, goterm.ColorDefault(), withUnderline)

	y++
	toggled := withUnderline.Toggle(goterm.StyleBold)
	screen.DrawText(4, y, "Toggle(Bold): Just Underline", goterm.ColorWhite, goterm.ColorDefault(), toggled)

	y++
	cleared := withUnderline.Clear(goterm.StyleUnderline)
	screen.DrawText(4, y, "Clear(Underline): Just Bold", goterm.ColorWhite, goterm.ColorDefault(), cleared)
}

func demoUnicode(screen *goterm.Screen) {
	screen.DrawText(4, 4, "Unicode Character Support:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	examples := []struct {
		text  string
		label string
		color goterm.Color
	}{
		{"ASCII: The quick brown fox jumps over the lazy dog", "ASCII", goterm.ColorWhite},
		{"Latin: Héllö Wörld! Ñoño Çédille", "Latin Extended", goterm.ColorYellow},
		{"Greek: Γεια σου κόσμε (Hello World)", "Greek", goterm.ColorCyan},
		{"Cyrillic: Привет мир (Hello World)", "Cyrillic", goterm.ColorGreen},
		{"Arabic: مرحبا بالعالم (Hello World)", "Arabic", goterm.ColorMagenta},
		{"Hebrew: שלום עולם (Hello World)", "Hebrew", goterm.ColorBlue},
		{"CJK Japanese: こんにちは世界 (Hello World)", "Japanese", goterm.ColorRed},
		{"CJK Chinese: 你好世界 (Hello World)", "Chinese", goterm.ColorYellow},
		{"CJK Korean: 안녕하세요 세계 (Hello World)", "Korean", goterm.ColorCyan},
		{"Emoji: 🎮 🚀 🌟 ⚡ 🔥 💻 🎨 🎭 🎪 🎯", "Emoji", goterm.ColorMagenta},
		{"Math: ∑ ∫ √ ≠ ≤ ≥ ∞ π θ α β γ", "Math Symbols", goterm.ColorGreen},
		{"Currency: $ € £ ¥ ₹ ₽ ₿ ¢", "Currency", goterm.ColorYellow},
		{"Arrows: ← → ↑ ↓ ↔ ↕ ⇐ ⇒ ⇑ ⇓", "Arrows", goterm.ColorCyan},
		{"Stars: ★ ☆ ✦ ✧ ✯ ✪ ⭐ 🌟", "Stars", goterm.ColorYellow},
	}

	y := 6
	for i, ex := range examples {
		screen.DrawText(4, y+i, fmt.Sprintf("%-18s", ex.label+":"), goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
		screen.DrawText(24, y+i, ex.text, ex.color, goterm.ColorDefault(), goterm.StyleNone)
	}
}

func demoBoxDrawing(screen *goterm.Screen) {
	w, _ := screen.Size()

	screen.DrawText(4, 4, "Box Drawing Characters:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	// Single-line box
	y := 6
	screen.DrawText(4, y, "Single Line:", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	y++
	drawBox(screen, 4, y, 30, 5, goterm.ColorCyan, true)
	screen.DrawText(6, y+2, "Single-line box", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)

	// Double-line box
	screen.DrawText(40, 6, "Double Line:", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	drawBox(screen, 40, 7, 30, 5, goterm.ColorMagenta, false)
	screen.DrawText(42, 9, "Double-line box", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)

	// Nested boxes
	y = 14
	screen.DrawText(4, y, "Nested Boxes:", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	y++
	drawBox(screen, 4, y, 60, 8, goterm.ColorGreen, true)
	drawBox(screen, 8, y+2, 20, 4, goterm.ColorYellow, true)
	screen.DrawText(10, y+3, "Inner box", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	drawBox(screen, 32, y+2, 20, 4, goterm.ColorCyan, false)
	screen.DrawText(34, y+3, "Another inner", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)

	// Complex layout
	if w > 80 {
		y = 24
		screen.DrawText(4, y, "Complex Layout:", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
		y++

		// Main container
		drawBox(screen, 4, y, 70, 10, goterm.ColorWhite, true)

		// Title bar
		screen.DrawText(6, y+1, "┤ Application Title ├", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold)

		// Divider
		for i := 5; i < 73; i++ {
			screen.SetCell(i, y+2, goterm.NewCell('─', goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone))
		}
		screen.SetCell(4, y+2, goterm.NewCell('├', goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone))
		screen.SetCell(73, y+2, goterm.NewCell('┤', goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone))

		// Content
		screen.DrawText(6, y+4, "Content Area", goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleNone)
		screen.DrawText(6, y+5, "  • Item 1", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
		screen.DrawText(6, y+6, "  • Item 2", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)

		// Status bar
		for i := 5; i < 73; i++ {
			screen.SetCell(i, y+8, goterm.NewCell('─', goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone))
		}
		screen.SetCell(4, y+8, goterm.NewCell('├', goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone))
		screen.SetCell(73, y+8, goterm.NewCell('┤', goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone))
		screen.DrawText(6, y+9, "Status: Ready", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	}
}

func drawBox(screen *goterm.Screen, x, y, width, height int, color goterm.Color, singleLine bool) {
	var topLeft, topRight, bottomLeft, bottomRight, horizontal, vertical rune

	if singleLine {
		topLeft, topRight = '┌', '┐'
		bottomLeft, bottomRight = '└', '┘'
		horizontal, vertical = '─', '│'
	} else {
		topLeft, topRight = '╔', '╗'
		bottomLeft, bottomRight = '╚', '╝'
		horizontal, vertical = '═', '║'
	}

	// Corners
	screen.SetCell(x, y, goterm.NewCell(topLeft, color, goterm.ColorDefault(), goterm.StyleNone))
	screen.SetCell(x+width-1, y, goterm.NewCell(topRight, color, goterm.ColorDefault(), goterm.StyleNone))
	screen.SetCell(x, y+height-1, goterm.NewCell(bottomLeft, color, goterm.ColorDefault(), goterm.StyleNone))
	screen.SetCell(x+width-1, y+height-1, goterm.NewCell(bottomRight, color, goterm.ColorDefault(), goterm.StyleNone))

	// Horizontal lines
	for i := 1; i < width-1; i++ {
		screen.SetCell(x+i, y, goterm.NewCell(horizontal, color, goterm.ColorDefault(), goterm.StyleNone))
		screen.SetCell(x+i, y+height-1, goterm.NewCell(horizontal, color, goterm.ColorDefault(), goterm.StyleNone))
	}

	// Vertical lines
	for i := 1; i < height-1; i++ {
		screen.SetCell(x, y+i, goterm.NewCell(vertical, color, goterm.ColorDefault(), goterm.StyleNone))
		screen.SetCell(x+width-1, y+i, goterm.NewCell(vertical, color, goterm.ColorDefault(), goterm.StyleNone))
	}
}

func demoScreenBuffer(screen *goterm.Screen) {
	screen.DrawText(4, 4, "Screen Buffer Management:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	w, h := screen.Size()

	y := 6
	screen.DrawText(4, y, fmt.Sprintf("Current Screen Size: %dx%d", w, h), goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)

	y += 2
	screen.DrawText(4, y, "Screen buffer stores cells in memory before rendering", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)
	y++
	screen.DrawText(4, y, "Each cell contains: character, foreground, background, style", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)

	y += 2
	screen.DrawText(4, y, "Buffer Operations:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
	y++

	operations := []string{
		"• SetCell(x, y, cell) - Set individual cell",
		"• GetCell(x, y) - Read cell from buffer",
		"• DrawText(x, y, text, ...) - Draw text string",
		"• Clear() - Reset all cells to default",
		"• Resize(w, h) - Change buffer dimensions",
		"• Show() - Render buffer to terminal",
	}

	for _, op := range operations {
		screen.DrawText(6, y, op, goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleNone)
		y++
	}

	// Demonstrate buffer with pattern
	y += 2
	screen.DrawText(4, y, "Buffer Pattern Example:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
	y++

	// Create a pattern in a box
	boxX, boxY := 6, y
	boxW, boxH := 40, 8
	drawBox(screen, boxX, boxY, boxW, boxH, goterm.ColorYellow, true)

	// Fill with pattern
	for dy := 1; dy < boxH-1; dy++ {
		for dx := 1; dx < boxW-1; dx++ {
			char := rune('A' + (dx+dy)%26)
			color := goterm.ColorIndex(uint8(1 + (dx+dy)%7))
			screen.SetCell(boxX+dx, boxY+dy, goterm.NewCell(char, color, goterm.ColorDefault(), goterm.StyleNone))
		}
	}
}

func demoCellManipulation(screen *goterm.Screen) {
	screen.DrawText(4, 4, "Cell Manipulation:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	y := 6
	screen.DrawText(4, y, "Individual cell operations allow fine-grained control", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)

	y += 2
	screen.DrawText(4, y, "Example 1: Creating a cell", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleBold)
	y++
	screen.DrawText(6, y, "cell := goterm.NewCell('X', ColorRed, ColorBlue, StyleBold)", goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleNone)
	y++
	screen.DrawText(6, y, "Result: ", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	screen.SetCell(15, y, goterm.NewCell('X', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold))

	y += 2
	screen.DrawText(4, y, "Example 2: Cell equality check", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleBold)
	y++
	cell1 := goterm.NewCell('A', goterm.ColorRed, goterm.ColorDefault(), goterm.StyleBold)
	cell2 := goterm.NewCell('A', goterm.ColorRed, goterm.ColorDefault(), goterm.StyleBold)
	screen.DrawText(6, y, fmt.Sprintf("cell1.Equal(cell2) = %v", cell1.Equal(cell2)), goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleNone)

	y += 2
	screen.DrawText(4, y, "Example 3: Grid pattern with cells", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleBold)
	y++

	// Create a checkerboard pattern
	gridX, gridY := 6, y
	gridSize := 16
	for row := 0; row < gridSize; row++ {
		for col := 0; col < gridSize*2; col += 2 {
			isEven := (row+col/2)%2 == 0
			var cell goterm.Cell
			if isEven {
				cell = goterm.NewCell('█', goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
			} else {
				cell = goterm.NewCell('█', goterm.ColorBlack, goterm.ColorDefault(), goterm.StyleNone)
			}
			screen.SetCell(gridX+col, gridY+row, cell)
			screen.SetCell(gridX+col+1, gridY+row, cell)
		}
	}

	screen.DrawText(gridX+gridSize*2+4, gridY+gridSize/2-1, "←", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(gridX+gridSize*2+6, gridY+gridSize/2-1, "Checkerboard", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(gridX+gridSize*2+6, gridY+gridSize/2, "pattern using", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(gridX+gridSize*2+6, gridY+gridSize/2+1, "individual cells", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
}

func demoResize(screen *goterm.Screen) {
	w, h := screen.Size()

	screen.DrawText(4, 4, "Dynamic Resize Support:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	y := 6
	screen.DrawText(4, y, fmt.Sprintf("Current terminal size: %dx%d", w, h), goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)

	y += 2
	screen.DrawText(4, y, "The Screen.Resize() method allows changing buffer dimensions", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)
	y++
	screen.DrawText(4, y, "Content is preserved where it fits in the new size", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)

	y += 2
	screen.DrawText(4, y, "Features:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
	y++
	features := []string{
		"• Content preservation during resize",
		"• Automatic clipping for smaller dimensions",
		"• Default cell filling for larger dimensions",
		"• Thread-safe resize operations",
	}
	for _, f := range features {
		screen.DrawText(6, y, f, goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleNone)
		y++
	}

	y += 2
	screen.DrawText(4, y, "Try resizing your terminal window!", goterm.ColorMagenta, goterm.ColorDefault(), goterm.StyleBold|goterm.StyleItalic)

	// Draw a reference grid to show resize
	y += 2
	for i := 0; i < w && i < 80; i += 10 {
		screen.DrawText(i, y, fmt.Sprintf("%d", i), goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
		for j := 0; j < h-y-3 && j < 20; j++ {
			screen.SetCell(i, y+1+j, goterm.NewCell('│', goterm.ColorBlue, goterm.ColorDefault(), goterm.StyleDim))
		}
	}
}

func demoGradients(screen *goterm.Screen) {
	w, _ := screen.Size()

	screen.DrawText(4, 4, "Color Gradients:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	y := 6
	screen.DrawText(4, y, "Horizontal RGB Gradients:", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)

	// Red to Green gradient
	y += 2
	screen.DrawText(4, y, "Red → Green:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	for i := 0; i < 64 && (18+i) < w; i++ {
		r := uint8(255 - (i * 255 / 63))
		g := uint8(i * 255 / 63)
		color := goterm.ColorRGB(r, g, 0)
		screen.SetCell(18+i, y, goterm.NewCell('█', color, goterm.ColorDefault(), goterm.StyleNone))
	}

	// Green to Blue gradient
	y++
	screen.DrawText(4, y, "Green → Blue:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	for i := 0; i < 64 && (18+i) < w; i++ {
		g := uint8(255 - (i * 255 / 63))
		b := uint8(i * 255 / 63)
		color := goterm.ColorRGB(0, g, b)
		screen.SetCell(18+i, y, goterm.NewCell('█', color, goterm.ColorDefault(), goterm.StyleNone))
	}

	// Blue to Red gradient
	y++
	screen.DrawText(4, y, "Blue → Red:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	for i := 0; i < 64 && (18+i) < w; i++ {
		b := uint8(255 - (i * 255 / 63))
		r := uint8(i * 255 / 63)
		color := goterm.ColorRGB(r, 0, b)
		screen.SetCell(18+i, y, goterm.NewCell('█', color, goterm.ColorDefault(), goterm.StyleNone))
	}

	// Rainbow gradient
	y += 2
	screen.DrawText(4, y, "Rainbow:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	for i := 0; i < 64 && (18+i) < w; i++ {
		hue := float64(i) / 64.0
		r, g, b := hsvToRGB(hue, 1.0, 1.0)
		color := goterm.ColorRGB(r, g, b)
		screen.SetCell(18+i, y, goterm.NewCell('█', color, goterm.ColorDefault(), goterm.StyleNone))
	}

	// Grayscale gradient
	y++
	screen.DrawText(4, y, "Grayscale:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
	for i := 0; i < 64 && (18+i) < w; i++ {
		gray := uint8(i * 255 / 63)
		color := goterm.ColorRGB(gray, gray, gray)
		screen.SetCell(18+i, y, goterm.NewCell('█', color, goterm.ColorDefault(), goterm.StyleNone))
	}

	// 2D gradient
	y += 3
	screen.DrawText(4, y, "2D Gradient (X=Red, Y=Blue):", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	y++

	gridW, gridH := 48, 12
	for dy := 0; dy < gridH; dy++ {
		for dx := 0; dx < gridW; dx++ {
			r := uint8((dx * 255) / (gridW - 1))
			b := uint8((dy * 255) / (gridH - 1))
			color := goterm.ColorRGB(r, 0, b)
			if (6+dx) < w {
				screen.SetCell(6+dx, y+dy, goterm.NewCell('█', color, goterm.ColorDefault(), goterm.StyleNone))
			}
		}
	}
}

func hsvToRGB(h, s, v float64) (uint8, uint8, uint8) {
	if s == 0 {
		gray := uint8(v * 255)
		return gray, gray, gray
	}

	h = h * 6.0
	i := int(h)
	f := h - float64(i)
	p := v * (1.0 - s)
	q := v * (1.0 - s*f)
	t := v * (1.0 - s*(1.0-f))

	var r, g, b float64
	switch i % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}

func demoAnimation(screen *goterm.Screen) {
	screen.DrawText(4, 4, "Animation Demo:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

	screen.DrawText(4, 6, "Simulated animation frames (shown for 3 seconds total)", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)

	// Note: This is a single frame showing what animation looks like
	// Real animation would call Show() multiple times

	y := 8

	// Bouncing ball positions
	positions := []int{10, 15, 18, 20, 18, 15, 10}
	screen.DrawText(4, y, "Bouncing Ball:", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold)
	y++

	for i, pos := range positions {
		x := 6 + i*4
		// Draw ball at position
		screen.SetCell(x, y+20-pos, goterm.NewCell('●', goterm.ColorRed, goterm.ColorDefault(), goterm.StyleBold))
		// Draw shadow
		screen.SetCell(x, y+20, goterm.NewCell('○', goterm.ColorBlack, goterm.ColorDefault(), goterm.StyleDim))
		// Draw frame number
		screen.DrawText(x, y+21, fmt.Sprintf("%d", i+1), goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	}

	// Spinner animation
	spinners := []rune{'|', '/', '─', '\\', '|', '/', '─', '\\'}
	screen.DrawText(40, y, "Spinner:", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold)
	for i, spin := range spinners {
		x := 40 + 10 + i*3
		screen.SetCell(x, y+1, goterm.NewCell(spin, goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleBold))
		screen.DrawText(x-1, y+2, fmt.Sprintf("%d", i+1), goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	}

	// Color wave
	screen.DrawText(4, y+24, "Color Wave:", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold)
	for i := 0; i < 60; i++ {
		hue := float64(i) / 60.0
		r, g, b := hsvToRGB(hue, 1.0, 1.0)
		color := goterm.ColorRGB(r, g, b)
		screen.SetCell(4+i, y+25, goterm.NewCell('█', color, goterm.ColorDefault(), goterm.StyleNone))
	}

	screen.DrawText(4, y+27, "Note: Real animations call Show() in a loop for smooth motion", goterm.ColorMagenta, goterm.ColorDefault(), goterm.StyleItalic)
}

func demoFinalShowcase(screen *goterm.Screen) {
	w, h := screen.Size()

	// Create a stunning final showcase

	// Title
	title := "goterm - Terminal Graphics Library"
	screen.DrawText((w-len(title))/2, 4, title, goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleBold)

	// Feature highlights in boxes
	features := []struct {
		title string
		items []string
		color goterm.Color
	}{
		{
			title: "Colors",
			items: []string{
				"16 ANSI colors",
				"256 color palette",
				"24-bit RGB (16M)",
				"Auto conversion",
			},
			color: goterm.ColorRed,
		},
		{
			title: "Styles",
			items: []string{
				"Bold, Italic",
				"Underline",
				"Reverse, Dim",
				"Combinable",
			},
			color: goterm.ColorGreen,
		},
		{
			title: "Features",
			items: []string{
				"Screen buffer",
				"Cell-based",
				"Thread-safe",
				"Unicode/CJK",
			},
			color: goterm.ColorBlue,
		},
	}

	startX := (w - 75) / 2
	if startX < 2 {
		startX = 2
	}
	y := 7

	for i, feat := range features {
		x := startX + i*25

		// Draw box
		drawBox(screen, x, y, 22, 8, feat.color, true)

		// Title
		titleX := x + (22-len(feat.title))/2
		screen.DrawText(titleX, y+1, feat.title, feat.color, goterm.ColorDefault(), goterm.StyleBold)

		// Items
		for j, item := range feat.items {
			screen.DrawText(x+2, y+3+j, item, goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
		}
	}

	// Rainbow divider
	y = 17
	for i := 0; i < w-8; i++ {
		hue := float64(i) / float64(w-8)
		r, g, b := hsvToRGB(hue, 1.0, 1.0)
		color := goterm.ColorRGB(r, g, b)
		screen.SetCell(4+i, y, goterm.NewCell('═', color, goterm.ColorDefault(), goterm.StyleBold))
	}

	// Key features list
	y = 19
	keyFeatures := []string{
		"✓ Multiple color modes with automatic conversion",
		"✓ Rich text styling with bitwise operations",
		"✓ Full Unicode and CJK character support",
		"✓ Thread-safe screen buffer management",
		"✓ Efficient cell-based rendering",
		"✓ Cross-platform terminal support",
	}

	centerX := (w - 50) / 2
	if centerX < 4 {
		centerX = 4
	}

	for i, feat := range keyFeatures {
		if y+i < h-4 {
			screen.DrawText(centerX, y+i, feat, goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
		}
	}

	// Final message
	finalY := h - 6
	if finalY > y+len(keyFeatures)+2 {
		msg := "Ready to build amazing terminal UIs!"
		msgX := (w - len(msg)) / 2
		screen.DrawText(msgX, finalY, msg, goterm.ColorMagenta, goterm.ColorDefault(), goterm.StyleBold|goterm.StyleItalic)

		repo := "github.com/dshills/goterm"
		repoX := (w - len(repo)) / 2
		screen.DrawText(repoX, finalY+2, repo, goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleUnderline)
	}
}
