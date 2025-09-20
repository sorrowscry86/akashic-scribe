# Akashic Scribe Developer Guide

This guide provides comprehensive information for developers who want to contribute to Akashic Scribe or build upon its functionality.

## Table of Contents

1. [Development Environment Setup](#development-environment-setup)
2. [Architecture Overview](#architecture-overview)
3. [Code Organization](#code-organization)
4. [Building and Testing](#building-and-testing)
5. [Contributing Guidelines](#contributing-guidelines)
6. [Plugin Development](#plugin-development)

## Development Environment Setup

### Prerequisites

- Go 1.24 or later
- Git
- Fyne dependencies (see below)

### Setting Up Fyne

Fyne requires some platform-specific dependencies:

#### Windows
- MinGW-w64 or MSYS2 with gcc
- Go environment properly configured

#### macOS
- Xcode Command Line Tools
- Go environment properly configured

#### Linux
- X11 and XCB development libraries
- Go environment properly configured

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/sorrowscry86/akashic-scribe.git
   ```

2. Navigate to the project directory:
   ```bash
   cd akashic-scribe/akashic_scribe
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run .
   ```

## Architecture Overview

Akashic Scribe follows a clean architecture approach with clear separation of concerns:

- **GUI Layer**: User interface components built with Fyne
- **Core Layer**: Business logic and processing engines
- **Plugin System**: Extensibility framework for custom functionality

### Component Diagram

```
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│                 │      │                 │      │                 │
│  GUI Layer      │◄────►│  Core Layer     │◄────►│  Plugin System  │
│  (Fyne)         │      │  (Engines)      │      │  (Extensions)   │
│                 │      │                 │      │                 │
└─────────────────┘      └─────────────────┘      └─────────────────┘
```

## Code Organization

```
akashic_scribe/
├── main.go              # Application entry point
├── gui/                 # User interface components
│   ├── layout.go        # UI layout and widgets
│   ├── state.go         # State management
│   ├── theme.go         # Custom theming
│   └── engine.go        # GUI-Core interface
├── core/                # Business logic
│   ├── engine.go        # Core engine interface
│   ├── real_engine.go   # Production implementation
│   ├── mock_engine.go   # Testing implementation
│   └── options.go       # Configuration structure
└── assets/              # Static resources
    ├── icons/
    └── fonts/
```

## Building and Testing

### Building for Different Platforms

```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o akashic_scribe.exe .

# macOS
GOOS=darwin GOARCH=amd64 go build -o akashic_scribe .

# Linux
GOOS=linux GOARCH=amd64 go build -o akashic_scribe .
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

## Contributing Guidelines

We welcome contributions to Akashic Scribe! Please follow these guidelines:

1. **Fork the repository** and create a feature branch
2. **Write tests** for your changes
3. **Ensure code quality** by running linters
4. **Submit a pull request** with a clear description of your changes

### Code Style

- Follow standard Go conventions
- Use meaningful variable and function names
- Write comprehensive comments
- Ensure all tests pass

## Plugin Development

Akashic Scribe supports plugins for extending functionality. To create a plugin:

1. Implement the appropriate plugin interface
2. Register your plugin with the plugin manager
3. Package your plugin according to the plugin specification

For detailed plugin development instructions, see the [Plugin Development Guide](plugins.md).

---

For API documentation, please refer to the [API Documentation](../api/index.md).