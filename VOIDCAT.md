# Project: Akashic Scribe

## Project Overview

Akashic Scribe is a desktop application designed for video transcription, translation, and dubbing. It provides a user-friendly interface for selecting a video file (either local or from a URL), choosing source and target languages, and configuring output options like subtitles and voice models.

The application is built using Go and the `fyne.io/fyne/v2` GUI toolkit. The architecture, as detailed in the `".github/Akashic Scribe_ Complete Development Plan (1).md"`, emphasizes a strict separation of concerns:

*   **`main.go`**: The application's entry point, responsible for initializing the Fyne app and window.
*   **`gui/`**: Contains all UI-related code, including layout construction (`layout.go`) and state management (`state.go`).
*   **`core/`**: Intended for the backend business logic (transcription, translation, etc.), though this is not yet implemented.
*   **`assets/`**: For static resources like icons, fonts, and configuration files.

The UI is designed as a three-step process, guiding the user through input, configuration, and execution. A central `ScribeOptions` struct (`gui/state.go`) acts as the single source of truth, capturing all user selections before triggering the (currently simulated) backend process.

## Building and Running

As a standard Go project, the following commands can be used from within the `akashic_scribe` directory:

*   **Run Dependencies:**
    ```sh
    go mod tidy
    ```

*   **Run the application:**
    ```sh
    go run .
    ```

*   **Build the application:**
    ```sh
    go build -o akashic_scribe.exe .
    ```

*   **Run tests:**
    ```sh
    go test ./...
    ```

## Development Conventions

The project follows the conventions laid out in the detailed development plan:

*   **Directory Structure:** Logic is separated into `gui`, `core`, and `assets` directories to maintain a clean and scalable codebase.
*   **State Management:** A single, pointer-shared instance of the `ScribeOptions` struct is used to manage the application's state. UI widgets update this struct directly via their `OnChanged` callbacks. This ensures the state is always consistent with the UI.
*   **UI Construction:** The UI is built programmatically using the Fyne toolkit. The layout is composed of nested containers and widgets, with a clear, card-based, three-step workflow.
*   **Backend Interaction:** The UI is designed to interact with a `core` engine (not yet implemented) in a non-blocking manner by launching processing tasks in a separate goroutine. Progress updates are intended to be communicated back to the UI via Go channels.
