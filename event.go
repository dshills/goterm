package goterm

// Event is the interface for all terminal events
type Event interface {
	isEvent()
}

// MouseButton represents mouse button identifiers
type MouseButton int

// Mouse button constants
const (
	MouseLeft      MouseButton = 0 // Left mouse button
	MouseMiddle    MouseButton = 1 // Middle mouse button (wheel click)
	MouseRight     MouseButton = 2 // Right mouse button
	MouseNone      MouseButton = 3 // Motion without button
	MouseWheelUp   MouseButton = 4 // Mouse wheel scrolled up
	MouseWheelDown MouseButton = 5 // Mouse wheel scrolled down
)

// MouseAction represents the type of mouse event
type MouseAction int

// Mouse action constants
const (
	MousePress   MouseAction = 0 // Button pressed down
	MouseRelease MouseAction = 1 // Button released
	MouseMotion  MouseAction = 2 // Mouse moved with button held (drag)
	MouseScroll  MouseAction = 3 // Wheel scrolled
)

// Modifier represents keyboard modifier keys
type Modifier uint8

// Keyboard modifier constants (can be combined with bitwise OR)
const (
	ModShift Modifier = 1 << 0 // Shift key held
	ModAlt   Modifier = 1 << 1 // Alt/Option key held
	ModCtrl  Modifier = 1 << 2 // Control key held
)

// MouseEvent represents a mouse interaction
type MouseEvent struct {
	X, Y      int         // Cell coordinates
	Button    MouseButton // Which button
	Modifiers Modifier    // Keyboard modifiers held
	Action    MouseAction // Press/Release/Motion/Scroll
}

func (MouseEvent) isEvent() {}

// ResizeEvent represents a terminal dimension change
type ResizeEvent struct {
	Width  int // New width in columns
	Height int // New height in rows
}

func (ResizeEvent) isEvent() {}
