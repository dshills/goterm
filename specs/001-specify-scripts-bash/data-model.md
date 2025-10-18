# Data Model: Terminal Graphics Library

**Feature**: Terminal Graphics Library
**Branch**: `001-specify-scripts-bash`
**Date**: 2025-10-18

## Overview

This document defines the core data structures and their relationships for the terminal graphics library. The library operates on an in-memory screen buffer representing the terminal display, with no persistent storage required.

## Core Entities

### Screen Buffer

**Purpose**: Represents the terminal display as a 2D grid of cells that can be modified and rendered.

**Attributes**:
- Width (integer) - Number of columns in the terminal
- Height (integer) - Number of rows in the terminal
- Cells (2D array of Cell) - Grid representing each character position
- Dirty flags (boolean array) - Tracks which cells have changed since last render
- Previous buffer (Cell array) - For delta rendering comparison
- Current cursor position (row, column)
- Cursor visible (boolean)

**Relationships**:
- Contains: Multiple Cell instances (width × height)
- Uses: Terminal Capabilities for initialization
- Interacts with: Rendering pipeline for output

**State Transitions**:
```
Created → Initialized → Active → Closed
```

### Cell

**Purpose**: Individual unit in the screen buffer containing character data, colors, and styling.

**Attributes**:
- Character (rune) - Unicode character to display
- Foreground color (Color) - Text color
- Background color (Color) - Cell background color
- Style flags (Style bitmask) - Text attributes (bold, italic, etc.)
- Dirty flag (boolean) - Whether cell changed since last render

**Validation Rules**:
- Character must be valid Unicode code point
- Colors must be within supported color mode range
- Style flags must be valid combinations

**Memory Layout**:
```
Rune:       4 bytes (int32)
FG Color:   8 bytes (Color struct)
BG Color:   8 bytes (Color struct)
Style:      2 bytes (uint16 bitfield)
Dirty:      1 byte (bool)
Padding:    1 byte (alignment)
Total:      24 bytes per cell
```

### Color

**Purpose**: Represents colors in multiple modes with graceful degradation support.

**Attributes**:
- Mode (ColorMode enum) - Color representation mode
- Value (variant):
  - RGB (r, g, b uint8) - 24-bit true color
  - Index (uint8) - 8/16/256 color palette index
  - Default (special value) - Terminal default color

**Color Modes**:
```
Default     = 0  (terminal's default)
Basic16     = 1  (16 ANSI colors)
Extended256 = 2  (256-color palette)
TrueColor   = 3  (24-bit RGB)
```

**Validation Rules**:
- RGB values must be 0-255
- Palette index must be valid for mode (0-15 for Basic16, 0-255 for Extended256)
- Mode transitions must degrade gracefully

**Color Mapping**:
- TrueColor → Extended256: Convert RGB to nearest palette color
- Extended256 → Basic16: Map to closest ANSI color
- Basic16 → Default: Use terminal's color scheme

### Style

**Purpose**: Text attributes that can be combined via bitmasking.

**Style Flags** (bitmask):
```
None           = 0x0000  (no styling)
Bold           = 0x0001  (increased intensity)
Dim            = 0x0002  (decreased intensity)
Italic         = 0x0004  (slanted text)
Underline      = 0x0008  (underlined)
SlowBlink      = 0x0010  (blinking text)
RapidBlink     = 0x0020  (rapid blinking)
Reverse        = 0x0040  (swap fg/bg colors)
Conceal        = 0x0080  (hidden text)
Strikethrough  = 0x0100  (crossed-out text)
```

**Validation Rules**:
- Flags can be combined via bitwise OR
- Some combinations may not be supported by all terminals
- Mutually exclusive: Bold and Dim (last one wins)
- Mutually exclusive: SlowBlink and RapidBlink (last one wins)

**Terminal Compatibility**:
- Bold: Universally supported
- Italic: Modern terminals only
- Blink: Often disabled by users
- Reverse: Widely supported
- Strikethrough: Modern terminals only

### Mouse Event

**Purpose**: Information about mouse interactions delivered to application.

**Attributes**:
- Position (row, column integers) - Cell coordinates
- Button (MouseButton enum) - Which button triggered event
- Modifiers (Modifier bitmask) - Keyboard modifiers held
- Action (MouseAction enum) - Type of mouse event

**Mouse Buttons**:
```
Left      = 0
Middle    = 1
Right     = 2
None      = 3  (motion without button)
WheelUp   = 4
WheelDown = 5
WheelLeft = 6  (if supported)
WheelRight = 7  (if supported)
```

**Mouse Actions**:
```
Press    = 0  (button pressed down)
Release  = 1  (button released)
Motion   = 2  (mouse moved with button held)
Scroll   = 3  (wheel scrolled)
```

**Modifiers** (bitmask):
```
Shift = 0x01
Alt   = 0x02
Ctrl  = 0x04
```

**Validation Rules**:
- Position must be within current terminal bounds (0 ≤ row < height, 0 ≤ col < width)
- Button and Action must be valid combinations
- Scroll action only valid with Wheel buttons

### Resize Event

**Purpose**: Information about terminal dimension changes.

**Attributes**:
- New width (integer)
- New height (integer)
- Timestamp (time.Time) - When resize occurred

**Validation Rules**:
- Width and height must be positive (≥ 1)
- Typical minimum: 24 rows × 80 columns
- Typical maximum: 999 rows × 999 columns (SGR mouse protocol limit)

**Handling**:
```
Receive event → Validate dimensions → Resize screen buffer → Clear/Redraw
```

### Terminal Capabilities

**Purpose**: Information about what the current terminal supports.

**Attributes**:
- Color mode (ColorMode) - Highest supported color mode
- Mouse protocol (MouseProtocol enum) - Available mouse tracking
- Supports alternate screen (boolean)
- Supports cursor hide/show (boolean)
- Terminal type (string) - Value of $TERM
- Terminal emulator (string) - Detected emulator name

**Color Capability Detection**:
```
Check environment variables ($COLORTERM, $TERM) →
Query terminfo database →
Dynamic capability probing →
Fallback to built-in database
```

**Mouse Protocol Support**:
```
SGR (1006)    - Modern, unlimited coordinates (PREFERRED)
X11 (1000)    - Legacy, 223×223 limit (FALLBACK)
None          - No mouse support
```

## Data Relationships

### Composition Hierarchy

```
Screen Buffer
├── Cells [width × height]
│   ├── Character (rune)
│   ├── Foreground Color
│   ├── Background Color
│   └── Style
├── Previous Buffer [width × height] (for delta rendering)
└── Terminal Capabilities
```

### Event Flow

```
Terminal Input
    ↓
Mouse Events → Application Event Handler
Resize Events → Screen Buffer Resize → Application Event Handler
    ↓
Application updates Screen Buffer Cells
    ↓
Render Pipeline (delta rendering)
    ↓
Terminal Output
```

## Memory Considerations

### Screen Buffer Size

For an 80×24 terminal:
- Cells: 24 bytes/cell × 1920 cells = 46,080 bytes (~45 KB)
- Previous buffer: 46,080 bytes (~45 KB)
- Dirty flags: 1920 bytes (~2 KB)
- **Total**: ~92 KB per screen buffer

For a 200×60 terminal (modern large display):
- Cells: 24 bytes/cell × 12000 cells = 288,000 bytes (~281 KB)
- Previous buffer: 288,000 bytes (~281 KB)
- Dirty flags: 12,000 bytes (~12 KB)
- **Total**: ~574 KB per screen buffer

### Memory Optimization

**Pre-allocation**:
- Allocate screen buffer once at initialization
- Reuse existing buffer on resize when possible
- Use linear array instead of 2D slice for cache locality

**Delta Rendering**:
- Compare cells against previous buffer
- Build list of dirty cell indices
- Render only changed cells

**Escape Sequence Buffer**:
- Pre-allocate 8KB buffer for escape sequence generation
- Reuse buffer across frames
- Grow only if needed

## Implementation Notes

### Cell Access Pattern

**Linear indexing** (cache-friendly):
```go
index := row * width + col
cell := cells[index]
```

**Advantages**:
- Single allocation
- Better CPU cache utilization
- Simpler bounds checking

### Concurrent Access

**Thread Safety**:
- Screen buffer is NOT thread-safe by default
- Applications must synchronize access if using multiple goroutines
- Typical pattern: Single render goroutine owns buffer

**Recommended Pattern**:
```
Game Loop (single goroutine):
  1. Update game state
  2. Update screen buffer
  3. Render to terminal
  4. Sleep until next frame
```

### Platform-Specific Considerations

**Windows**:
- May need separate code path for Console API vs VT100 mode
- Console API uses different structures (CHAR_INFO vs escape sequences)

**Unix/Linux/macOS**:
- Unified code path using termios and VT100 sequences
- Consistent behavior across platforms

## Validation Summary

### Cell Validation
- ✅ Character is valid Unicode
- ✅ Colors within range for mode
- ✅ Style flags are valid

### Event Validation
- ✅ Mouse coordinates within bounds
- ✅ Button/Action combinations valid
- ✅ Resize dimensions positive

### Buffer Validation
- ✅ Dimensions match terminal size
- ✅ Cell array correctly allocated
- ✅ Dirty tracking synchronized

## Performance Targets

### Memory Access
- Cell access: O(1) via linear indexing
- Dirty check: O(n) where n = width × height
- Delta rendering: O(d) where d = number of dirty cells

### Memory Footprint
- Small terminal (80×24): ~100 KB
- Large terminal (200×60): ~600 KB
- Acceptable for games: ✅ Yes

### CPU Efficiency
- Delta rendering reduces I/O by 80-95%
- Pre-allocated buffers eliminate GC pressure
- Linear memory layout improves cache hits
