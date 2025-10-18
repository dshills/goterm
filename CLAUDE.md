# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**goterm** is a Go project located at `github.com/dshills/goterm` using Go 1.25.3.

This repository uses the Speckit workflow for feature development, which provides structured commands for specification-driven development through custom slash commands.

## Development Commands

### Building and Testing
```bash
# Build the project
go build ./...

# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a single test
go test -v -run TestName ./path/to/package

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Format code
go fmt ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

### Module Management
```bash
# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify

# Download dependencies
go mod download
```

## Speckit Workflow

This project uses Speckit for structured feature development. The workflow follows these stages:

1. **Specification** (`/speckit.specify`) - Create or update feature specifications from natural language descriptions
2. **Clarification** (`/speckit.clarify`) - Identify underspecified areas and ask targeted questions
3. **Planning** (`/speckit.plan`) - Execute implementation planning workflow
4. **Task Generation** (`/speckit.tasks`) - Generate actionable, dependency-ordered tasks
5. **Analysis** (`/speckit.analyze`) - Perform consistency and quality analysis across artifacts
6. **Implementation** (`/speckit.implement`) - Execute the implementation plan

### Key Speckit Commands
- `/speckit.constitution` - Create or update project principles and guidelines
- `/speckit.checklist` - Generate custom checklists for features

### Artifacts Location
- Specification templates: `.specify/templates/`
- Project constitution: `.specify/memory/constitution.md`
- Slash commands: `.claude/commands/`

## Code Architecture

This is a new Go project. As the codebase develops, follow Go standard project layout conventions:
- `cmd/` - Main applications for this project
- `pkg/` - Library code that can be used by external applications
- `internal/` - Private application and library code

## Go Best Practices for This Project

- Use standard Go tooling (`go fmt`, `go vet`, `go test`)
- Follow effective Go guidelines for naming and structure
- Write table-driven tests where appropriate
- Keep packages focused and cohesive
- Use Go modules for dependency management
