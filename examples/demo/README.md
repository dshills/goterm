# goterm Comprehensive Demo

This demo showcases all the features of the goterm terminal graphics library.

## Running the Demo

```bash
go run ./examples/demo
```

Or build and run:

```bash
go build -o goterm-demo ./examples/demo
./goterm-demo
```

## What's Demonstrated

The demo runs through 15 different screens, each highlighting specific features:

### 1. Welcome Screen
- Introduction to the library
- Display of terminal size
- Feature overview

### 2. 16 ANSI Colors
- Standard 8 colors (Black, Red, Green, Yellow, Blue, Magenta, Cyan, White)
- Bright variants (8-15)
- Foreground and background colors

### 3. 256 Color Palette
- Standard colors (0-15)
- 216-color cube (16-231)
- Grayscale ramp (232-255)

### 4. True Color (RGB)
- 24-bit RGB color support
- Custom color examples
- RGB gradients (Red, Green, Blue)

### 5. Color Conversions
- RGB to 256-color conversion
- 256-color to 16-color conversion
- Automatic color degradation examples

### 6. Text Styles
- All available text styles:
  - Normal, Bold, Dim, Italic
  - Underline, Blink (Slow/Rapid)
  - Reverse, Conceal, Strikethrough

### 7. Style Combinations
- Bitwise OR combinations
- Style operations (Set, Clear, Toggle, Has)
- Multiple styles applied simultaneously

### 8. Unicode Support
- ASCII characters
- Latin extended characters
- Greek, Cyrillic, Arabic, Hebrew
- CJK (Japanese, Chinese, Korean)
- Emoji support
- Mathematical symbols
- Currency symbols
- Arrows and decorative characters

### 9. Box Drawing
- Single-line boxes
- Double-line boxes
- Nested boxes
- Complex layouts with dividers

### 10. Screen Buffer
- Buffer size information
- Buffer operations (SetCell, GetCell, DrawText, Clear, Resize, Show)
- Pattern filling demonstration

### 11. Cell Manipulation
- Creating cells
- Cell equality checking
- Checkerboard pattern using individual cells

### 12. Dynamic Resize
- Screen resize capabilities
- Content preservation
- Reference grid

### 13. Color Gradients
- Red ↔ Green gradient
- Green ↔ Blue gradient
- Blue ↔ Red gradient
- Rainbow gradient
- Grayscale gradient
- 2D gradients

### 14. Animation
- Bouncing ball frame sequence
- Spinner animation frames
- Color wave
- Animation concepts

### 15. Final Showcase
- Complete feature summary
- Professional layout
- Rainbow divider
- Key features highlight

## Features Highlighted

### Color Support
- **16 Colors**: Standard ANSI colors with bright variants
- **256 Colors**: Extended palette with color cube and grayscale
- **True Color**: 24-bit RGB (16.7 million colors)
- **Conversions**: Automatic downgrade for terminal compatibility

### Text Styling
- **9 Style Flags**: Bold, Dim, Italic, Underline, Blink, Reverse, Conceal, Strikethrough
- **Combinable**: Use bitwise OR to combine multiple styles
- **Operations**: Set, Clear, Toggle, and Has methods

### Screen Management
- **Buffer-based**: All rendering goes through memory buffer first
- **Cell-based**: Individual cell control with character, colors, and style
- **Thread-safe**: RWMutex protection for concurrent access
- **Efficient**: Only changed cells need updating

### Character Support
- **Unicode**: Full Unicode support including emojis
- **CJK**: Chinese, Japanese, Korean characters
- **Box Drawing**: Single and double-line box drawing characters
- **Multi-byte**: Proper handling of multi-byte UTF-8 sequences

## Technical Details

Each demo screen:
- Clears the screen buffer
- Draws a styled header with the demo name
- Demonstrates specific features
- Draws a footer with instructions
- Shows the screen (renders to terminal)
- Waits 3 seconds before the next demo

The demo uses:
- `goterm.Init()` - Initialize terminal in raw mode
- `screen.Clear()` - Clear the screen buffer
- `screen.DrawText()` - Draw text strings
- `screen.SetCell()` - Set individual cells
- `screen.Show()` - Render buffer to terminal
- `screen.Close()` - Restore terminal to normal mode

## Requirements

- Terminal with color support
- Go 1.25.3 or later
- Unix-like system (macOS, Linux, BSD)

## Notes

- Some style effects (blink, conceal) may not be supported by all terminals
- True color support varies by terminal emulator
- The demo runs in sequence automatically
- Press Ctrl+C to exit at any time
- Best viewed in a terminal with at least 80x24 characters

## Duration

The complete demo takes approximately 45 seconds (15 screens × 3 seconds each).
