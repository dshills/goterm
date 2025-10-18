package goterm

import "golang.org/x/term"

// Terminal interface defines cross-platform terminal operations
// Platform-specific implementations will be in platform/ directory
type Terminal interface {
	// MakeRaw puts terminal into raw mode and returns previous state
	MakeRaw(fd int) (*term.State, error)

	// Restore restores terminal to previous state
	Restore(fd int, state *term.State) error

	// GetSize returns terminal dimensions (width, height)
	GetSize(fd int) (width, height int, err error)

	// IsTerminal checks if file descriptor is a terminal
	IsTerminal(fd int) bool
}
