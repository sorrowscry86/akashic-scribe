---
description: AI rules derived by SpecStory from the project AI interaction history
globs: *
---

## HEADERS

## TECH STACK
*   Go
*   Fyne toolkit

## PROJECT DOCUMENTATION & CONTEXT SYSTEM
*   Project Name: Akashic Scribe

## PROJECT STRUCTURE
*   `akashic_scribe/`
    *   `main.go` (The entry point)
    *   `gui/` (Directory for UI elements)
        *   `layout.go`
        *   `theme.go`
    *   `assets/` (For icons or fonts)

## CODING STANDARDS
*   Use clean and logical structure.
*   Follow Fyne toolkit conventions.

## GUI GUIDELINES
*   Use `fyne.io/fyne/v2` for UI elements.
*   `CreateMainLayout` function in `gui/layout.go` constructs the main UI structure.
*   The main layout consists of a navigation menu on the left and a content area on the right, using `container.NewHSplit`.
*   Navigation is built using `createNavigation` function.
*   The main content area will display different views, such as the video translation view.
*   `createVideoTranslationView` function builds the interface for video translation.
*   The video translation view follows a 3-step layout: Input, Configuration, Execution.

## WORKFLOW & RELEASE RULES
*   Ignore `.github` and `.specstory` directories during development and deployment.

## DEBUGGING

## BEST PRACTICES