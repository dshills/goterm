package goterm

import "errors"

// Error types for terminal operations
var (
	// ErrNotATerminal indicates that stdout is not connected to a terminal
	ErrNotATerminal = errors.New("not a terminal")

	// ErrTerminalSetupFailed indicates that terminal initialization failed
	ErrTerminalSetupFailed = errors.New("terminal setup failed")

	// ErrTerminalRestoreFailed indicates that terminal restoration failed
	ErrTerminalRestoreFailed = errors.New("terminal restore failed")
)
