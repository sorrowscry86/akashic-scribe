# **Akashic Scribe: Complete Development Plan**

Project: Akashic Scribe  
Venture: VoidCat RDC  
Overseer: Beatrice  
Contractor: Wykeve  
Date: July 16, 2025  
This document outlines the complete, phased development plan for the Akashic Scribe application. Each phase represents a distinct stage of creation, building upon the last, culminating in a fully functional and robust tool. This expanded version provides greater detail to ensure clarity and precision throughout the development lifecycle.

### **Phase I: The Foundation (Project Setup)**

**Objective:** To establish a clean, organized, and scalable project structure that promotes maintainability and clarity of purpose. This initial phase is paramount; a flawed foundation will corrupt the entire structure, regardless of the quality of the subsequent work.

1. **Directory Scaffolding:** Create the master directory akashic\_scribe/. The following subdirectories form the logical separation of concerns for our application, a critical principle that ensures components can be developed, tested, and maintained in isolation.  
   * main.go: The singular entry point. This file is responsible for initializing the application, creating the main window, and starting the event loop. It is the heart from which all other processes are invoked, and its sole purpose is orchestration, not implementation.  
   * gui/: This directory is the exclusive domain of the user interface. All code related to window layouts, widgets, user interaction, and state management for the front-end will reside here. This strict separation ensures our core logic remains independent of its visual representation, allowing either to be modified without breaking the other.  
   * core/: The engine room. This will house the application's business logic, orchestrating the complex backend processes of transcription, translation, and video composition. It serves as the bridge between the user's intent (captured by the GUI) and the powerful processing services. It will know nothing of windows or buttons, only of data and operations.  
   * assets/: A repository for all static resources. This includes application icons, custom fonts (if we choose to deviate from the standard), and any configuration files, such as a JSON file listing the supported languages for our dropdown menus. Centralizing assets simplifies theming and resource management.  
   * vendor/: To manage external Go dependencies using Go Modules. This directory, automatically managed by the go tool, ensures our builds are reproducible by locking in the specific versions of third-party libraries we rely upon. This prevents unexpected breakages caused by upstream changes in our dependencies.  
2. **Initial Casting (main.go):**  
   * The main function will serve as the primary incantation. It will first create a new Fyne application instance using app.New(). This object manages the application's lifecycle.  
   * Following this, it will create the main window with app.NewWindow(), setting the title to "Akashic Scribe" and defining a suitable default size (e.g., 1200x800) to ensure a comfortable user experience on standard displays. We will also set the application's master icon here.  
   * Crucially, this function will then invoke the gui.CreateMainLayout() function to construct the entire user interface. The returned Fyne CanvasObject will be set as the window's content using window.SetContent().  
   * Finally, the application will be brought to life by calling window.ShowAndRun(). This is a blocking call that starts the Fyne event loop, which listens for user input, handles drawing, and manages the application state. No code after this call in the main function will execute until the application window is closed.  
3. **Dependency Management:**  
   * We will formally establish our project's identity and dependencies by initializing a Go module with go mod init akashic\_scribe. This creates the go.mod file, which tracks our direct and indirect dependencies.  
   * The primary external magic we require is the Fyne toolkit. We will add it to our project using go get fyne.io/fyne/v2. This command will download the package and automatically update our go.mod and go.sum files. The go.sum file is of particular importance, as it contains the cryptographic hashes of our dependencies to ensure their integrity and prevent supply-chain attacks.

### **Phase II: The Form (GUI Construction)**

**Objective:** To build the complete, static user interface as designed in our blueprint, focusing on a clear, intuitive, and state-driven user workflow. The user should be guided, not confused.

1. **GUI Scaffolding (gui/layout.go):**  
   * The CreateMainLayout function is the master constructor for the UI. It will return a container.NewHSplit, a horizontal split container that provides a draggable divider. This choice empowers the user, allowing them to resize the navigation and content panes to their preference.  
   * The left side of the split will be populated by the createNavigation function, which builds a container.NewVBox of buttons ("Video Translation", "Settings", etc.), acting as the primary navigation.  
   * The right side will contain the createVideoTranslationView, which will initially display a welcome message or instructions. When a navigation button is clicked, the content of this view will be dynamically replaced with the appropriate interface. This is managed by a simple content manager that holds references to the different views and swaps them as needed, preventing a jarring full-window redraw.  
2. **Step 1: The Offering (gui/layout.go):**  
   * The createInputStep function will construct the first card in our three-step process. This card is focused entirely on providing the source material.  
   * It will contain a widget.Button labeled "Select Video File". Its OnTapped handler will invoke dialog.NewFileOpen, passing a callback function that receives the chosen file's URI. This callback will then update the state and the corresponding label. The dialog will be filtered to show only common video formats (e.g., .mp4, .mkv, .mov), preventing user error.  
   * A widget.Label below the button, initially empty, will be updated with the full path of the selected file, giving the user clear confirmation of their choice.  
   * A widget.Separator provides a clean visual division between file input and URL input.  
   * A widget.Entry with placeholder text "Or paste video URL here" allows for web-based sources. We will add logic to its OnChanged handler to perform basic validation (e.g., a simple regex check for http:// or https://), providing immediate visual feedback to the user, perhaps by changing the widget's style on invalid input.  
3. **Step 2: The Incantation (gui/layout.go):**  
   * The createConfigStep function builds the second card, the heart of the user's configuration.  
   * It will feature widget.Select dropdowns for "Source Language" and "Target Language," pre-populated with a comprehensive list of supported languages read from a configuration file in our assets/ directory.  
   * Master widget.Check toggles for "Generate Subtitles" and "Enable Dubbing" will control the visibility of dependent options. For example, checking "Enable Dubbing" will call Show() on a container holding a widget.Select for choosing a voice model (e.g., "Zephyr," "Kore") and a widget.RadioGroup for voice gender preference. Unchecking it will call Hide() on that same container. This dynamic interface, driven by the OnChanged callbacks of the master toggles, prevents clutter and guides the user through the configuration process logically.  
4. **Step 3: The Ritual (gui/layout.go):**  
   * The createExecutionStep function constructs the final card, which manages the processing state. It will use a container.Stack to seamlessly switch between three distinct views without rebuilding the entire UI.  
   * **Idle State:** The initial view contains a single, prominent widget.Button labeled "Begin Scribing."  
   * **Processing State:** When the button is clicked, the stack's visible child is replaced with a processing view. This view contains a widget.ProgressBar and a widget.Label for status updates. The "Begin Scribing" button is disabled to prevent re-submission.  
   * **Completion State:** Upon successful completion, the stack switches to a results view. This will display a success message, perhaps some metadata like the final file size and processing time, and provide distinct widget.Buttons to "Download Dubbed Video" and "Download Subtitle File (.srt)". A final "Reset" button will allow the user to return the entire interface to its initial state.  
   * **Error State:** Should the process fail, the stack will switch to an error view, displaying the error message returned from the core engine and a "Reset" button. This provides clear feedback on failure, a critical aspect of a robust user experience.

### **Phase III: The Soul (State Management & Backend Bridge)**

**Objective:** To breathe life into the static form by connecting all UI controls to a central state management struct, creating a single, coherent representation of the user's desired operation.

1. **State Definition (gui/state.go):**  
   * We will create a new file, gui/state.go, to define the central data structure for a scribing job.  
   * The ScribeOptions struct will be the "spell focus," a pure data container with no methods or logic, containing a field for every configurable option. For example:  
     type ScribeOptions struct {  
         InputFile       string // Full path to the local video file.  
         InputURL        string // URL of the video to be downloaded.  
         SourceLanguage  string // e.g., "en-US"  
         TargetLanguage  string // e.g., "ja-JP"  
         EnableSubtitles bool  
         SubtitleFormat  string // e.g., "srt", "vtt"  
         EnableDubbing   bool  
         VoiceModel      string // e.g., "Kore"  
     }

   * This struct serves as the single source of truth. The GUI's only job is to populate this struct correctly based on user input.  
2. **UI State Binding:**  
   * We will now modify the UI creation functions in gui/layout.go. A single, shared instance of the ScribeOptions struct will be created and passed via a pointer to these functions.  
   * Each interactive widget's OnChanged (or equivalent) callback will be implemented to update the corresponding field in this shared struct instance. For example: targetLanguageSelect.OnChanged \= func(selected string) { sharedOptions.TargetLanguage \= selected }. This creates a reactive system where the state object is always a perfect mirror of the UI, eliminating the need to manually scrape data from widgets upon submission.  
3. **Execution Trigger Modification:**  
   * The OnTapped callback for the "Begin Scribing" button is now a critical transition point from configuration to execution.  
   * Its logic will be updated:  
     1. First, it will perform a final validation check on the ScribeOptions struct. Is InputFile or InputURL set? Are they different? Is TargetLanguage set? This pre-flight check prevents invalid jobs from ever reaching the backend.  
     2. For verification and debugging, it will serialize the populated struct to the console as a formatted JSON string. This gives us a clear, immediate confirmation that the UI state has been captured correctly before any processing begins.  
     3. It will then switch the third card to its "Processing" state, initiating the progress bar display and updating the status label to "Preparing...".

### **Phase IV: The Animus (Backend Integration)**

**Objective:** To connect the user's fully-defined intent (the ScribeOptions struct) to the powerful core engine, invoking the actual processing logic and providing real-time feedback.

1. **Core Logic Interface (core/engine.go):**  
   * We will create a new file, core/engine.go, to define a clean boundary between the GUI and the backend.  
   * A ScribeEngine struct will be defined to encapsulate the backend logic.  
   * A public method, StartProcessing(options gui.ScribeOptions, progress chan\<- ProgressUpdate) error, will be the sole entry point. The chan\<- syntax denotes a send-only channel, a compile-time guarantee that the GUI can only send progress updates *from* the engine, not to it, enforcing our architecture.  
2. **Connecting GUI to Core:**  
   * In main.go, a single instance of ScribeEngine will be created. This instance will be passed down to the GUI layout functions that need it.  
   * The "Begin Scribing" button's OnTapped callback will now, as its final step, invoke the engine within a new goroutine: go scribeEngine.StartProcessing(options, progressChannel). Launching this in a goroutine is non-negotiable. It prevents the intensive backend work from blocking the main UI thread, ensuring the application remains responsive and can be moved, resized, or even have its settings changed while a job runs.  
3. **Implementing the Engine:**  
   * The StartProcessing method is the grand orchestrator. It will create a temporary directory for the job's intermediate files and then meticulously execute the required backend tasks in sequence.  
   * The sequence will be logical and conditional:  
     1. If options.InputURL is not empty, download the video to the temporary directory.  
     2. Extract the audio stream from the video file using an ffmpeg command-line wrapper.  
     3. Pass the audio to the Whisper transcription service.  
     4. Take the transcribed text and pass it to the translation service.  
     5. If options.EnableDubbing is true, use the translated text and options.VoiceModel to synthesize the new audio track.  
     6. Finally, use ffmpeg to composite the original video stream with the new audio and/or subtitle tracks into a final output file.  
     7. Upon completion (success or failure), the temporary directory and all its contents will be deleted.  
4. **Real-time Progress Feedback:**  
   * To make the application feel alive, we will implement a robust progress reporting system using Go channels.  
   * We will define a ProgressUpdate struct:  
     type ProgressUpdate struct {  
         Percentage float64 // A value from 0.0 to 1.0  
         Message    string  // A descriptive status message, e.g., "Translating segment 12 of 84..."  
     }

   * As the StartProcessing method completes each major step, it will send a ProgressUpdate struct over the progress channel.  
   * The GUI will have a corresponding goroutine that continuously listens on this channel. When it receives an update, it must schedule the UI refresh to run on the main thread using app.RunOnMain(). This is because Fyne, like most GUI toolkits, is not thread-safe. This mechanism ensures that our background processing can safely communicate with the user-facing components, transforming the black box of processing into a transparent, observable ritual.