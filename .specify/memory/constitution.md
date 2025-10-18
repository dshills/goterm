<!--
Sync Impact Report - Constitution Update
Version: 0.0.0 → 1.0.0 (Initial ratification)
Date: 2025-10-18

Changes:
- Initial constitution created for goterm project
- Added 5 core principles: Go Idioms, Terminal Interface Design, Testing Discipline, Simplicity, API Stability
- Added Development Standards section
- Added Quality Gates section
- Governance rules established

Templates Status:
✅ plan-template.md - Constitution Check section compatible
✅ spec-template.md - Requirements sections align with principles
✅ tasks-template.md - Task organization supports testing discipline and parallel development
⚠️ No updates needed - templates are generic and compatible

Follow-up: None
-->

# goterm Constitution

## Core Principles

### I. Go Idioms and Best Practices

All code MUST follow established Go conventions and idioms:
- Use `gofmt` and `go vet` on all code before commits
- Follow effective Go guidelines for naming, structure, and error handling
- Prefer composition over inheritance; use interfaces for abstraction
- Handle errors explicitly; never ignore errors without documented justification
- Use standard library first; external dependencies require justification

**Rationale**: Go's strength lies in its simplicity and consistency. Following community conventions ensures maintainability and reduces cognitive load.

### II. Terminal Interface Design

Terminal interfaces MUST prioritize usability and composability:
- Support both human-readable and machine-parseable output formats
- Accept input via stdin, command arguments, and flags
- Write normal output to stdout, errors to stderr
- Return meaningful exit codes (0 = success, non-zero = specific failure)
- Enable piping and composition with other Unix tools
- Provide `--help` and `--version` flags for all commands

**Rationale**: Terminal tools are building blocks. They must be composable, scriptable, and follow Unix philosophy for maximum utility.

### III. Testing Discipline

Testing is mandatory for all production code:
- Write tests BEFORE or DURING implementation (not after)
- Achieve minimum 80% code coverage for new packages
- Use table-driven tests for multiple input scenarios
- Include both unit tests and integration tests where appropriate
- Ensure all tests pass before merging to main branch
- Tests MUST be fast (<1s for unit tests, <10s for integration)

**Rationale**: Go's testing tools make TDD natural. Early testing catches design issues and prevents regressions, ensuring code quality and maintainability.

### IV. Simplicity and Clarity (YAGNI)

Start simple and add complexity only when justified:
- Implement only what is needed for current requirements
- Prefer straightforward solutions over clever abstractions
- Keep functions small and focused (prefer <50 lines)
- Avoid premature optimization; measure before optimizing
- Document WHY for non-obvious decisions, not WHAT

**Rationale**: Go was designed for simplicity. Premature abstraction creates maintenance burden. Simple code is easier to understand, test, and modify.

### V. API Stability and Versioning

Public APIs require stability guarantees:
- Use semantic versioning (MAJOR.MINOR.PATCH)
- Document breaking changes explicitly in CHANGELOG
- Maintain backward compatibility within major versions
- Deprecate before removal (minimum one minor version warning)
- Mark experimental features clearly in documentation
- Use Go modules for dependency management

**Rationale**: Consumers of goterm libraries and tools need predictability. Semantic versioning and stability commitments build trust and enable safe upgrades.

## Development Standards

### Code Quality
- All code MUST pass `go fmt`, `go vet`, and configured linters
- All exported functions and types MUST have godoc comments
- Critical paths MUST have error handling with context
- Resource cleanup MUST use `defer` appropriately

### Performance
- Avoid unnecessary allocations in hot paths
- Use benchmarking (`go test -bench`) for performance-critical code
- Profile before optimizing (use `pprof` for CPU/memory analysis)
- Document performance characteristics for public APIs

### Security
- Never log sensitive data (credentials, tokens, PII)
- Validate all external inputs
- Use Go's crypto packages for cryptographic operations
- Keep dependencies updated for security patches

## Quality Gates

All features MUST pass these gates before merging:

1. **Code Review**: At least one approval from project maintainer
2. **Tests Pass**: All tests green on CI (unit + integration)
3. **Coverage**: New code meets 80% coverage minimum
4. **Lint Clean**: No warnings from `go vet` or configured linters
5. **Documentation**: Public APIs documented with examples where helpful
6. **Constitution Compliance**: Violates no principles, or violations justified in PR

### Complexity Justification

Violations of Simplicity (Principle IV) or introduction of non-standard patterns MUST be justified in pull requests with:
- Problem statement: What cannot be solved simply?
- Simpler alternatives considered and why they were insufficient
- Maintenance plan: How will this complexity be managed?

## Governance

### Amendment Process
- Constitution changes require documentation of rationale
- Breaking principle changes require team consensus
- Version bump follows semantic versioning:
  - **MAJOR**: Principle removal or backward-incompatible governance change
  - **MINOR**: New principle added or significant clarification
  - **PATCH**: Typo fixes, wording improvements, non-semantic updates

### Enforcement
- All pull requests MUST verify constitution compliance
- Reviewers MUST check Quality Gates before approval
- CI pipeline MUST enforce automated gates (tests, lint, coverage)

### Living Document
- Constitution should evolve with project needs
- Amendments MUST be committed with version update
- Use CLAUDE.md for development guidance that doesn't warrant constitution inclusion

**Version**: 1.0.0 | **Ratified**: 2025-10-18 | **Last Amended**: 2025-10-18
