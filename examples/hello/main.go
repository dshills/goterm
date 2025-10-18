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

	// Clear the screen
	screen.Clear()

	// Draw a simple message with colors and styles
	screen.DrawText(2, 2, "Hello, Terminal Graphics!", goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleBold)
	screen.DrawText(2, 4, "This is a simple demo using goterm", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)

	// Draw a colored box
	screen.DrawText(2, 6, "┌────────────────────────┐", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(2, 7, "│  Basic Rendering Demo  │", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(2, 8, "└────────────────────────┘", goterm.ColorYellow, goterm.ColorDefault(), goterm.StyleNone)

	// Show various colors
	screen.DrawText(2, 10, "Colors:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
	colors := []struct {
		color goterm.Color
		name  string
	}{
		{goterm.ColorRed, "Red"},
		{goterm.ColorGreen, "Green"},
		{goterm.ColorYellow, "Yellow"},
		{goterm.ColorBlue, "Blue"},
		{goterm.ColorMagenta, "Magenta"},
		{goterm.ColorCyan, "Cyan"},
		{goterm.ColorWhite, "White"},
	}

	for i, c := range colors {
		screen.DrawText(4, 11+i, c.name, c.color, goterm.ColorDefault(), goterm.StyleNone)
	}

	// Show various styles
	screen.DrawText(15, 10, "Styles:", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
	screen.DrawText(17, 11, "Bold", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
	screen.DrawText(17, 12, "Italic", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleItalic)
	screen.DrawText(17, 13, "Underline", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleUnderline)
	screen.DrawText(17, 14, "Reverse", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleReverse)
	screen.DrawText(17, 15, "Bold+Underline", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold|goterm.StyleUnderline)

	// Unicode support
	screen.DrawText(2, 19, "Unicode: 🎮 ⚡ 🚀 ★ ♥ ♦ ☺", goterm.ColorMagenta, goterm.ColorDefault(), goterm.StyleNone)
	screen.DrawText(2, 20, "CJK: 日本語 中文 한국어", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)

	// Show the screen
	if err := screen.Show(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to show screen: %v\n", err)
		return
	}

	// Keep the screen visible for 5 seconds
	time.Sleep(5 * time.Second)

	// Clear and show exit message
	screen.Clear()
	screen.DrawText(2, 2, "Goodbye!", goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleBold)
	if err := screen.Show(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to show exit message: %v\n", err)
	}

	time.Sleep(1 * time.Second)
}
