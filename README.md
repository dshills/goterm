# goterm

[![Go Version](https://img.shields.io/badge/Go-1.25.3+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

 > **⚠️ ALPHA SOFTWARE WARNING**
 >
 > This project is currently in alpha stage. While we strive for quality and stability, you should expect:
 > - Potential bugs and rough edges
 > - Breaking API changes without prior notice
 > - Incomplete documentation in some areas
 > - Features that may be refined or removed
 >
 > **Not recommended for production use** without thorough testing and understanding of the codebase. We welcome early adopters,
 + testers, and contributors! Please report issues on our [issue tracker](https://github.com/dshills/gokeys/issues).

A powerful, feature-rich Go library for building terminal user interfaces with comprehensive color support, text styling, and screen buffer management.


## Features

### 🎨 Multi-Mode Color Support
- **16 ANSI Colors** - Standard terminal colors with bright variants
- **256 Color Palette** - Extended palette with 6x6x6 color cube and 24-shade grayscale
- **24-bit True Color** - RGB color support (16.7 million colors)
- **Automatic Conversion** - Seamless degradation between color modes

### ✨ Rich Text Styling
- **9 Style Types** - Bold, Dim, Italic, Underline, Blink, Reverse, Conceal, Strikethrough
- **Combinable Styles** - Use bitwise OR to combine multiple styles
- **Style Operations** - Set, Clear, Toggle, and Has methods for fine control

### 📝 Unicode & Character Support
- **Full Unicode** - Complete Unicode support including emoji 🎮 🚀 ⚡
- **CJK Characters** - Chinese, Japanese, Korean character support
- **Box Drawing** - Single and double-line box drawing characters
- **Multi-byte Safe** - Proper handling of multi-byte UTF-8 sequences

### 🖥️ Screen Management
- **Buffer-based Rendering** - Off-screen buffer for flicker-free updates
- **Cell-based Model** - Individual control over each terminal cell
- **Thread-safe** - Concurrent access protection with RWMutex
- **Dynamic Resize** - Handle terminal size changes gracefully

## Installation

```bash
go get github.com/dshills/goterm
```

## Quick Start

```go
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
        fmt.Fprintf(os.Stderr, "Failed to initialize: %v\n", err)
        os.Exit(1)
    }
    defer func() {
        if err := screen.Close(); err != nil {
            fmt.Fprintf(os.Stderr, "Failed to close: %v\n", err)
        }
    }()

    // Clear the screen
    screen.Clear()

    // Draw colorful text
    screen.DrawText(2, 2, "Hello, Terminal!",
        goterm.ColorGreen,
        goterm.ColorDefault(),
        goterm.StyleBold)

    // Draw with RGB colors
    orange := goterm.ColorRGB(255, 128, 0)
    screen.DrawText(2, 4, "True Color Support!",
        orange,
        goterm.ColorDefault(),
        goterm.StyleNone)

    // Draw a box
    screen.DrawText(2, 6, "┌────────────┐",
        goterm.ColorCyan,
        goterm.ColorDefault(),
        goterm.StyleNone)
    screen.DrawText(2, 7, "│  Box Demo  │",
        goterm.ColorCyan,
        goterm.ColorDefault(),
        goterm.StyleNone)
    screen.DrawText(2, 8, "└────────────┘",
        goterm.ColorCyan,
        goterm.ColorDefault(),
        goterm.StyleNone)

    // Show the screen
    if err := screen.Show(); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to show: %v\n", err)
        return
    }

    time.Sleep(3 * time.Second)
}
```

## Examples

### Running the Demo

See all features in action with the comprehensive demo:

```bash
go run ./examples/demo
```

The demo showcases 15 different feature screens including colors, styles, Unicode, animations, and more!

### Simple Example

```bash
go run ./examples/hello
```

## Color Usage

### ANSI 16 Colors

```go
// Named colors
screen.DrawText(0, 0, "Red text", goterm.ColorRed, goterm.ColorDefault(), goterm.StyleNone)
screen.DrawText(0, 1, "Blue text", goterm.ColorBlue, goterm.ColorDefault(), goterm.StyleNone)

// Indexed colors (0-15)
brightRed := goterm.ColorIndex(9)
screen.DrawText(0, 2, "Bright red", brightRed, goterm.ColorDefault(), goterm.StyleNone)
```

### 256 Color Palette

```go
// Use any color from the 256-color palette
color := goterm.ColorIndex(208) // Orange from palette
screen.DrawText(0, 0, "Orange text", color, goterm.ColorDefault(), goterm.StyleNone)

// Grayscale (232-255)
gray := goterm.ColorIndex(240)
screen.DrawText(0, 1, "Gray text", gray, goterm.ColorDefault(), goterm.StyleNone)
```

### True Color (RGB)

```go
// Create custom RGB colors
purple := goterm.ColorRGB(128, 0, 128)
teal := goterm.ColorRGB(0, 128, 128)
pink := goterm.ColorRGB(255, 192, 203)

screen.DrawText(0, 0, "Purple", purple, goterm.ColorDefault(), goterm.StyleNone)
screen.DrawText(0, 1, "Teal", teal, goterm.ColorDefault(), goterm.StyleNone)
screen.DrawText(0, 2, "Pink", pink, goterm.ColorDefault(), goterm.StyleNone)
```

### Color Conversions

```go
// Automatic color mode conversion
rgb := goterm.ColorRGB(128, 200, 255)

// Convert to 256-color mode
color256 := rgb.To256()

// Convert to 16-color mode
color16 := rgb.To16()

// Use the converted colors
screen.DrawText(0, 0, "RGB", rgb, goterm.ColorDefault(), goterm.StyleNone)
screen.DrawText(0, 1, "256", color256, goterm.ColorDefault(), goterm.StyleNone)
screen.DrawText(0, 2, "16", color16, goterm.ColorDefault(), goterm.StyleNone)
```

### Background Colors

```go
// Use any color as background
screen.DrawText(0, 0, "White on Blue",
    goterm.ColorWhite,
    goterm.ColorBlue,
    goterm.StyleBold)

// RGB background
bg := goterm.ColorRGB(40, 40, 40)
fg := goterm.ColorRGB(255, 255, 0)
screen.DrawText(0, 1, "Yellow on dark gray", fg, bg, goterm.StyleNone)
```

## Text Styling

### Basic Styles

```go
// Individual styles
screen.DrawText(0, 0, "Bold", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)
screen.DrawText(0, 1, "Italic", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleItalic)
screen.DrawText(0, 2, "Underline", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleUnderline)
screen.DrawText(0, 3, "Reverse", goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleReverse)
```

### Combining Styles

```go
// Use bitwise OR to combine styles
combined := goterm.StyleBold | goterm.StyleUnderline
screen.DrawText(0, 0, "Bold and Underlined",
    goterm.ColorWhite,
    goterm.ColorDefault(),
    combined)

// Triple combination
triple := goterm.StyleBold | goterm.StyleItalic | goterm.StyleUnderline
screen.DrawText(0, 1, "Bold, Italic, Underlined",
    goterm.ColorWhite,
    goterm.ColorDefault(),
    triple)
```

### Style Operations

```go
// Start with bold
style := goterm.StyleBold

// Add underline
style = style.Set(goterm.StyleUnderline)

// Toggle italic
style = style.Toggle(goterm.StyleItalic)

// Check if has bold
if style.Has(goterm.StyleBold) {
    // Bold is active
}

// Remove underline
style = style.Clear(goterm.StyleUnderline)
```

## Screen Buffer Operations

### Individual Cell Control

```go
// Create a cell
cell := goterm.NewCell('X', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold)

// Set cell at position
screen.SetCell(10, 5, cell)

// Get cell from position
cell = screen.GetCell(10, 5)

// Check cell properties
if cell.Ch == 'X' {
    // Character matches
}

// Compare cells
cell2 := goterm.NewCell('X', goterm.ColorRed, goterm.ColorBlue, goterm.StyleBold)
if cell.Equal(cell2) {
    // Cells are identical
}
```

### Drawing Text

```go
// Simple text drawing
screen.DrawText(5, 10, "Hello, World!",
    goterm.ColorGreen,
    goterm.ColorDefault(),
    goterm.StyleNone)

// Text automatically wraps at screen edge
screen.DrawText(0, 0, "This is a very long line that will be clipped at the screen edge",
    goterm.ColorWhite,
    goterm.ColorDefault(),
    goterm.StyleNone)
```

### Screen Operations

```go
// Clear entire screen
screen.Clear()

// Get screen dimensions
width, height := screen.Size()

// Resize screen buffer
screen.Resize(100, 50)

// Render buffer to terminal
if err := screen.Show(); err != nil {
    // Handle error
}

// Flush output
if err := screen.Sync(); err != nil {
    // Handle error
}
```

## Building Complex UIs

### Drawing Boxes

```go
// Single-line box
screen.DrawText(0, 0, "┌──────┐", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)
screen.DrawText(0, 1, "│ Box  │", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)
screen.DrawText(0, 2, "└──────┘", goterm.ColorCyan, goterm.ColorDefault(), goterm.StyleNone)

// Double-line box
screen.DrawText(10, 0, "╔══════╗", goterm.ColorMagenta, goterm.ColorDefault(), goterm.StyleNone)
screen.DrawText(10, 1, "║ Box  ║", goterm.ColorMagenta, goterm.ColorDefault(), goterm.StyleNone)
screen.DrawText(10, 2, "╚══════╝", goterm.ColorMagenta, goterm.ColorDefault(), goterm.StyleNone)
```

### Creating Gradients

```go
// Horizontal gradient
for i := 0; i < 64; i++ {
    r := uint8((i * 255) / 63)
    color := goterm.ColorRGB(r, 0, 255-r)
    screen.SetCell(i, 10, goterm.NewCell('█', color, goterm.ColorDefault(), goterm.StyleNone))
}
```

### Animation

```go
// Animation loop
for frame := 0; frame < 60; frame++ {
    screen.Clear()

    // Draw animated content
    x := frame % width
    screen.DrawText(x, 10, "->", goterm.ColorGreen, goterm.ColorDefault(), goterm.StyleBold)

    // Render frame
    if err := screen.Show(); err != nil {
        break
    }

    time.Sleep(50 * time.Millisecond)
}
```

## API Reference

### Initialization

```go
// Initialize terminal and create screen
screen, err := goterm.Init()

// Close and restore terminal (always defer this)
defer screen.Close()
```

### Colors

```go
// Create colors
goterm.ColorDefault()              // Terminal default color
goterm.ColorRGB(r, g, b uint8)     // 24-bit RGB color
goterm.ColorIndex(index uint8)     // 256-color palette (0-255)

// Named colors (indices 0-7)
goterm.ColorBlack, goterm.ColorRed, goterm.ColorGreen, goterm.ColorYellow
goterm.ColorBlue, goterm.ColorMagenta, goterm.ColorCyan, goterm.ColorWhite

// Color methods
color.Mode() ColorMode             // Get color mode
color.RGB() (r, g, b uint8)       // Get RGB values
color.Index() uint8               // Get palette index
color.To256() Color               // Convert to 256-color
color.To16() Color                // Convert to 16-color
```

### Styles

```go
// Style constants
goterm.StyleNone, goterm.StyleBold, goterm.StyleDim, goterm.StyleItalic
goterm.StyleUnderline, goterm.StyleSlowBlink, goterm.StyleRapidBlink
goterm.StyleReverse, goterm.StyleConceal, goterm.StyleStrikethrough

// Style methods
style.Has(flag Style) bool        // Check if style flag is set
style.Set(flag Style) Style       // Add style flag
style.Clear(flag Style) Style     // Remove style flag
style.Toggle(flag Style) Style    // Toggle style flag
```

### Cells

```go
// Create cell
goterm.NewCell(ch rune, fg, bg Color, style Style) Cell

// Cell methods
cell.Clear()                      // Reset to default
cell.Equal(other Cell) bool       // Compare cells
```

### Screen

```go
// Screen methods
screen.Size() (width, height int)
screen.SetCell(x, y int, cell Cell)
screen.GetCell(x, y int) Cell
screen.Clear()
screen.DrawText(x, y int, text string, fg, bg Color, style Style)
screen.Resize(width, height int)
screen.Show() error
screen.Sync() error
screen.Close() error
```

## Error Handling

```go
// Common errors
goterm.ErrNotATerminal             // stdout is not a terminal
goterm.ErrTerminalSetupFailed      // Terminal initialization failed
goterm.ErrTerminalRestoreFailed    // Terminal restoration failed

// Error handling example
screen, err := goterm.Init()
if err != nil {
    if errors.Is(err, goterm.ErrNotATerminal) {
        fmt.Fprintln(os.Stderr, "Not running in a terminal")
    } else {
        fmt.Fprintf(os.Stderr, "Init failed: %v\n", err)
    }
    os.Exit(1)
}
```

## Thread Safety

The `Screen` type uses `sync.RWMutex` for thread-safe operations:

```go
// Safe to call from multiple goroutines
go func() {
    screen.SetCell(10, 10, cell1)
}()

go func() {
    screen.SetCell(20, 20, cell2)
}()

// Show() can be called after all updates complete
screen.Show()
```

## Best Practices

### Always Close the Screen

```go
screen, err := goterm.Init()
if err != nil {
    // handle error
}
defer func() {
    if err := screen.Close(); err != nil {
        fmt.Fprintf(os.Stderr, "Close failed: %v\n", err)
    }
}()
```

### Check Terminal Capabilities

```go
// Some terminals don't support all features
// Test on your target terminals

// True color support varies
rgb := goterm.ColorRGB(128, 200, 255)
// Fallback for terminals without true color
color256 := rgb.To256()
```

### Buffer Before Showing

```go
// Bad: Multiple Show() calls cause flickering
screen.DrawText(0, 0, "Line 1", ...)
screen.Show()
screen.DrawText(0, 1, "Line 2", ...)
screen.Show()

// Good: Draw everything, then show once
screen.DrawText(0, 0, "Line 1", ...)
screen.DrawText(0, 1, "Line 2", ...)
screen.DrawText(0, 2, "Line 3", ...)
screen.Show()
```

### Handle Resize Events

```go
// Listen for terminal resize (using SIGWINCH)
// Then call screen.Resize() with new dimensions

// Example with signal handling
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, syscall.SIGWINCH)

go func() {
    for range sigCh {
        // Get new terminal size
        width, height, _ := term.GetSize(int(os.Stdout.Fd()))
        screen.Resize(width, height)
        // Redraw content
        redraw(screen)
        screen.Show()
    }
}()
```

## Performance Considerations

### Minimize Show() Calls

```go
// Each Show() renders the entire buffer
// For animations, target 30-60 FPS max

ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
for range ticker.C {
    updateGameState()
    renderFrame(screen)
    screen.Show()
}
```

### Use Cell Manipulation for Efficiency

```go
// For small updates, use SetCell instead of DrawText
// Changes only what's needed

// Update single character
screen.SetCell(x, y, goterm.NewCell('X', color, bg, style))
```

### Reuse Colors and Styles

```go
// Create color/style constants instead of recreating
var (
    titleColor = goterm.ColorRGB(100, 200, 255)
    titleStyle = goterm.StyleBold | goterm.StyleUnderline
    normalColor = goterm.ColorWhite
)

// Reuse throughout your application
screen.DrawText(0, 0, "Title", titleColor, goterm.ColorDefault(), titleStyle)
screen.DrawText(0, 2, "Content", normalColor, goterm.ColorDefault(), goterm.StyleNone)
```

## Requirements

- **Go**: 1.25.3 or later
- **OS**: Unix-like systems (Linux, macOS, BSD)
- **Terminal**: Any terminal with color support

## Project Structure

```
goterm/
├── cell.go              # Cell structure and operations
├── color.go             # Color types and conversions
├── errors.go            # Error definitions
├── event.go             # Event system (mouse, keyboard)
├── screen.go            # Screen buffer management
├── style.go             # Text styling
├── terminal.go          # Terminal interface
├── examples/
│   ├── hello/          # Simple example
│   └── demo/           # Comprehensive demo
└── tests/
    ├── unit/           # Unit tests
    └── integration/    # Integration tests
```

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test suite
go test ./tests/unit
go test ./tests/integration

# Run benchmarks
go test -bench=. ./tests/benchmark/...
```

## Development

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Run tests with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Using the Makefile

```bash
make help        # Show available commands
make test        # Run all tests
make lint        # Run linters
make fmt         # Format code
make coverage    # Generate coverage report
make clean       # Clean build artifacts
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [golang.org/x/term](https://pkg.go.dev/golang.org/x/term) for terminal control
- Inspired by terminal UI libraries in other languages
- Special thanks to the Go community

## Related Projects

- [tcell](https://github.com/gdamore/tcell) - Full-featured terminal UI library
- [termbox-go](https://github.com/nsf/termbox-go) - Minimalist terminal library
- [tview](https://github.com/rivo/tview) - Terminal UI widgets

## Support

- **Issues**: [GitHub Issues](https://github.com/dshills/goterm/issues)
- **Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/dshills/goterm)
- **Examples**: See `examples/` directory

## Roadmap

- [ ] Keyboard event handling
- [ ] Mouse event handling
- [ ] Platform-specific optimizations
- [ ] Additional box drawing utilities
- [ ] Color palette customization
- [ ] Alternative screen buffer support

---

**Built with ❤️ by the goterm team**

*Making terminal UIs beautiful, one cell at a time.*
