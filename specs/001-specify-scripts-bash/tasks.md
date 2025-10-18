# Tasks: Terminal Graphics Library

**Input**: Design documents from `/specs/001-specify-scripts-bash/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are REQUIRED per constitution (TDD with 80% coverage minimum). Tests must be written BEFORE implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions
- **Go library**: Files at repository root (e.g., `screen.go`, `cell.go`)
- **Platform-specific**: `platform/unix.go`, `platform/windows.go`
- **Tests**: `tests/unit/`, `tests/integration/`, `tests/benchmark/`
- **Examples**: `examples/hello/`, `examples/colors/`, `examples/game/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Initialize Go module at repository root with go.mod already present
- [X] T002 Add golang.org/x/term dependency to go.mod
- [X] T003 [P] Add golang.org/x/sys dependency to go.mod
- [X] T004 [P] Create platform/ directory for OS-specific implementations
- [X] T005 [P] Create tests/unit/ directory for unit tests
- [X] T006 [P] Create tests/integration/ directory for integration tests
- [X] T007 [P] Create tests/benchmark/ directory for performance benchmarks
- [X] T008 [P] Create examples/ directory for example applications
- [X] T009 [P] Configure golangci-lint with .golangci.yml at repository root
- [X] T010 [P] Create Makefile with test, bench, lint, and fmt targets at repository root

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core types and interfaces that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T011 [P] Create errors.go with error types (ErrNotATerminal, ErrTerminalSetupFailed, ErrTerminalRestoreFailed) at repository root
- [X] T012 [P] Create color.go with Color type, ColorMode enum, and color constructors (ColorDefault, ColorRGB, ColorIndex) at repository root
- [X] T013 [P] Create style.go with Style type (bitmask) and style constants (StyleBold, StyleItalic, etc.) at repository root
- [X] T014 [P] Create cell.go with Cell struct (rune, fg Color, bg Color, style Style) at repository root
- [X] T015 [P] Create event.go with Event interface, MouseEvent, ResizeEvent types at repository root
- [X] T016 Create terminal.go with Terminal interface defining cross-platform contract (MakeRaw, Restore, GetSize, etc.) at repository root
- [X] T017 [P] Write unit tests for Color type in tests/unit/color_test.go (test ColorRGB, ColorIndex, mode detection)
- [X] T018 [P] Write unit tests for Style type in tests/unit/style_test.go (test bitmask operations, combinations)
- [X] T019 [P] Write unit tests for Cell type in tests/unit/cell_test.go (test cell creation, attribute setting)

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Basic Screen Rendering (Priority: P1) 🎯 MVP

**Goal**: Enable drawing text and graphics to terminal with colors, positioning, and styling

**Independent Test**: Create simple program that draws colored text at specific coordinates and verify output in terminal emulators

### Tests for User Story 1 ⚠️

**NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T020 [P] [US1] Write unit tests for Screen buffer allocation in tests/unit/screen_test.go
- [X] T021 [P] [US1] Write unit tests for SetCell/GetCell operations in tests/unit/screen_test.go
- [X] T022 [P] [US1] Write unit tests for Clear() functionality in tests/unit/screen_test.go
- [X] T023 [P] [US1] Write integration test for basic rendering in tests/integration/render_test.go

### Implementation for User Story 1

- [X] T024 [US1] Create screen.go with Screen struct (width, height, cells []Cell, dirty []bool) at repository root
- [X] T025 [US1] Implement Init() function in screen.go (terminal detection, size query, buffer allocation)
- [X] T026 [US1] Implement Close() method in screen.go (cleanup, terminal restoration)
- [X] T027 [US1] Implement Size() method in screen.go (return width, height)
- [X] T028 [US1] Implement SetCell(x, y int, ch rune, fg, bg Color, style Style) method in screen.go
- [X] T029 [US1] Implement GetCell(x, y int) method in screen.go
- [X] T030 [US1] Implement Clear() method in screen.go (reset all cells to defaults, mark dirty)
- [X] T031 [US1] Implement DrawText(x, y int, text string, fg, bg Color, style Style) method in screen.go
- [X] T032 [US1] Create buffer.go with rendering buffer management (escape sequence generation) at repository root
- [X] T033 [US1] Implement Show() method in screen.go (render dirty cells to terminal via buffer)
- [X] T034 [US1] Implement Sync() method in screen.go (force full screen redraw)
- [X] T035 [US1] Add godoc comments for all exported functions in screen.go
- [X] T036 [US1] Verify all User Story 1 tests pass with go test ./...

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Cross-Platform Terminal Support (Priority: P2)

**Goal**: Make library work consistently across Linux, macOS, Windows and multiple terminal emulators

**Independent Test**: Run same application on Windows, macOS, Linux terminals and verify identical behavior

### Tests for User Story 2 ⚠️

- [ ] T037 [P] [US2] Write unit tests for terminfo capability detection in tests/unit/capabilities_test.go
- [ ] T038 [P] [US2] Write unit tests for Unix platform implementation in tests/unit/platform_unix_test.go
- [ ] T039 [P] [US2] Write unit tests for Windows platform implementation in tests/unit/platform_windows_test.go
- [ ] T040 [P] [US2] Write integration tests for cross-platform terminal init in tests/integration/platform_test.go

### Implementation for User Story 2

- [ ] T041 [P] [US2] Create platform/capabilities.go with terminal capability detection (color mode, TERM env vars) at repository root
- [ ] T042 [P] [US2] Create platform/terminfo.go with minimal terminfo implementation (built-in database for common terminals) at repository root
- [ ] T043 [US2] Create platform/unix.go with Unix terminal implementation (uses golang.org/x/term, build tag: unix) at repository root
- [ ] T044 [US2] Create platform/windows.go with Windows terminal implementation (VT100 mode + Console API fallback, build tag: windows) at repository root
- [ ] T045 [US2] Update screen.go Init() to use platform-specific terminal initialization
- [ ] T046 [US2] Implement resize event handling with SIGWINCH signal (Unix) and console events (Windows) in screen.go
- [ ] T047 [US2] Update screen.go to handle terminal capability detection and graceful degradation
- [ ] T048 [US2] Add escape sequence generation based on detected terminal capabilities in buffer.go
- [ ] T049 [US2] Verify all User Story 2 tests pass with go test ./...

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently across platforms

---

## Phase 5: User Story 3 - High-Performance Rendering (Priority: P3)

**Goal**: Achieve 30-60 fps rendering with <10% CPU usage through delta rendering and optimization

**Independent Test**: Create animation updating full screen at 60 fps, measure CPU usage and visual smoothness

### Tests for User Story 3 ⚠️

- [ ] T050 [P] [US3] Write unit tests for dirty cell tracking in tests/unit/buffer_test.go
- [ ] T051 [P] [US3] Write unit tests for delta rendering algorithm in tests/unit/delta_test.go
- [ ] T052 [P] [US3] Write benchmark for full-screen rendering in tests/benchmark/render_bench.go
- [ ] T053 [P] [US3] Write benchmark for delta rendering in tests/benchmark/delta_bench.go

### Implementation for User Story 3

- [ ] T054 [US3] Implement dirty cell tracking in screen.go (mark cells dirty on SetCell, compare with previous buffer)
- [ ] T055 [US3] Implement delta rendering in buffer.go (only generate escape sequences for dirty cells)
- [ ] T056 [US3] Add double buffering to screen.go (current and previous screen buffers for comparison)
- [ ] T057 [US3] Optimize escape sequence batching in buffer.go (combine multiple SGR codes, use relative cursor positioning)
- [ ] T058 [US3] Pre-allocate buffers in buffer.go to avoid allocations in render loop (use sync.Pool if needed)
- [ ] T059 [US3] Add buffered I/O in buffer.go (use bufio.Writer with 8KB buffer, single flush per frame)
- [ ] T060 [US3] Profile rendering with pprof and optimize hot paths to meet <16ms target
- [ ] T061 [US3] Verify benchmarks meet performance targets (60 fps, <10% CPU) with go test -bench ./...

**Checkpoint**: All rendering should now be high-performance with delta updates and minimal CPU usage

---

## Phase 6: User Story 4 - Color and Styling Support (Priority: P4)

**Goal**: Support full color spectrum (16/256/truecolor) and text styles (bold, italic, underline, etc.)

**Independent Test**: Render text with all color modes and all style attributes, verify correct display

### Tests for User Story 4 ⚠️

- [ ] T062 [P] [US4] Write unit tests for color mode detection in tests/unit/color_test.go (test RGB → 256 → 16 degradation)
- [ ] T063 [P] [US4] Write unit tests for style rendering in tests/unit/style_test.go (test SGR code generation)
- [ ] T064 [P] [US4] Write integration test for color rendering in tests/integration/color_test.go

### Implementation for User Story 4

- [ ] T065 [P] [US4] Implement RGB to 256-color conversion in color.go (color cube + grayscale mapping)
- [ ] T066 [P] [US4] Implement 256-color to 16-color conversion in color.go (nearest ANSI color)
- [ ] T067 [US4] Update buffer.go to generate appropriate SGR codes based on detected color mode
- [ ] T068 [US4] Implement style SGR code generation in buffer.go (bold=1, italic=3, underline=4, etc.)
- [ ] T069 [US4] Add style combination support in buffer.go (combine multiple styles in single SGR sequence)
- [ ] T070 [US4] Update capability detection in platform/capabilities.go to detect truecolor support ($COLORTERM env var)
- [ ] T071 [US4] Verify all User Story 4 tests pass and colors render correctly

**Checkpoint**: Color and styling should work across different terminal capabilities with graceful degradation

---

## Phase 7: User Story 5 - Mouse and Keyboard Input Integration (Priority: P5)

**Goal**: Provide mouse support (clicks, movement, wheel) that integrates with gokeys for complete input handling

**Independent Test**: Create program responding to mouse clicks/movements, verify coordinates and button states

### Tests for User Story 5 ⚠️

- [ ] T072 [P] [US5] Write unit tests for SGR mouse event parsing in tests/unit/mouse_test.go
- [ ] T073 [P] [US5] Write unit tests for X11 mouse event parsing in tests/unit/mouse_test.go
- [ ] T074 [P] [US5] Write integration test for mouse input in tests/integration/mouse_test.go

### Implementation for User Story 5

- [ ] T075 [US5] Create mouse.go with mouse protocol initialization (SGR mode, button-event tracking) at repository root
- [ ] T076 [US5] Implement SGR mouse event parser in mouse.go (parse ESC[<Cb;Cx;CyM/m format)
- [ ] T077 [US5] Implement X11 mouse event parser in mouse.go (fallback for legacy terminals, parse ESC[Mbxy format)
- [ ] T078 [US5] Add mouse event queue to screen.go (buffered channel for events)
- [ ] T079 [US5] Implement PollEvent() method in screen.go (non-blocking event retrieval)
- [ ] T080 [US5] Implement Events() method in screen.go (return event channel)
- [ ] T081 [US5] Add mouse enable/disable in screen.go Init() and Close() (send escape sequences)
- [ ] T082 [US5] Start event reader goroutine in screen.go Init() (read stdin, parse events, send to channel)
- [ ] T083 [US5] Verify all User Story 5 tests pass and mouse events work correctly

**Checkpoint**: Mouse input should work alongside keyboard input without conflicts

---

## Phase 8: User Story 6 - gokeys Compatibility (Priority: P6)

**Goal**: Seamless integration with gokeys library - coordinate terminal settings without conflicts

**Independent Test**: Create application using both gokeys and goterm, verify no initialization conflicts or event interference

### Tests for User Story 6 ⚠️

- [ ] T084 [P] [US6] Write integration test for gokeys compatibility in tests/integration/gokeys_test.go
- [ ] T085 [P] [US6] Write integration test for coordinated terminal settings in tests/integration/gokeys_test.go

### Implementation for User Story 6

- [ ] T086 [US6] Add WithExternalKeyboard() option to screen.go Init() (skip keyboard event reading)
- [ ] T087 [US6] Update terminal initialization in screen.go to set OPOST flag correctly (needed for graphics, different from gokeys)
- [ ] T088 [US6] Ensure screen.go only reads from stdin for mouse events, not keyboard (gokeys handles keyboard)
- [ ] T089 [US6] Add RawFd() method to screen.go for external libraries to access terminal file descriptor
- [ ] T090 [US6] Update resize event handling to be observable by external libraries
- [ ] T091 [US6] Create example in examples/game/ demonstrating goterm + gokeys integration
- [ ] T092 [US6] Verify all User Story 6 tests pass and integration works without conflicts

**Checkpoint**: All user stories should now be independently functional with gokeys integration complete

---

## Phase 9: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories, examples, and documentation

- [ ] T093 [P] Create example in examples/hello/ showing basic "Hello World" with colors
- [ ] T094 [P] Create example in examples/colors/ demonstrating all color modes and styles
- [ ] T095 [P] Update examples/game/ with complete game using mouse, gokeys, and rendering
- [ ] T096 [P] Add godoc comments to all exported types and functions across all files
- [ ] T097 [P] Create README.md at repository root with installation, quick start, and examples
- [ ] T098 Run go fmt ./... to format all code
- [ ] T099 Run go vet ./... and fix any issues
- [ ] T100 Run golangci-lint run and fix any issues
- [ ] T101 Verify 80% test coverage with go test -cover ./...
- [ ] T102 Run all benchmarks and verify performance targets met (60 fps, <10% CPU)
- [ ] T103 Test manually on 5 different terminals (xterm, iTerm2, Windows Terminal, gnome-terminal, alacritty)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-8)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3 → P4 → P5 → P6)
- **Polish (Phase 9)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1 - Basic Rendering)**: Can start after Foundational (Phase 2) - No dependencies on other stories ✅ **MVP**
- **User Story 2 (P2 - Cross-Platform)**: Can start after Foundational - Enhances US1 but independently testable
- **User Story 3 (P3 - Performance)**: Can start after Foundational - Optimizes US1 but independently testable
- **User Story 4 (P4 - Colors/Styling)**: Can start after Foundational - Enhances US1 but independently testable
- **User Story 5 (P5 - Mouse Input)**: Can start after Foundational - Independent input system
- **User Story 6 (P6 - gokeys Integration)**: Can start after Foundational - Coordination layer, integrates with US5

### Within Each User Story

- Tests MUST be written and FAIL before implementation (TDD requirement from constitution)
- Tests can run in parallel if marked [P] (different files, no dependencies)
- Implementation tasks run after tests are written
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, multiple user stories can be worked on in parallel:
  - Developer A: User Story 1
  - Developer B: User Story 2
  - Developer C: User Story 4
- Tests within a user story marked [P] can run in parallel
- Polish tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Write unit tests for Screen buffer allocation in tests/unit/screen_test.go"
Task: "Write unit tests for SetCell/GetCell operations in tests/unit/screen_test.go"
Task: "Write unit tests for Clear() functionality in tests/unit/screen_test.go"
Task: "Write integration test for basic rendering in tests/integration/render_test.go"

# After tests are written and failing, implementation proceeds sequentially (tasks depend on each other)
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Basic Rendering)
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Verify all tests pass, coverage ≥80%, benchmarks acceptable
6. You now have a working terminal graphics library! 🎉

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → **MVP Released!**
3. Add User Story 2 → Test independently → Cross-platform support added
4. Add User Story 3 → Test independently → Performance optimized
5. Add User Story 4 → Test independently → Rich colors/styling added
6. Add User Story 5 → Test independently → Mouse input added
7. Add User Story 6 → Test independently → gokeys integration complete
8. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Basic Rendering) - PRIORITY
   - Developer B: User Story 4 (Colors/Styling) - Can work in parallel
   - Developer C: User Story 2 (Cross-Platform) - Can work in parallel
3. After US1/US4/US2 complete:
   - Developer A: User Story 3 (Performance) - Optimizes US1
   - Developer B: User Story 5 (Mouse Input) - Independent
   - Developer C: User Story 6 (gokeys Integration) - Coordinates US5
4. Stories integrate and test together

---

## Notes

- **[P] tasks**: Different files, no dependencies, safe to run in parallel
- **[Story] label**: Maps task to specific user story for traceability
- **Tests first**: Per constitution TDD requirement, write and fail tests before implementing
- **80% coverage minimum**: Required by constitution quality gates
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Run `go test ./...` frequently to ensure tests pass
- Run `go test -bench ./...` to verify performance targets
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence

---

## Success Metrics

**Total Tasks**: 103
**Tasks per User Story**:
- Setup (Phase 1): 10 tasks
- Foundational (Phase 2): 9 tasks
- User Story 1 (P1): 17 tasks (4 test + 13 implementation)
- User Story 2 (P2): 13 tasks (4 test + 9 implementation)
- User Story 3 (P3): 12 tasks (4 test + 8 implementation)
- User Story 4 (P4): 10 tasks (3 test + 7 implementation)
- User Story 5 (P5): 12 tasks (3 test + 9 implementation)
- User Story 6 (P6): 9 tasks (2 test + 7 implementation)
- Polish (Phase 9): 11 tasks

**Parallel Opportunities Identified**: 45 tasks marked [P]
**MVP Scope**: Phase 1 + Phase 2 + Phase 3 (User Story 1) = 36 tasks
**Independent Test Criteria**: Defined for each user story

**Constitution Compliance**:
- ✅ TDD enforced (tests before implementation)
- ✅ 80% coverage target (T101)
- ✅ Performance benchmarks (T052, T053, T102)
- ✅ Go idioms (go fmt, go vet, golangci-lint)
- ✅ Platform abstraction (build tags, interface-based)
