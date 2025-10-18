# Implementation Plan: Terminal Graphics Library

**Branch**: `001-specify-scripts-bash` | **Date**: 2025-10-18 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-specify-scripts-bash/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a cross-platform terminal graphics library for creating terminal-based games and interactive applications. The library provides high-performance cell-based rendering, color support, mouse input, and seamless integration with the gokeys keyboard library. Target use case is terminal game development requiring 30-60 fps rendering with minimal CPU overhead.

## Technical Context

**Language/Version**: Go 1.25.3 (as specified in go.mod)
**Primary Dependencies**: golang.org/x/term + golang.org/x/sys (quasi-standard library) + custom minimal terminfo implementation
**Storage**: N/A (library operates on in-memory screen buffer)
**Testing**: Go standard testing (`go test`), table-driven tests, benchmarks for performance validation
**Target Platform**: Cross-platform (Linux, macOS, Windows) - multiple terminal emulators (xterm, gnome-terminal, iTerm2, Windows Terminal, Command Prompt, PowerShell)
**Project Type**: Single library project
**Performance Goals**: 60 fps full-screen rendering, sub-16ms frame times, delta rendering for efficiency
**Constraints**: <10% CPU at 60 fps, <5% CPU at 30 fps, flicker-free updates, sub-1s unit tests, sub-10s integration tests
**Scale/Scope**: Library (not application), support 5+ terminal emulators, handle terminals from 24x80 to modern large displays

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Go Idioms and Best Practices
- ✅ **PASS**: Will use standard Go project layout, gofmt, go vet
- ✅ **PASS**: Standard library first approach (syscalls for terminal control where needed)
- ✅ **PASS**: Composition over inheritance (interfaces for platform abstraction)
- ✅ **PASS**: Explicit error handling required for all terminal operations

### Principle II: Terminal Interface Design
- ⚠️ **LIBRARY EXEMPTION**: This IS a library for terminal interfaces, not a CLI tool
- ✅ **APPLICABLE**: Library will enable applications to follow this principle
- ✅ **PASS**: Design supports building composable terminal applications

### Principle III: Testing Discipline
- ✅ **PASS**: TDD approach required - tests before/during implementation
- ✅ **PASS**: 80% coverage minimum, table-driven tests for cross-platform scenarios
- ✅ **PASS**: Unit tests (<1s) for logic, integration tests (<10s) for rendering
- ✅ **PASS**: Performance benchmarks using `go test -bench`

### Principle IV: Simplicity and Clarity (YAGNI)
- ✅ **PASS**: Start with P1 (basic rendering) MVP
- ✅ **PASS**: Functions <50 lines, straightforward rendering pipeline
- ⚠️ **JUSTIFICATION NEEDED**: Platform abstraction complexity required for cross-platform support
  - **Why**: Windows, Unix, and terminal emulator differences require abstraction layer
  - **Simpler alternative rejected**: Single-platform library insufficient per requirements
  - **Mitigation**: Use interfaces to keep platform-specific code isolated and testable

### Principle V: API Stability and Versioning
- ✅ **PASS**: Will use semantic versioning (starting 0.1.0 for initial development)
- ✅ **PASS**: Public API will be carefully designed with backward compatibility in mind
- ✅ **PASS**: Go modules already configured in project

### Development Standards Compliance
- ✅ **Code Quality**: gofmt, go vet, golangci-lint configured
- ✅ **Performance**: Benchmarking required for hot paths (rendering, delta calculation)
- ✅ **Security**: Validate terminal escape sequences, handle malformed input safely

### Quality Gates Compliance
- ✅ **Tests Pass**: CI will run `go test ./...`
- ✅ **Coverage**: 80% minimum tracked
- ✅ **Lint Clean**: go vet and linters in CI
- ✅ **Documentation**: godoc for all exported types/functions

### Complexity Justification Required
| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Platform abstraction layer | Windows uses different terminal APIs than Unix (Console API vs termios) | Single codebase impossible - fundamental OS differences in terminal control |
| Multiple color mode support | Terminals have varying capabilities (16/256/truecolor) | Users expect colors to work everywhere - graceful degradation required |

**Gate Status**: ✅ **PASS** - Complexity justified by cross-platform requirement

## Project Structure

### Documentation (this feature)

```
specs/001-specify-scripts-bash/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   └── api.md           # Public API contract
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```
goterm/                  # Root of library
├── screen.go            # Screen buffer and rendering coordination
├── cell.go              # Cell representation (char, colors, style)
├── color.go             # Color type and palette management
├── style.go             # Text styling (bold, italic, etc.)
├── event.go             # Event types (mouse, resize)
├── terminal.go          # Terminal interface (cross-platform abstraction)
├── platform/            # Platform-specific implementations
│   ├── unix.go          # Unix/Linux/macOS terminal control (termios, ioctl)
│   ├── windows.go       # Windows Console API implementation
│   └── capabilities.go  # Terminal capability detection
├── mouse.go             # Mouse event handling
├── buffer.go            # Double buffering and delta rendering
└── examples/            # Example applications
    ├── hello/           # Simple hello world
    ├── colors/          # Color demonstration
    └── game/            # Simple game with gokeys integration

tests/
├── unit/                # Unit tests (isolated logic)
│   ├── cell_test.go
│   ├── color_test.go
│   ├── buffer_test.go
│   └── delta_test.go
├── integration/         # Integration tests (actual terminal operations)
│   ├── render_test.go
│   ├── mouse_test.go
│   └── platform_test.go
└── benchmark/           # Performance benchmarks
    ├── render_bench.go
    └── delta_bench.go
```

**Structure Decision**: Single Go library project at repository root. This follows standard Go library conventions where the library itself is at the module root, not in a `pkg/` subdirectory. Platform-specific code uses Go build tags to compile only the necessary implementation. Examples directory demonstrates usage patterns. Tests are organized by type (unit/integration/benchmark) for clear separation of concerns.

## Complexity Tracking

*Filled because Constitution Check identified justified complexity*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Platform abstraction (Principle IV) | Windows Console API fundamentally different from Unix termios/pty | Cannot achieve cross-platform support without platform-specific implementations |
| Color mode detection and adaptation | Terminal capabilities vary widely (16/256/truecolor) | Hard-coding one mode fails on many terminals; library must adapt to environment |
| Event multiplexing | Need to handle stdin (keyboard from gokeys), mouse, resize signals simultaneously | Simple blocking I/O insufficient for interactive applications requiring 60fps |

**Maintenance Plan**:
- Platform-specific code isolated behind interfaces (`terminal.go` defines contract)
- Build tags keep platform code separate (`// +build unix`, `// +build windows`)
- Extensive testing on each platform to catch regressions
- Keep complexity contained to `platform/` package - rest of library platform-agnostic
