# Quickstart Guide: Terminal Graphics Library

**Feature**: Terminal Graphics Library
**Branch**: `001-specify-scripts-bash`
**Date**: 2025-10-18

## Overview

This guide will help you get started with the terminal graphics library in 5 minutes. You'll learn to create a simple terminal application, render text and colors, handle mouse input, and integrate with gokeys for keyboard handling.

## Installation

```bash
go get github.com/dshills/goterm
```

## Hello World

Create `main.go`:

```go
package main

import (
    "github.com/dshills/goterm"
    "time"
)

func main() {
    // Initialize terminal
    screen, err := goterm.Init()
    if err != nil {
        panic(err)
    }
    defer screen.Close()

    // Draw text
    screen.Clear()
    screen.DrawText(10, 5, "Hello, World!",
        goterm.ColorYellow,
        goterm.ColorBlue,
        goterm.StyleBold)

    // Render to terminal
    screen.Show()

    // Wait 3 seconds
    time.Sleep(3 * time.Second)
}
```

Run it:
```bash
go run main.go
```

You should see yellow "Hello, World!" text on a blue background in the center of your terminal.

## Basic Concepts

### Screen Buffer

The screen is a 2D grid of cells. Each cell contains:
- A character (rune)
- Foreground color
- Background color
- Text styling (bold, italic, etc.)

### Rendering Pipeline

```
1. Clear or update cells
   screen.SetCell(x, y, 'A', fg, bg, style)

2. Show changes
   screen.Show()  // Only renders changed cells!
```

### Coordinate System

```
(0,0)  ────────────> X (columns)
  │
  │
  │
  ▼
  Y (rows)
```

## Drawing Shapes

```go
package main

import (
    "github.com/dshills/goterm"
    "time"
)

func main() {
    screen, err := goterm.Init()
    if err != nil {
        panic(err)
    }
    defer screen.Close()

    screen.Clear()

    // Draw a box
    drawBox(screen, 5, 2, 20, 10,
        goterm.ColorWhite,
        goterm.ColorBlack)

    screen.Show()
    time.Sleep(3 * time.Second)
}

func drawBox(s *goterm.Screen, x, y, w, h int, fg, bg goterm.Color) {
    // Top and bottom borders
    for col := x; col < x+w; col++ {
        s.SetCell(col, y, '─', fg, bg, goterm.StyleNone)
        s.SetCell(col, y+h-1, '─', fg, bg, goterm.StyleNone)
    }

    // Left and right borders
    for row := y; row < y+h; row++ {
        s.SetCell(x, row, '│', fg, bg, goterm.StyleNone)
        s.SetCell(x+w-1, row, '│', fg, bg, goterm.StyleNone)
    }

    // Corners
    s.SetCell(x, y, '┌', fg, bg, goterm.StyleNone)
    s.SetCell(x+w-1, y, '┐', fg, bg, goterm.StyleNone)
    s.SetCell(x, y+h-1, '└', fg, bg, goterm.StyleNone)
    s.SetCell(x+w-1, y+h-1, '┘', fg, bg, goterm.StyleNone)
}
```

## Colors

### Named Colors

```go
fg := goterm.ColorRed
bg := goterm.ColorBlack
```

Available named colors:
- `ColorBlack`, `ColorRed`, `ColorGreen`, `ColorYellow`
- `ColorBlue`, `ColorMagenta`, `ColorCyan`, `ColorWhite`

### RGB Colors (True Color)

```go
orange := goterm.ColorRGB(255, 165, 0)
purple := goterm.ColorRGB(128, 0, 128)
```

The library automatically degrades to 256-color or 16-color on terminals without true color support.

### Indexed Colors

```go
brightRed := goterm.ColorIndex(9)   // 0-15: ANSI colors
customColor := goterm.ColorIndex(45) // 16-255: Extended palette
```

## Text Styling

```go
// Single style
screen.SetCell(x, y, 'B', fg, bg, goterm.StyleBold)

// Multiple styles (combine with |)
style := goterm.StyleBold | goterm.StyleUnderline
screen.SetCell(x, y, 'C', fg, bg, style)
```

Available styles:
- `StyleNone` - No styling
- `StyleBold` - Bold/bright
- `StyleDim` - Dimmed
- `StyleItalic` - Italic
- `StyleUnderline` - Underlined
- `StyleReverse` - Swap fg/bg
- `StyleStrikethrough` - Crossed out

## Animation

```go
package main

import (
    "github.com/dshills/goterm"
    "time"
)

func main() {
    screen, err := goterm.Init()
    if err != nil {
        panic(err)
    }
    defer screen.Close()

    x, y := 10, 5
    dx := 1

    // 60 fps game loop
    ticker := time.NewTicker(16 * time.Millisecond)
    defer ticker.Stop()

    for i := 0; i < 300; i++ {  // 5 seconds at 60 fps
        <-ticker.C

        // Clear screen
        screen.Clear()

        // Update position
        x += dx
        w, _ := screen.Size()
        if x <= 0 || x >= w-1 {
            dx = -dx  // Bounce
        }

        // Draw
        screen.SetCell(x, y, '@',
            goterm.ColorYellow,
            goterm.ColorDefault(),
            goterm.StyleBold)

        // Render (only changed cells sent to terminal)
        screen.Show()
    }
}
```

## Mouse Input

```go
package main

import (
    "fmt"
    "github.com/dshills/goterm"
    "time"
)

func main() {
    screen, err := goterm.Init()
    if err != nil {
        panic(err)
    }
    defer screen.Close()

    screen.Clear()
    screen.DrawText(0, 0, "Click anywhere! Press Ctrl+C to quit.",
        goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
    screen.Show()

    ticker := time.NewTicker(16 * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {
        event := screen.PollEvent()
        if event == nil {
            continue
        }

        switch e := event.(type) {
        case *goterm.MouseEvent:
            if e.Action == goterm.MousePress {
                // Draw at click position
                screen.SetCell(e.X, e.Y, 'X',
                    goterm.ColorRed,
                    goterm.ColorDefault(),
                    goterm.StyleBold)
                screen.Show()
            }
        }
    }
}
```

## Keyboard Input (with gokeys)

```go
package main

import (
    "github.com/dshills/goterm"
    "github.com/dshills/gokeys/input"
    "time"
)

func main() {
    // Initialize graphics
    screen, err := goterm.Init(goterm.WithExternalKeyboard())
    if err != nil {
        panic(err)
    }
    defer screen.Close()

    // Initialize input
    inp := input.New()
    if err := inp.Start(); err != nil {
        panic(err)
    }
    defer inp.Stop()

    // Game state
    x, y := 10, 5

    // Game loop
    ticker := time.NewTicker(16 * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {
        // Handle input
        if event := inp.Next(); event != nil {
            switch event.Key {
            case input.KeyEscape:
                return
            case input.KeyLeft, input.KeyA:
                x--
            case input.KeyRight, input.KeyD:
                x++
            case input.KeyUp, input.KeyW:
                y--
            case input.KeyDown, input.KeyS:
                y++
            }

            // Keep in bounds
            w, h := screen.Size()
            if x < 0 {
                x = 0
            }
            if x >= w {
                x = w - 1
            }
            if y < 0 {
                y = 0
            }
            if y >= h {
                y = h - 1
            }
        }

        // Render
        screen.Clear()
        screen.DrawText(0, 0, "WASD or Arrows to move, ESC to quit",
            goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
        screen.SetCell(x, y, '@',
            goterm.ColorYellow,
            goterm.ColorDefault(),
            goterm.StyleBold)
        screen.Show()
    }
}
```

## Handling Resize

```go
package main

import (
    "fmt"
    "github.com/dshills/goterm"
    "time"
)

func main() {
    screen, err := goterm.Init()
    if err != nil {
        panic(err)
    }
    defer screen.Close()

    ticker := time.NewTicker(16 * time.Millisecond)
    defer ticker.Stop()

    for range ticker.C {
        // Check for resize events
        event := screen.PollEvent()
        if e, ok := event.(*goterm.ResizeEvent); ok {
            // Terminal was resized to e.Width × e.Height
            // Redraw everything
        }

        // Render
        screen.Clear()
        w, h := screen.Size()
        msg := fmt.Sprintf("Terminal size: %d × %d", w, h)
        screen.DrawText(0, 0, msg,
            goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleNone)
        screen.Show()
    }
}
```

## Complete Game Example

```go
package main

import (
    "fmt"
    "github.com/dshills/goterm"
    "github.com/dshills/gokeys/input"
    "time"
)

type Player struct {
    X, Y int
}

type Game struct {
    player  Player
    score   int
    running bool
}

func NewGame(w, h int) *Game {
    return &Game{
        player:  Player{X: w / 2, Y: h / 2},
        running: true,
    }
}

func (g *Game) HandleInput(event *input.Event) {
    switch event.Key {
    case input.KeyEscape, input.KeyQ:
        g.running = false
    case input.KeyLeft, input.KeyA:
        g.player.X--
    case input.KeyRight, input.KeyD:
        g.player.X++
    case input.KeyUp, input.KeyW:
        g.player.Y--
    case input.KeyDown, input.KeyS:
        g.player.Y++
    case input.KeySpace:
        g.score++
    }
}

func (g *Game) Update(w, h int) {
    // Keep player in bounds
    if g.player.X < 0 {
        g.player.X = 0
    }
    if g.player.X >= w {
        g.player.X = w - 1
    }
    if g.player.Y < 0 {
        g.player.Y = 0
    }
    if g.player.Y >= h {
        g.player.Y = h - 1
    }
}

func (g *Game) Render(screen *goterm.Screen) {
    screen.Clear()

    // Draw UI
    info := fmt.Sprintf("Score: %d | WASD=Move | SPACE=Score | ESC=Quit", g.score)
    screen.DrawText(0, 0, info,
        goterm.ColorWhite, goterm.ColorDefault(), goterm.StyleBold)

    // Draw player
    screen.SetCell(g.player.X, g.player.Y, '@',
        goterm.ColorYellow,
        goterm.ColorDefault(),
        goterm.StyleBold)

    screen.Show()
}

func main() {
    // Initialize graphics
    screen, err := goterm.Init(goterm.WithExternalKeyboard())
    if err != nil {
        panic(err)
    }
    defer screen.Close()

    // Initialize input
    inp := input.New()
    if err := inp.Start(); err != nil {
        panic(err)
    }
    defer inp.Stop()

    // Create game
    w, h := screen.Size()
    game := NewGame(w, h)

    // Game loop (60 fps)
    ticker := time.NewTicker(16 * time.Millisecond)
    defer ticker.Stop()

    for game.running {
        <-ticker.C

        // Handle input
        if event := inp.Next(); event != nil {
            game.HandleInput(event)
        }

        // Handle resize
        if event := screen.PollEvent(); event != nil {
            if e, ok := event.(*goterm.ResizeEvent); ok {
                w, h = e.Width, e.Height
            }
        }

        // Update
        game.Update(w, h)

        // Render
        game.Render(screen)
    }
}
```

## Performance Tips

### 1. Use Delta Rendering

`Show()` automatically renders only changed cells. Don't worry about optimization - it's built-in!

```go
// Only changed cells are sent to terminal
screen.SetCell(x, y, 'A', fg, bg, style)
screen.Show()  // Fast!
```

### 2. Batch Updates

Update multiple cells before calling `Show()`:

```go
// Good: Single Show() call
for i := 0; i < 100; i++ {
    screen.SetCell(i, 0, 'X', fg, bg, style)
}
screen.Show()  // Once

// Bad: Multiple Show() calls
for i := 0; i < 100; i++ {
    screen.SetCell(i, 0, 'X', fg, bg, style)
    screen.Show()  // Don't do this!
}
```

### 3. Limit Frame Rate

Use a ticker for consistent frame timing:

```go
ticker := time.NewTicker(16 * time.Millisecond)  // 60 fps
defer ticker.Stop()

for range ticker.C {
    // Game loop runs at 60 fps
}
```

### 4. Clear Smartly

Only clear if needed. If you're redrawing everything anyway, skip `Clear()`:

```go
// If drawing full screen every frame:
// screen.Clear()  // Not needed!
for y := 0; y < h; y++ {
    for x := 0; x < w; x++ {
        screen.SetCell(x, y, tile[y][x], fg, bg, style)
    }
}
screen.Show()
```

## Common Patterns

### Menu System

```go
type Menu struct {
    items    []string
    selected int
}

func (m *Menu) Render(screen *goterm.Screen) {
    for i, item := range m.items {
        style := goterm.StyleNone
        fg := goterm.ColorWhite
        if i == m.selected {
            style = goterm.StyleReverse
            fg = goterm.ColorYellow
        }
        screen.DrawText(10, 5+i, item, fg, goterm.ColorDefault(), style)
    }
}
```

### Double Buffering is Built-in

The library handles double buffering automatically. You don't need to manage buffers yourself!

### Partial Screen Updates

```go
// Only update game world area
for y := gameY; y < gameY+gameH; y++ {
    for x := gameX; x < gameX+gameW; x++ {
        screen.SetCell(x, y, world[y][x], fg, bg, style)
    }
}
// UI stays unchanged - won't be rerendered
screen.Show()
```

## Debugging

### Check Terminal Size

```go
w, h := screen.Size()
fmt.Printf("Terminal: %d×%d\n", w, h)
```

### Force Full Redraw

If display gets corrupted:

```go
screen.Sync()  // Force complete redraw
```

### Test Without Terminal

Use test screen for unit tests:

```go
screen := goterm.NewTestScreen(80, 24)
// Test your game logic without a real terminal
```

## Next Steps

- Read the [full API documentation](contracts/api.md)
- Check the [data model](data-model.md) for implementation details
- See `examples/` directory for more complex examples
- Read the [research document](research.md) for technical background

## Troubleshooting

### "not a terminal" error

Make sure you're running in an actual terminal, not redirecting output:

```bash
go run main.go           # ✅ Good
go run main.go > out.txt # ❌ Won't work (stdout redirected)
```

### Colors don't work

Check your terminal's color support:

```bash
echo $TERM              # Should be xterm-256color or better
echo $COLORTERM         # Should be "truecolor" for RGB colors
```

### Mouse doesn't work

Some terminals disable mouse by default. Try a modern terminal like:
- iTerm2 (macOS)
- Windows Terminal (Windows)
- gnome-terminal (Linux)
- alacritty (cross-platform)

### Performance is slow

- Limit frame rate to 30-60 fps
- Make sure you're not calling `Show()` in nested loops
- Use profiling: `go test -bench . -cpuprofile=cpu.prof`

## Happy Coding!

You're now ready to build amazing terminal games and applications. Have fun! 🎮
