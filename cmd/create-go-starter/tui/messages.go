package tui

// Messages define all the events that can occur in the application.
// Following Elm Architecture, these are the inputs to the Update function.

// WelcomeMenuSelectedMsg is sent when user selects an option from the welcome menu.
type WelcomeMenuSelectedMsg struct {
	Action string // "create", "help", or "exit"
}

// ProjectNameSubmittedMsg is sent when the user submits the project name.
type ProjectNameSubmittedMsg struct {
	Name string
}

// TemplateSelectedMsg is sent when the user selects a template.
type TemplateSelectedMsg struct {
	Template string
}

// DatabaseSelectedMsg is sent when the user selects a database.
type DatabaseSelectedMsg struct {
	Database string
}

// FrameworkSelectedMsg is sent when the user selects a web framework.
type FrameworkSelectedMsg struct {
	Framework string
}

// ObservabilitySelectedMsg is sent when the user selects observability level.
type ObservabilitySelectedMsg struct {
	Observability string
}

// GenerationStartedMsg is sent when file generation begins.
type GenerationStartedMsg struct {
	TotalFiles int
}

// FileGeneratedMsg is sent when a single file has been generated.
// Story 10.7 Task 3: Enhanced with file size and step information for real-time stats
type FileGeneratedMsg struct {
	FilePath string
	Index    int    // Current file number (1-based)
	Size     int64  // Size of the generated file in bytes
	Step     string // Current generation step (e.g., "Creating directories...")
}

// GenerationCompleteMsg is sent when all files have been generated.
type GenerationCompleteMsg struct {
	TotalFiles int
	Success    bool
	Error      error
}

// StepCompletedMsg is sent when a generation step completes (for progress tracking).
type StepCompletedMsg struct {
	StepName string
	Progress float64 // 0.0 to 1.0
}

// NavigateBackMsg is sent when the user wants to go back to the previous screen.
type NavigateBackMsg struct{}

// ConfirmGenerationMsg is sent when the user confirms the configuration and starts generation.
type ConfirmGenerationMsg struct{}
