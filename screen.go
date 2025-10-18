package goterm

import (
	"fmt"
	"io"
	"os"
	"sync"

	"golang.org/x/term"
)

// Screen represents a terminal screen buffer with cells for rendering
type Screen struct {
	width  int
	height int
	cells  []Cell
	mu     sync.RWMutex

	// Terminal state
	fd       int
	oldState *term.State
	out      io.Writer
}

// NewScreen creates a new screen buffer with the specified dimensions
// Panics if width or height are <= 0
func NewScreen(width, height int) *Screen {
	if width <= 0 || height <= 0 {
		panic(fmt.Sprintf("invalid screen dimensions: width=%d, height=%d", width, height))
	}

	s := &Screen{
		width:  width,
		height: height,
		cells:  make([]Cell, width*height),
		out:    os.Stdout,
	}

	// Initialize all cells to defaults
	s.Clear()

	return s
}

// Size returns the current screen dimensions
func (s *Screen) Size() (width, height int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.width, s.height
}

// SetCell sets the cell at the specified position
// Does nothing if x, y are out of bounds
func (s *Screen) SetCell(x, y int, cell Cell) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if x < 0 || y < 0 || x >= s.width || y >= s.height {
		return
	}

	s.cells[y*s.width+x] = cell
}

// GetCell returns the cell at the specified position
// Returns a default empty cell if x, y are out of bounds
func (s *Screen) GetCell(x, y int) Cell {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if x < 0 || y < 0 || x >= s.width || y >= s.height {
		return NewCell(' ', ColorDefault(), ColorDefault(), StyleNone)
	}

	return s.cells[y*s.width+x]
}

// Clear resets all cells to their default state
func (s *Screen) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	defaultCell := NewCell(' ', ColorDefault(), ColorDefault(), StyleNone)
	for i := range s.cells {
		s.cells[i] = defaultCell
	}
}

// DrawText draws text at the specified position with the given colors and style
// Text that extends beyond the screen width is clipped
func (s *Screen) DrawText(x, y int, text string, fg, bg Color, style Style) {
	runes := []rune(text)
	for i, ch := range runes {
		cell := NewCell(ch, fg, bg, style)
		s.SetCell(x+i, y, cell)
	}
}

// Resize changes the screen dimensions
// Content is preserved where it fits in the new dimensions
func (s *Screen) Resize(width, height int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if width <= 0 || height <= 0 {
		return
	}

	// Create new buffer
	newCells := make([]Cell, width*height)
	defaultCell := NewCell(' ', ColorDefault(), ColorDefault(), StyleNone)
	for i := range newCells {
		newCells[i] = defaultCell
	}

	// Copy existing content that fits
	minWidth := width
	if s.width < minWidth {
		minWidth = s.width
	}
	minHeight := height
	if s.height < minHeight {
		minHeight = s.height
	}

	for y := 0; y < minHeight; y++ {
		for x := 0; x < minWidth; x++ {
			oldIdx := y*s.width + x
			newIdx := y*width + x
			newCells[newIdx] = s.cells[oldIdx]
		}
	}

	s.width = width
	s.height = height
	s.cells = newCells
}

// Show renders the screen buffer to the terminal
// This is where the actual terminal escape sequences are written
func (s *Screen) Show() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Move cursor to home position
	if _, err := fmt.Fprint(s.out, "\x1b[H"); err != nil {
		return fmt.Errorf("failed to move cursor: %w", err)
	}

	var lastFg, lastBg Color
	var lastStyle Style
	needsReset := false

	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			cell := s.cells[y*s.width+x]

			// Output color/style changes only when needed
			if cell.Fg != lastFg || cell.Bg != lastBg || cell.Style != lastStyle || needsReset {
				// Reset before applying new attributes
				if _, err := fmt.Fprint(s.out, "\x1b[0m"); err != nil {
					return fmt.Errorf("failed to reset attributes: %w", err)
				}

				// Set foreground color
				if cell.Fg.Mode() != ColorModeDefault {
					if _, err := fmt.Fprint(s.out, cell.Fg.ansiCode(true)); err != nil {
						return fmt.Errorf("failed to set foreground color: %w", err)
					}
				}

				// Set background color
				if cell.Bg.Mode() != ColorModeDefault {
					if _, err := fmt.Fprint(s.out, cell.Bg.ansiCode(false)); err != nil {
						return fmt.Errorf("failed to set background color: %w", err)
					}
				}

				// Set style
				if cell.Style != StyleNone {
					if _, err := fmt.Fprint(s.out, cell.Style.ansiCode()); err != nil {
						return fmt.Errorf("failed to set style: %w", err)
					}
				}

				lastFg = cell.Fg
				lastBg = cell.Bg
				lastStyle = cell.Style
				needsReset = false
			}

			// Output the character
			if _, err := fmt.Fprint(s.out, string(cell.Ch)); err != nil {
				return fmt.Errorf("failed to write character: %w", err)
			}
		}

		// Move to next line if not last line
		if y < s.height-1 {
			if _, err := fmt.Fprint(s.out, "\r\n"); err != nil {
				return fmt.Errorf("failed to write newline: %w", err)
			}
		}
	}

	// Reset attributes at end
	if _, err := fmt.Fprint(s.out, "\x1b[0m"); err != nil {
		return fmt.Errorf("failed to reset final attributes: %w", err)
	}

	return nil
}

// Sync flushes any buffered output to the terminal
func (s *Screen) Sync() error {
	if f, ok := s.out.(*os.File); ok {
		return f.Sync()
	}
	return nil
}

// Close restores the terminal to its previous state
func (s *Screen) Close() error {
	if s.oldState != nil && s.fd > 0 {
		return term.Restore(s.fd, s.oldState)
	}
	return nil
}

// Init initializes the terminal for screen rendering
// Returns a Screen initialized to the terminal's current size
func Init() (*Screen, error) {
	fd := int(os.Stdout.Fd())

	// Check if stdout is a terminal
	if !term.IsTerminal(fd) {
		return nil, ErrNotATerminal
	}

	// Get terminal size
	width, height, err := term.GetSize(fd)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTerminalSetupFailed, err)
	}

	// Put terminal in raw mode
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTerminalSetupFailed, err)
	}

	screen := NewScreen(width, height)
	screen.fd = fd
	screen.oldState = oldState

	// Clear screen and hide cursor
	if _, err := fmt.Fprint(screen.out, "\x1b[2J\x1b[H\x1b[?25l"); err != nil {
		// Best effort cleanup
		_ = term.Restore(fd, oldState)
		return nil, fmt.Errorf("%w: failed to initialize screen: %v", ErrTerminalSetupFailed, err)
	}

	return screen, nil
}
