package goterm

// Style represents text styling attributes as a bitmask
type Style uint16

// Style constants that can be combined with bitwise OR
const (
	StyleNone          Style = 0
	StyleBold          Style = 1 << 0 // Increased intensity / bright
	StyleDim           Style = 1 << 1 // Decreased intensity
	StyleItalic        Style = 1 << 2 // Slanted text
	StyleUnderline     Style = 1 << 3 // Underlined
	StyleSlowBlink     Style = 1 << 4 // Blinking text (slow)
	StyleRapidBlink    Style = 1 << 5 // Blinking text (rapid)
	StyleReverse       Style = 1 << 6 // Swap foreground/background colors
	StyleConceal       Style = 1 << 7 // Hidden text
	StyleStrikethrough Style = 1 << 8 // Crossed-out text
)

// Has checks if a style flag is set
func (s Style) Has(flag Style) bool {
	return s&flag != 0
}

// Set adds a style flag
func (s Style) Set(flag Style) Style {
	return s | flag
}

// Clear removes a style flag
func (s Style) Clear(flag Style) Style {
	return s &^ flag
}

// Toggle flips a style flag
func (s Style) Toggle(flag Style) Style {
	return s ^ flag
}

// ansiCode returns the ANSI escape codes for this style
func (s Style) ansiCode() string {
	if s == StyleNone {
		return ""
	}

	codes := ""
	if s.Has(StyleBold) {
		codes += "\x1b[1m"
	}
	if s.Has(StyleDim) {
		codes += "\x1b[2m"
	}
	if s.Has(StyleItalic) {
		codes += "\x1b[3m"
	}
	if s.Has(StyleUnderline) {
		codes += "\x1b[4m"
	}
	if s.Has(StyleSlowBlink) {
		codes += "\x1b[5m"
	}
	if s.Has(StyleRapidBlink) {
		codes += "\x1b[6m"
	}
	if s.Has(StyleReverse) {
		codes += "\x1b[7m"
	}
	if s.Has(StyleConceal) {
		codes += "\x1b[8m"
	}
	if s.Has(StyleStrikethrough) {
		codes += "\x1b[9m"
	}

	return codes
}
