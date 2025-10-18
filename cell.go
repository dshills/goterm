package goterm

// Cell represents a single character cell in the terminal screen buffer
type Cell struct {
	Ch    rune  // Character to display
	Fg    Color // Foreground color
	Bg    Color // Background color
	Style Style // Text styling flags
}

// NewCell creates a new cell with the specified attributes
func NewCell(ch rune, fg, bg Color, style Style) Cell {
	return Cell{
		Ch:    ch,
		Fg:    fg,
		Bg:    bg,
		Style: style,
	}
}

// Clear resets the cell to default (space character, default colors, no style)
func (c *Cell) Clear() {
	c.Ch = ' '
	c.Fg = ColorDefault()
	c.Bg = ColorDefault()
	c.Style = StyleNone
}

// Equal checks if two cells are identical
func (c Cell) Equal(other Cell) bool {
	return c.Ch == other.Ch &&
		c.Fg == other.Fg &&
		c.Bg == other.Bg &&
		c.Style == other.Style
}
