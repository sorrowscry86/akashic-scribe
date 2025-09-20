package gui

import (
	"akashic_scribe/core"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// CreateMainLayout constructs the primary UI structure.
func CreateMainLayout(window fyne.Window, engine core.ScribeEngine) fyne.CanvasObject {
	// Centralized state object that all UI components will modify.
	options := &ScribeOptions{}
	// createVideoTranslationView includes navigation and content composition
	return createVideoTranslationView(window, options, engine)
}

// createNavigation builds the left-side navigation buttons.
func createNavigation(stack *fyne.Container, steps fyne.CanvasObject, settings fyne.CanvasObject) fyne.CanvasObject {
	// Simple left-side navigation with two views
	return container.NewVBox(
		widget.NewButtonWithIcon("Video Translation", theme.HomeIcon(), func() {
			steps.Show()
			settings.Hide()
			stack.Refresh()
		}),
		widget.NewButtonWithIcon("Settings", theme.SettingsIcon(), func() {
			settings.Show()
			steps.Hide()
			stack.Refresh()
		}),
	)
}

// createVideoTranslationView builds the interface for our primary feature.
func createVideoTranslationView(window fyne.Window, options *ScribeOptions, engine core.ScribeEngine) fyne.CanvasObject {
	// Each step of the UI is a card, and they are arranged vertically.
	steps := container.NewVBox(
		createInputStep(window, options),
		createConfigStep(window, options),
		createExecutionStep(window, options, engine),
	)

	// Settings panel
	settings := createSettingsView(window, options)

	// Stack to switch between steps and settings
	stack := container.NewStack(steps, settings)
	settings.Hide()

	// Wire up settings button to toggle view
	nav := createNavigation(stack, steps, settings)

	// Replace split to use our new stack
	return container.NewBorder(nil, nil, nav, nil, stack)
}

// createSettingsView builds a simple settings panel for output directory selection.
func createSettingsView(window fyne.Window, options *ScribeOptions) fyne.CanvasObject {
	current := widget.NewEntry()
	current.Disable()
	current.SetPlaceHolder("No folder selected")
	if options.OutputDir != "" {
		current.SetText(options.OutputDir)
	}

	pickBtn := widget.NewButtonWithIcon("Choose Output Folder", theme.FolderOpenIcon(), func() {
		dlg := dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if list == nil {
				return
			}
			options.OutputDir = list.Path()
			current.SetText(options.OutputDir)
		}, window)
		dlg.Show()
	})

	resetBtn := widget.NewButton("Use Default", func() {
		options.OutputDir = ""
		current.SetText("")
	})

	content := container.NewVBox(
		widget.NewLabelWithStyle("Output Directory", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		current,
		container.NewHBox(pickBtn, resetBtn),
	)

	return widget.NewCard("Settings", "Configure application preferences.", content)
}

// createInputStep builds the UI for Step 1: The Offering.
func createInputStep(window fyne.Window, options *ScribeOptions) *widget.Card {
	// A label to display the name of the selected file.
	selectedFileLabel := widget.NewLabel("No file selected, in fact.")
	selectedFileLabel.Alignment = fyne.TextAlignCenter
	selectedFileLabel.Wrapping = fyne.TextWrapWord

	// The button that opens the file selection dialog.
	fileSelectBtn := widget.NewButtonWithIcon("Select Video File", theme.FileVideoIcon(), func() {
		fileOpenDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				// User cancelled the dialog.
				return
			}
			defer reader.Close()

			// Update the label and the central options struct.
			selectedFileLabel.SetText("Selected: " + reader.URI().Name())
			options.InputFile = reader.URI().Path()
			// Clear the URL field if a file is selected.
			// urlEntry.SetText("")
			options.InputURL = ""
		}, window)
		fileOpenDialog.SetFilter(storage.NewExtensionFileFilter([]string{".mp4", ".ts", ".webm", ".mkv"}))
		fileOpenDialog.Show()
	})

	// The text entry field for a URL.
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("...or paste a video link here, I suppose.")
	urlEntry.OnChanged = func(text string) {
		options.InputURL = text
		if text != "" {
			// Clear the file selection if a URL is entered.
			selectedFileLabel.SetText("No file selected, in fact.")
			options.InputFile = ""
		}
	}

	// Assemble the components in a vertical box.
	inputContainer := container.NewVBox(
		fileSelectBtn,
		selectedFileLabel,
		widget.NewSeparator(),
		urlEntry,
	)

	return widget.NewCard("Step 1: The Offering", "Provide the source material.", inputContainer)
}

// getLanguageOptions provides the available language options for translation.
func getLanguageOptions() []string {
	return []string{
		"English", "日本語 (Japanese)", "简体中文 (Simplified Chinese)", "Русский язык (Russian)",
		"Deutsch (German)", "Français (French)", "Español (Spanish)", "Italiano (Italian)",
		"Português (Portuguese)", "한국어 (Korean)",
	}
}

// createConfigStep builds the UI for Step 2: The Incantation.
func createConfigStep(window fyne.Window, options *ScribeOptions) *widget.Card {
	// --- Language Selection ---
	languageOptions := getLanguageOptions()

	originLangSelect := widget.NewSelect(languageOptions, func(s string) {
		options.OriginLanguage = s
	})
	originLangSelect.PlaceHolder = "Select Original Language"

	targetLangSelect := widget.NewSelect(languageOptions, func(s string) {
		options.TargetLanguage = s
	})
	targetLangSelect.PlaceHolder = "Select Target Language"

	langContainer := container.New(layout.NewFormLayout(),
		widget.NewLabel("Original Language:"), originLangSelect,
		widget.NewLabel("Translate To:"), targetLangSelect,
	)

	// --- Subtitle Configuration ---
	subtitlePosition := widget.NewRadioGroup([]string{"Translation on Top", "Translation on Bottom"}, func(s string) {
		options.SubtitlePosition = s
	})
	subtitlePosition.Horizontal = true
	subtitlePosition.SetSelected("Translation on Top")
	options.SubtitlePosition = "Translation on Top" // Set initial state

	bilingualCheck := widget.NewCheck("Bilingual Subtitles", func(checked bool) {
		options.BilingualSubtitles = checked
		if checked {
			subtitlePosition.Enable()
		} else {
			subtitlePosition.Disable()
		}
	})
	bilingualCheck.SetChecked(true)
	options.BilingualSubtitles = true // Set initial state

	subtitleToggle := widget.NewCheck("Create Subtitles", func(checked bool) {
		options.CreateSubtitles = checked
		if checked {
			bilingualCheck.Enable()
			if bilingualCheck.Checked {
				subtitlePosition.Enable()
			}
		} else {
			bilingualCheck.Disable()
			subtitlePosition.Disable()
		}
	})
	subtitleToggle.SetChecked(true)
	options.CreateSubtitles = true // Set initial state

	subtitleContainer := container.NewVBox(subtitleToggle, container.NewPadded(bilingualCheck), container.NewPadded(subtitlePosition))

	// --- Dubbing Configuration ---
	voiceSelect := widget.NewSelect([]string{"alloy", "echo", "fable", "onyx", "nova", "shimmer"}, func(s string) {
		options.VoiceModel = s
	})
	voiceSelect.PlaceHolder = "Select Voice Model"

	voiceCloneBtn := widget.NewButton("Use Custom Voice...", func() {
		// This feature is not yet implemented.
		// When implemented, set options.UseCustomVoice = true and update any related state.
		options.UseCustomVoice = true
	})

	dubbingOptions := container.NewVBox(voiceSelect, voiceCloneBtn)
	dubbingOptions.Hide() // Initially hidden.

	dubbingToggle := widget.NewCheck("Create Dubbing", func(checked bool) {
		options.CreateDubbing = checked
		if checked {
			dubbingOptions.Show()
		} else {
			dubbingOptions.Hide()
		}
	})

	dubbingContainer := container.NewVBox(dubbingToggle, container.NewPadded(dubbingOptions))

	// --- Final Assembly of the Card ---
	configContent := container.NewVBox(
		langContainer,
		widget.NewSeparator(),
		subtitleContainer,
		widget.NewSeparator(),
		dubbingContainer,
	)

	return widget.NewCard("Step 2: The Incantation", "Define the transformation.", configContent)
}

// createExecutionStep builds the UI for Step 3: The Ritual.
//
// === BACKEND INTEGRATION PHASE ===
// This function is architected for seamless backend integration:
//   - All user options are collected in the ScribeOptions struct.
//   - Progress updates are handled via a channel (see ProgressUpdate).
//   - Replace the simulation code with a real call to engine.StartProcessing(*options, progressChan).
//   - UI updates are performed in response to progress channel events.
//   - On completion, results are displayed and download buttons are enabled.
//
// TODOs for integration:
//   - Replace ProgressUpdate simulation with actual backend progress reporting.
//   - Ensure error handling and UI reset on backend failure.
//   - Wire up download buttons to real output files.
//   - Remove simulation notifications after backend is live.
func createExecutionStep(window fyne.Window, options *ScribeOptions, engine core.ScribeEngine) *widget.Card {
	progress := widget.NewProgressBar()
	statusLabel := widget.NewLabel("Status: Awaiting command...")
	statusLabel.Alignment = fyne.TextAlignCenter
	progressContainer := container.NewVBox(progress, statusLabel)

	downloadContainer := container.NewVBox()
	viewStack := container.NewStack()

	var startButton *widget.Button
	resetButton := widget.NewButton("Scribe Another", func() {
		viewStack.Objects = []fyne.CanvasObject{startButton}
		viewStack.Refresh()
	})

	startButton = widget.NewButtonWithIcon("Begin Scribing", theme.ConfirmIcon(), func() {
		// The 'options' struct is now always up-to-date.
		fmt.Println("🏰 ALBEDO ASSISTANT - SCRIBE OPTIONS GATHERED:")
		fmt.Println(options.String())

		// Basic validation before starting.
		if options.InputFile == "" && options.InputURL == "" {
			dialog.ShowInformation("Missing Input", "Please select a file or provide a URL before starting.", window)
			return
		}
		if options.OriginLanguage == "" || options.TargetLanguage == "" {
			dialog.ShowInformation("Missing Language", "Please select both an original and a target language.", window)
			return
		}

		progress.SetValue(0)
		statusLabel.SetText("Status: Ready for backend integration...")
		downloadContainer.RemoveAll()
		viewStack.Objects = []fyne.CanvasObject{progressContainer}
		viewStack.Refresh()

		// === BACKEND INTEGRATION POINT ===
		// Use the real backend engine and progress channel
		progressChan := make(chan core.ProgressUpdate)

		// Start backend processing in a goroutine
		go func() {
			err := engine.StartProcessing(*options, progressChan)
			if err != nil {
				// Create user-friendly error message
				errorTitle := "Processing Error"
				errorMsg := err.Error()

				// Categorize common errors for better user experience
				if strings.Contains(errorMsg, "yt-dlp not found") {
					errorTitle = "Missing Dependency: yt-dlp"
					errorMsg = "yt-dlp is required to download videos from URLs. Please install it using:\nwinget install yt-dlp"
				} else if strings.Contains(errorMsg, "ffmpeg not found") {
					errorTitle = "Missing Dependency: FFmpeg"
					errorMsg = "FFmpeg is required to process video files. Please install it using:\nwinget install ffmpeg"
				} else if strings.Contains(errorMsg, "input file not found") {
					errorTitle = "File Not Found"
					errorMsg = "The selected video file could not be found. Please check the file path and try again."
				} else if strings.Contains(errorMsg, "failed to download") {
					errorTitle = "Download Failed"
					errorMsg = "Failed to download the video from the provided URL. Please check the URL and your internet connection."
				}

				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   errorTitle,
					Content: errorMsg,
				})
				fyne.Do(func() {
					statusLabel.SetText("Status: " + errorTitle)
					progress.SetValue(0)
					viewStack.Objects = []fyne.CanvasObject{startButton}
					viewStack.Refresh()
				})
			}
		}()

		// Listen for progress updates and update the UI accordingly
		go func() {
			defer close(progressChan) // Close when done reading
			var finalTranscription, finalTranslation, outputDir string
			for update := range progressChan {
				fyne.Do(func() {
					progress.SetValue(update.Percentage)
					statusLabel.SetText(fmt.Sprintf("Status: %s", update.Message))
				})
				// Try to parse result JSON if present in the final message
				if update.Percentage >= 1.0 && len(update.Message) > 0 {
					// Extract output directory
					lines := strings.Split(update.Message, "\n")
					for _, line := range lines {
						if strings.HasPrefix(line, "Output saved to: ") {
							outputDir = strings.TrimPrefix(line, "Output saved to: ")
							break
						}
					}
					// Look for result JSON
					if idx := strings.Index(update.Message, "{"); idx != -1 {
						resultJSON := update.Message[idx:]
						type resultStruct struct {
							Transcription string `json:"Transcription"`
							Translation   string `json:"Translation"`
						}
						var result resultStruct
						if err := json.Unmarshal([]byte(resultJSON), &result); err == nil {
							finalTranscription = result.Transcription
							finalTranslation = result.Translation
						}
					}
				}
			}
			// Show completion UI
			fyne.Do(func() {
				statusLabel.SetText("Status: Scribing complete.")
				downloadContainer.Add(widget.NewLabelWithStyle("Scribing Complete.", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}))
				if outputDir != "" {
					downloadContainer.Add(widget.NewLabelWithStyle("Files saved to: "+outputDir, fyne.TextAlignCenter, fyne.TextStyle{}))
				}
				entry := widget.NewMultiLineEntry()
				if finalTranscription != "" || finalTranslation != "" {
					entry.SetText(fmt.Sprintf("Transcription:\n%s\n\nTranslation:\n%s", finalTranscription, finalTranslation))
				} else {
					entry.SetText("Transcription:\n(See logs)\n\nTranslation:\n(See logs)")
				}
				entry.Disable()
				downloadContainer.Add(entry)
				downloadContainer.Add(widget.NewButton("Open Output Folder", func() {
					if outputDir != "" {
						if err := exec.Command("explorer", outputDir).Start(); err != nil {
							dialog.ShowError(err, window)
						}
					}
				}))
				downloadContainer.Add(widget.NewSeparator())
				downloadContainer.Add(resetButton)
				viewStack.Objects = []fyne.CanvasObject{downloadContainer}
				viewStack.Refresh()
			})
		}()
	})
	startButton.Importance = widget.HighImportance
	viewStack.Add(startButton)

	return widget.NewCard("Step 3: The Ritual", "Initiate the process and receive the results.", viewStack)
}
