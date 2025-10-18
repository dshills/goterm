# Public API Contract: Terminal Graphics Library

**Feature**: Terminal Graphics Library
**Branch**: `001-specify-scripts-bash`
**Date**: 2025-10-18
**Version**: 0.1.0 (initial)

## Overview

This document defines the public API contract for the terminal graphics library. This is a library API (not HTTP/REST), so "contracts" refer to Go package interfaces and function signatures that external applications will use.

## Package Structure

```
github.com/dshills/goterm
├── screen.go       - Main screen interface
├── cell.go         - Cell and cell manipulation
├── color.go        - Color types and constants
├── style.go        - Style types and constants
├── event.go        - Event types
└── errors.go       - Error types
```

## Core API

### Initialization and Lifecycle

#### `Init() (*Screen, error)`

Initialize the terminal graphics library and return a Screen instance.

**Preconditions**:
- Must be called from a process with a controlling terminal
- stdin/stdout must be connected to a terminal (not piped)

**Postconditions**:
- Terminal is in raw mode
- Alternate screen buffer is active
- Mouse tracking is enabled
- Cursor is hidden
- Screen buffer is allocated and cleared

**Errors**:
- `ErrNotATerminal`: stdout is not a terminal
- `ErrTerminalSetupFailed`: failed to configure terminal

**Example**:
```go
screen, err := goterm.Init()
if err != nil {
    log.Fatal(err)
}
defer screen.Close()
```

#### `(*Screen) Close() error`

Clean up and restore terminal to original state.

**Preconditions**:
- Screen was previously initialized

**Postconditions**:
- Terminal restored to original state
- Alternate screen buffer exited
- Mouse tracking disabled
- Cursor shown
- Screen buffer deallocated

**Errors**:
- `ErrTerminalRestoreFailed`: failed to restore terminal state

**Example**:
```go
if err := screen.Close(); err != nil {
    log.Printf("failed to close screen: %v", err)
}
```

### Screen Information

#### `(*Screen) Size() (width, height int)`

Get current terminal dimensions.

**Returns**:
- width: Number of columns (≥ 1)
- height: Number of rows (≥ 1)

**Example**:
```go
w, h := screen.Size()
fmt.Printf("Terminal is %d×%d\n", w, h)
```

#### `(*Screen) Clear()`

Clear entire screen buffer to default colors.

**Postconditions**:
- All cells set to space character
- All cells use default foreground/background colors
- All cells have no style
- All cells marked dirty

**Example**:
```go
screen.Clear()
```

### Cell Manipulation

#### `(*Screen) SetCell(x, y int, ch rune, fg, bg Color, style Style)`

Set a single cell in the screen buffer.

**Parameters**:
- x: Column position (0-based, 0 ≤ x < width)
- y: Row position (0-based, 0 ≤ y < height)
- ch: Character to display
- fg: Foreground color
- bg: Background color
- style: Text styling flags

**Preconditions**:
- Coordinates must be within bounds
- Character must be valid Unicode

**Postconditions**:
- Cell at (x, y) updated with new values
- Cell marked dirty

**Panics**:
- If coordinates are out of bounds (programming error)

**Example**:
```go
screen.SetCell(10, 5, '@',
    goterm.ColorRGB(255, 255, 0),  // Yellow
    goterm.ColorDefault(),
    goterm.StyleBold)
```

#### `(*Screen) GetCell(x, y int) (ch rune, fg, bg Color, style Style)`

Get current cell content.

**Parameters**:
- x, y: Cell coordinates

**Returns**:
- ch: Current character
- fg: Foreground color
- bg: Background color
- style: Current style

**Panics**:
- If coordinates are out of bounds

**Example**:
```go
ch, fg, bg, style := screen.GetCell(10, 5)
```

#### `(*Screen) SetContent(x, y int, ch rune, fg, bg Color, style Style, comb ...rune)`

Set cell with combining characters (for complex Unicode).

**Parameters**:
- x, y: Cell coordinates
- ch: Primary character
- fg, bg: Colors
- style: Text styling
- comb: Optional combining characters

**Example**:
```go
// Draw accented character
screen.SetContent(x, y, 'e', fg, bg, style, '\u0301')  // é
```

### Text Rendering

#### `(*Screen) DrawText(x, y int, text string, fg, bg Color, style Style)`

Draw a string of text starting at position (x, y).

**Parameters**:
- x, y: Starting coordinates
- text: String to draw
- fg, bg: Colors for all characters
- style: Styling for all characters

**Behavior**:
- Text wraps to next line if it exceeds screen width
- Out-of-bounds characters are clipped
- Handles Unicode properly (multi-byte characters count as one cell)

**Example**:
```go
screen.DrawText(0, 0, "Hello, World!",
    goterm.ColorWhite(),
    goterm.ColorBlue(),
    goterm.StyleBold)
```

### Rendering

#### `(*Screen) Show()`

Render all dirty cells to the terminal.

**Preconditions**:
- Screen is initialized

**Postconditions**:
- All dirty cells rendered to terminal using escape sequences
- Dirty flags cleared
- Previous buffer updated with current state
- Terminal output flushed

**Performance**:
- Only changed cells are rendered (delta rendering)
- Batched escape sequences for efficiency
- Target: <16ms for full-screen update on 80×24 terminal

**Example**:
```go
screen.Clear()
screen.DrawText(0, 0, "Frame 1", fg, bg, style)
screen.Show()  // Render to terminal
```

#### `(*Screen) Sync()`

Force full screen redraw (ignore dirty tracking).

**Use Cases**:
- Terminal corruption recovery
- After resume from suspend
- When unsure of terminal state

**Example**:
```go
screen.Sync()  // Full redraw
```

### Cursor Control

#### `(*Screen) ShowCursor()`

Make cursor visible.

**Example**:
```go
screen.ShowCursor()
```

#### `(*Screen) HideCursor()`

Make cursor invisible.

**Example**:
```go
screen.HideCursor()
```

#### `(*Screen) SetCursor(x, y int)`

Set cursor position (if visible).

**Parameters**:
- x, y: Cursor coordinates

**Example**:
```go
screen.SetCursor(10, 5)
screen.ShowCursor()
```

### Event Handling

#### `(*Screen) PollEvent() Event`

Poll for next event (non-blocking).

**Returns**:
- Event: Mouse or resize event, or nil if no events available

**Event Types**:
- `*MouseEvent`: Mouse interaction
- `*ResizeEvent`: Terminal resized
- `nil`: No events pending

**Example**:
```go
if event := screen.PollEvent(); event != nil {
    switch e := event.(type) {
    case *goterm.MouseEvent:
        handleMouse(e)
    case *goterm.ResizeEvent:
        handleResize(e)
    }
}
```

#### `(*Screen) Events() <-chan Event`

Get read-only event channel.

**Returns**:
- Channel that delivers events as they occur

**Example**:
```go
events := screen.Events()
for event := range events {
    // Handle event
}
```

## Color API

### Color Creation

#### `ColorDefault() Color`

Terminal's default color.

**Example**:
```go
fg := goterm.ColorDefault()
```

#### `ColorRGB(r, g, b uint8) Color`

Create true color (24-bit RGB).

**Parameters**:
- r, g, b: Red, green, blue components (0-255)

**Behavior**:
- Automatically degrades to 256-color or 16-color on terminals without true color support

**Example**:
```go
yellow := goterm.ColorRGB(255, 255, 0)
```

#### `ColorIndex(index uint8) Color`

Create indexed color (0-255).

**Parameters**:
- index: Palette index

**Ranges**:
- 0-7: Basic ANSI colors
- 8-15: Bright ANSI colors
- 16-231: 216-color cube
- 232-255: Grayscale ramp

**Example**:
```go
brightRed := goterm.ColorIndex(9)
```

### Named Colors (Constants)

```go
const (
    ColorBlack   = Color{...}
    ColorRed     = Color{...}
    ColorGreen   = Color{...}
    ColorYellow  = Color{...}
    ColorBlue    = Color{...}
    ColorMagenta = Color{...}
    ColorCyan    = Color{...}
    ColorWhite   = Color{...}
)
```

**Example**:
```go
screen.SetCell(x, y, 'X', goterm.ColorRed, goterm.ColorBlack, goterm.StyleNone)
```

## Style API

### Style Flags (Constants)

```go
const (
    StyleNone         Style = 0
    StyleBold         Style = 1 << 0
    StyleDim          Style = 1 << 1
    StyleItalic       Style = 1 << 2
    StyleUnderline    Style = 1 << 3
    StyleSlowBlink    Style = 1 << 4
    StyleRapidBlink   Style = 1 << 5
    StyleReverse      Style = 1 << 6
    StyleConceal      Style = 1 << 7
    StyleStrikethrough Style = 1 << 8
)
```

### Style Combination

Styles can be combined using bitwise OR:

```go
style := goterm.StyleBold | goterm.StyleUnderline
screen.SetCell(x, y, 'A', fg, bg, style)
```

## Event Types

### MouseEvent

```go
type MouseEvent struct {
    X, Y      int          // Cell coordinates
    Button    MouseButton  // Which button
    Modifiers Modifier     // Keyboard modifiers
    Action    MouseAction  // Press/Release/Motion/Scroll
}
```

**Mouse Buttons**:
```go
const (
    MouseLeft      MouseButton = 0
    MouseMiddle    MouseButton = 1
    MouseRight     MouseButton = 2
    MouseNone      MouseButton = 3  // Motion without button
    MouseWheelUp   MouseButton = 4
    MouseWheelDown MouseButton = 5
)
```

**Mouse Actions**:
```go
const (
    MousePress   MouseAction = 0
    MouseRelease MouseAction = 1
    MouseMotion  MouseAction = 2  // Drag
    MouseScroll  MouseAction = 3
)
```

**Modifiers**:
```go
const (
    ModShift Modifier = 1 << 0
    ModAlt   Modifier = 1 << 1
    ModCtrl  Modifier = 1 << 2
)
```

**Example**:
```go
func handleMouse(e *goterm.MouseEvent) {
    if e.Button == goterm.MouseLeft && e.Action == goterm.MousePress {
        fmt.Printf("Left click at (%d, %d)\n", e.X, e.Y)
    }
    if e.Modifiers & goterm.ModCtrl != 0 {
        fmt.Println("Ctrl was held")
    }
}
```

### ResizeEvent

```go
type ResizeEvent struct {
    Width  int
    Height int
}
```

**Example**:
```go
func handleResize(e *goterm.ResizeEvent) {
    fmt.Printf("Terminal resized to %d×%d\n", e.Width, e.Height)
    // Recreate game world for new size
}
```

## Error Types

```go
var (
    ErrNotATerminal        = errors.New("not a terminal")
    ErrTerminalSetupFailed = errors.New("terminal setup failed")
    ErrTerminalRestoreFailed = errors.New("terminal restore failed")
)
```

## Options API

### Initialization Options

#### `WithExternalKeyboard() Option`

Disable internal keyboard handling (for gokeys integration).

**Example**:
```go
screen, err := goterm.Init(goterm.WithExternalKeyboard())
```

#### `WithMouseTracking(enable bool) Option`

Enable or disable mouse tracking.

**Example**:
```go
screen, err := goterm.Init(goterm.WithMouseTracking(false))
```

#### `WithAlternateScreen(enable bool) Option`

Enable or disable alternate screen buffer.

**Example**:
```go
screen, err := goterm.Init(goterm.WithAlternateScreen(true))
```

## Usage Patterns

### Basic Game Loop

```go
package main

import (
    "time"
    "github.com/dshills/goterm"
    "github.com/dshills/gokeys/input"
)

func main() {
    // Initialize graphics
    screen, err := goterm.Init()
    if err != nil {
        panic(err)
    }
    defer screen.Close()

    // Initialize input (gokeys)
    inp := input.New()
    if err := inp.Start(); err != nil {
        panic(err)
    }
    defer inp.Stop()

    // Game state
    player := NewPlayer()

    // Game loop (60 fps)
    ticker := time.NewTicker(16 * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {
        // 1. Handle input
        if event := inp.Next(); event != nil {
            if event.Key == input.KeyEscape {
                break
            }
            handleInput(player, event)
        }

        // 2. Handle events (mouse, resize)
        if event := screen.PollEvent(); event != nil {
            switch e := event.(type) {
            case *goterm.MouseEvent:
                handleMouse(player, e)
            case *goterm.ResizeEvent:
                resizeGame(e.Width, e.Height)
            }
        }

        // 3. Update game logic
        player.Update()

        // 4. Render
        screen.Clear()
        drawPlayer(screen, player)
        drawWorld(screen)
        screen.Show()
    }
}
```

### Drawing Helpers

```go
// Draw a box
func DrawBox(s *goterm.Screen, x, y, w, h int, fg, bg goterm.Color) {
    for row := y; row < y+h; row++ {
        for col := x; col < x+w; col++ {
            var ch rune
            if row == y || row == y+h-1 {
                if col == x || col == x+w-1 {
                    ch = '+'
                } else {
                    ch = '-'
                }
            } else {
                if col == x || col == x+w-1 {
                    ch = '|'
                } else {
                    ch = ' '
                }
            }
            s.SetCell(col, row, ch, fg, bg, goterm.StyleNone)
        }
    }
}
```

## API Stability Guarantees

### Semantic Versioning

The library follows semantic versioning (MAJOR.MINOR.PATCH):
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Deprecation Policy

Before removing functionality:
1. Mark as deprecated with `// Deprecated:` comment
2. Provide alternative in deprecation notice
3. Keep for at least one minor version
4. Remove in next major version

### API Compatibility

**Version 0.x.x** (current):
- API may change between minor versions
- Breaking changes documented in CHANGELOG

**Version 1.0.0+** (future):
- API stability guaranteed within major version
- Breaking changes only in major version bumps

## Performance Contracts

### Rendering Performance

- `Show()` must complete in <16ms for 80×24 terminal at 60 fps
- `SetCell()` must be O(1)
- `Clear()` must be O(n) where n = width × height
- Delta rendering must reduce I/O by >80% for typical game frames

### Memory Usage

- Screen buffer: O(width × height)
- No allocations in render loop (after initialization)
- GC pressure: Zero allocations per frame

## Thread Safety

**Not thread-safe by default**:
- Screen buffer operations are NOT thread-safe
- Applications must synchronize access if using multiple goroutines

**Recommended pattern**:
- Single goroutine owns screen
- Other goroutines communicate via channels

## Platform Support

### Supported Platforms

- Linux (x86_64, arm64)
- macOS (x86_64, arm64)
- Windows (x86_64) - via VT100 emulation

### Supported Terminals

Minimum: VT100-compatible terminal
Recommended: xterm-256color or better

**Tested terminals**:
- xterm
- gnome-terminal
- iTerm2
- Windows Terminal
- kitty
- alacritty

## gokeys Integration

### Coordination

The library is designed to work seamlessly with gokeys:

1. **Terminal management**: goterm owns terminal state
2. **Keyboard input**: gokeys handles keyboard events
3. **Mouse input**: goterm handles mouse events
4. **Resize events**: goterm handles resize events

### Integration Example

```go
// goterm manages terminal, gokeys reads from it
screen, err := goterm.Init(goterm.WithExternalKeyboard())
if err != nil {
    panic(err)
}
defer screen.Close()

// gokeys processes input
inp := input.New()
if err := inp.Start(); err != nil {
    panic(err)
}
defer inp.Stop()

// Both work together without conflicts
```

## Testing Support

### Test Helpers

#### `NewTestScreen(width, height int) *Screen`

Create screen for testing (no terminal required).

```go
screen := goterm.NewTestScreen(80, 24)
screen.DrawText(0, 0, "test", fg, bg, style)
// Verify cell contents
ch, _, _, _ := screen.GetCell(0, 0)
assert.Equal(t, 't', ch)
```

## Examples Repository

Full examples available at: `github.com/dshills/goterm/examples/`

- `examples/hello/` - Simple hello world
- `examples/colors/` - Color demonstration
- `examples/mouse/` - Mouse interaction
- `examples/game/` - Complete game with gokeys

## Contract Version

**Version**: 0.1.0
**Status**: Draft
**Last Updated**: 2025-10-18

This contract will be finalized before 1.0.0 release with stability guarantees.
