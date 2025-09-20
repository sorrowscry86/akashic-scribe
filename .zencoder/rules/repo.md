---
description: Repository Information Overview
alwaysApply: true
---

# Akashic Scribe Information

## Summary
Akashic Scribe is a desktop application built with Go and the Fyne toolkit. It provides a user interface for video translation, allowing users to select video files, configure translation settings, and generate subtitles and dubbed audio.

## Structure
- `akashic_scribe/`: Main application directory
  - `main.go`: Application entry point
  - `gui/`: UI components and layout
    - `layout.go`: Main UI layout and components
    - `state.go`: Data structures for application state
    - `theme.go`: Custom theme placeholder
  - `assets/`: Directory for static resources like icons and fonts
  - `test_state.go`: Test file for state management

## Language & Runtime
**Language**: Go
**Version**: Go 1.24.4
**Build System**: Go Modules
**Package Manager**: Go Modules

## Dependencies
**Main Dependencies**:
- `fyne.io/fyne/v2` v2.6.1: Cross-platform GUI toolkit for Go

**Indirect Dependencies**:
- `fyne.io/systray` v1.11.0: System tray integration
- `github.com/go-gl/glfw/v3.3/glfw`: OpenGL framework
- Various rendering and UI component libraries

## Build & Installation
```bash
cd akashic_scribe
go mod download
go build
```

## Main Application Components
- **GUI Layer**: Built with Fyne toolkit, providing a responsive user interface
- **State Management**: Centralized state handling via the `ScribeOptions` struct
- **Layout System**: Three-step workflow (Input, Configuration, Execution)

### Main Entry Point
The application starts in `main.go`, which initializes the Fyne application, creates the main window, and sets up the UI layout.

### UI Structure
- **Navigation**: Left sidebar with buttons for different views
- **Content Area**: Right side containing the active view
- **Video Translation View**: Three-step card layout
  - Step 1: Input selection (file or URL)
  - Step 2: Configuration (languages, subtitles, dubbing)
  - Step 3: Execution and results

### State Management
The application uses a centralized state management approach:
- `ScribeOptions` struct in `state.go` defines all user configuration options
- `ScribeWidgets` struct in `layout.go` holds references to UI widgets
- Widget references are collected and used to populate the state

## Testing Framework
**E2E Testing**: Go native testing framework with testify for assertions
**Test Structure**: Comprehensive test suites covering:
- End-to-end application workflows
- GUI component testing with Fyne test framework
- Core engine integration testing
- Performance and memory usage baselines

**Test Files**:
- `e2e_test.go`: Main end-to-end workflow tests
- `gui/gui_e2e_test.go`: GUI component and interaction tests
- `core/integration_test.go`: Core engine integration tests
- `core/mock_engine_test.go`: Unit tests for mock engine

**Test Execution**:
```bash
# Run all tests
go test ./...

# Run E2E tests only
go test -run TestE2EScenarios

# Run with verbose output
go test -v ./...

# Skip long-running tests
go test -short ./...
```

## Future Development
According to the development plan, the following components are planned but not yet implemented:
- `core/`: Backend processing logic for transcription, translation, and video composition
- Custom theme implementation
- Integration with external services for transcription and translation